package datastore

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/newrelic/go-agent/v3/integrations/nrmongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	logLine      = "MONGO %s %s: %+v\n"
	logErrorLine = "MONGO %s %s: %e: %+v\n"
)

// saveWithMongo will save a given struct to MongoDB
func (c *Client) saveWithMongo(
	ctx context.Context,
	model interface{},
	newRecord bool,
) (err error) {
	collectionName := GetModelTableName(model)
	if collectionName == nil {
		return ErrUnknownCollection
	}

	// Set the collection
	collection := c.options.mongoDB.Collection(
		setPrefix(c.options.mongoDBConfig.TablePrefix, *collectionName),
	)

	// Create or update
	if newRecord {
		c.DebugLog(ctx, fmt.Sprintf(logLine, "insert", *collectionName, model))
		_, err = collection.InsertOne(ctx, model)
	} else {
		id := GetModelStringAttribute(model, sqlIDFieldProper)
		update := bson.M{conditionSet: model}
		unset := GetModelUnset(model)
		if len(unset) > 0 {
			update = bson.M{conditionSet: model, conditionUnSet: unset}
		}

		c.DebugLog(ctx, fmt.Sprintf(logLine, "update", *collectionName, model))

		_, err = collection.UpdateOne(
			ctx, bson.M{mongoIDField: *id}, update,
		)
	}

	// Check for duplicate key (insert error, record exists)
	if mongo.IsDuplicateKeyError(err) {
		c.DebugLog(ctx, fmt.Sprintf(logErrorLine, "error", *collectionName, ErrDuplicateKey, model))
		return ErrDuplicateKey
	}

	if err != nil {
		c.DebugLog(ctx, fmt.Sprintf(logErrorLine, "error", *collectionName, err, model))
	}

	return
}

// incrementWithMongo will save a given struct to MongoDB
func (c *Client) incrementWithMongo(
	ctx context.Context,
	model interface{},
	fieldName string,
	increment int64,
) (newValue int64, err error) {
	collectionName := GetModelTableName(model)
	if collectionName == nil {
		return newValue, ErrUnknownCollection
	}

	// Set the collection
	collection := c.options.mongoDB.Collection(
		setPrefix(c.options.mongoDBConfig.TablePrefix, *collectionName),
	)

	id := GetModelStringAttribute(model, sqlIDFieldProper)
	if id == nil {
		return newValue, errors.New("can only increment by " + sqlIDField)
	}
	update := bson.M{conditionIncrement: bson.M{fieldName: increment}}

	c.DebugLog(ctx, fmt.Sprintf(logLine, "increment", *collectionName, model))

	result := collection.FindOneAndUpdate(
		ctx, bson.M{mongoIDField: *id}, update,
	)
	if result.Err() != nil {
		return newValue, result.Err()
	}
	var rawValue bson.Raw
	if rawValue, err = result.DecodeBytes(); err != nil {
		return
	}
	var newModel map[string]interface{}
	_ = bson.Unmarshal(rawValue, &newModel) // todo: cannot check error, breaks code atm

	newValue = newModel[fieldName].(int64) + increment

	if err != nil {
		c.DebugLog(ctx, fmt.Sprintf(logErrorLine, "error", *collectionName, err, model))
	}

	return
}

// CreateInBatchesMongo insert multiple models vai bulk.Write
func (c *Client) CreateInBatchesMongo(
	ctx context.Context,
	models interface{},
	batchSize int,
) error {
	collectionName := GetModelTableName(models)
	if collectionName == nil {
		return ErrUnknownCollection
	}

	mongoModels := make([]mongo.WriteModel, 0)
	collection := c.GetMongoCollection(*collectionName)
	bulkOptions := options.BulkWrite().SetOrdered(true)
	count := 0

	if reflect.TypeOf(models).Kind() == reflect.Slice {
		s := reflect.ValueOf(models)
		for i := 0; i < s.Len(); i++ {
			m := mongo.NewInsertOneModel()
			m.SetDocument(s.Index(i).Interface())
			mongoModels = append(mongoModels, m)
			count++

			if count%batchSize == 0 {
				_, err := collection.BulkWrite(ctx, mongoModels, bulkOptions)
				if err != nil {
					return err
				}
				// reset the bulk
				mongoModels = make([]mongo.WriteModel, 0)
			}
		}
	}

	if count%batchSize != 0 {
		_, err := collection.BulkWrite(ctx, mongoModels, bulkOptions)
		if err != nil {
			return err
		}
	}

	return nil
}

// getWithMongo will get given struct(s) from MongoDB
func (c *Client) getWithMongo(
	ctx context.Context,
	models interface{},
	conditions map[string]interface{},
	fieldResult interface{},
	queryParams *QueryParams,
) error {
	queryConditions := getMongoQueryConditions(models, conditions, c.GetMongoConditionProcessor())
	collectionName := GetModelTableName(models)
	if collectionName == nil {
		return ErrUnknownCollection
	}

	// Set the collection
	collection := c.options.mongoDB.Collection(
		setPrefix(c.options.mongoDBConfig.TablePrefix, *collectionName),
	)

	var fields []string
	if fieldResult != nil {
		fields = getFieldNames(fieldResult)
	}

	if IsModelSlice(models) {
		c.DebugLog(ctx, fmt.Sprintf(logLine, "findMany", *collectionName, queryConditions))

		var opts []*options.FindOptions
		if fields != nil {
			projection := bson.D{}
			for _, field := range fields {
				projection = append(projection, bson.E{Key: field, Value: 1})
			}
			opts = append(opts, options.Find().SetProjection(projection))
		}

		if queryParams.Page > 0 {
			opts = append(opts, options.Find().SetLimit(int64(queryParams.PageSize)).SetSkip(int64(queryParams.PageSize*(queryParams.Page-1))))
		}

		if queryParams.OrderByField == sqlIDField {
			queryParams.OrderByField = mongoIDField // use Mongo _id instead of default id field
		}
		if queryParams.OrderByField != "" {
			sortOrder := 1
			if queryParams.SortDirection == SortDesc {
				sortOrder = -1
			}
			opts = append(opts, options.Find().SetSort(bson.D{{Key: queryParams.OrderByField, Value: sortOrder}}))
		}

		cursor, err := collection.Find(ctx, queryConditions, opts...)
		if err != nil {
			return err
		}
		if err = cursor.Err(); errors.Is(err, mongo.ErrNoDocuments) {
			return ErrNoResults
		} else if err != nil {
			return cursor.Err()
		}

		if fieldResult != nil {
			if err = cursor.All(ctx, fieldResult); err != nil {
				return err
			}
		} else {
			if err = cursor.All(ctx, models); err != nil {
				return err
			}
		}
	} else {
		c.DebugLog(ctx, fmt.Sprintf(logLine, "find", *collectionName, queryConditions))

		var opts []*options.FindOneOptions
		if fields != nil {
			projection := bson.D{}
			for _, field := range fields {
				projection = append(projection, bson.E{Key: field, Value: 1})
			}
			opts = append(opts, options.FindOne().SetProjection(projection))
		}

		result := collection.FindOne(ctx, queryConditions, opts...)
		if err := result.Err(); errors.Is(err, mongo.ErrNoDocuments) {
			c.DebugLog(ctx, fmt.Sprintf(logLine, "result", *collectionName, "no result"))
			return ErrNoResults
		} else if err != nil {
			c.DebugLog(ctx, fmt.Sprintf(logLine, "result error", *collectionName, err))
			return result.Err()
		}

		if fieldResult != nil {
			if err := result.Decode(fieldResult); err != nil {
				c.DebugLog(ctx, fmt.Sprintf(logLine, "result error", *collectionName, err))
				return err
			}
		} else {
			if err := result.Decode(models); err != nil {
				c.DebugLog(ctx, fmt.Sprintf(logLine, "result error", *collectionName, err))
				return err
			}
		}
	}

	return nil
}

// countWithMongo will get a count of all models matching the conditions
func (c *Client) countWithMongo(
	ctx context.Context,
	models interface{},
	conditions map[string]interface{},
) (int64, error) {
	queryConditions := getMongoQueryConditions(models, conditions, c.GetMongoConditionProcessor())
	collectionName := GetModelTableName(models)
	if collectionName == nil {
		return 0, ErrUnknownCollection
	}

	// Set the collection
	collection := c.options.mongoDB.Collection(
		setPrefix(c.options.mongoDBConfig.TablePrefix, *collectionName),
	)

	c.DebugLog(ctx, fmt.Sprintf(logLine, accumulationCountField, *collectionName, queryConditions))

	count, err := collection.CountDocuments(ctx, queryConditions)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// aggregateWithMongo will get a count of all models aggregate by aggregateColumn matching the conditions
func (c *Client) aggregateWithMongo(
	ctx context.Context,
	models interface{},
	conditions map[string]interface{},
	aggregateColumn string,
	timeout time.Duration,
) (map[string]interface{}, error) {
	queryConditions := getMongoQueryConditions(models, conditions, c.GetMongoConditionProcessor())
	collectionName := GetModelTableName(models)
	if collectionName == nil {
		return nil, ErrUnknownCollection
	}

	// Set the collection
	collection := c.options.mongoDB.Collection(
		setPrefix(c.options.mongoDBConfig.TablePrefix, *collectionName),
	)

	c.DebugLog(ctx, fmt.Sprintf(logLine, accumulationCountField, *collectionName, queryConditions))

	// Marshal the data
	var matchStage bson.D
	data, err := bson.Marshal(queryConditions)
	if err != nil {
		return nil, err
	}

	// Unmarshal the bson
	if err = bson.Unmarshal(data, &matchStage); err != nil {
		return nil, err
	}

	aggregateOn := bson.E{
		Key:   mongoIDField,
		Value: "$" + aggregateColumn,
	} // default

	// Check for date field
	if StringInSlice(aggregateColumn, DateFields) {
		aggregateOn = bson.E{
			Key: mongoIDField,
			Value: bson.D{
				{
					Key: conditionDateToString,
					Value: bson.D{
						{Key: "format", Value: "%Y%m%d"},
						{Key: "date", Value: "$" + aggregateColumn},
					},
				},
			},
		}
	}

	// Grouping
	groupStage := bson.D{{Key: conditionGroup, Value: bson.D{
		aggregateOn, {
			Key: accumulationCountField,
			Value: bson.D{{
				Key:   conditionSum,
				Value: 1,
			}},
		},
	}}}

	pipeline := mongo.Pipeline{
		bson.D{
			{Key: conditionMatch, Value: matchStage},
		}, groupStage,
	}

	// anonymous struct for unmarshalling result bson
	var results []struct {
		ID    string `bson:"_id"`
		Count int64  `bson:"count"`
	}

	var aggregateCursor *mongo.Cursor
	aggregateCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Get the aggregation
	if aggregateCursor, err = collection.Aggregate(
		aggregateCtx, pipeline,
	); err != nil {
		return nil, err
	}

	// Cursor: All
	if err = aggregateCursor.All(ctx, &results); err != nil {
		return nil, err
	}

	// Create the result
	aggregateResult := make(map[string]interface{})
	for _, result := range results {
		aggregateResult[result.ID] = result.Count
	}

	return aggregateResult, nil
}

// GetMongoCollection will get the mongo collection for the given tableName
func (c *Client) GetMongoCollection(
	collectionName string,
) *mongo.Collection {
	return c.options.mongoDB.Collection(
		setPrefix(c.options.mongoDBConfig.TablePrefix, collectionName),
	)
}

// GetMongoCollectionByTableName will get the mongo collection for the given tableName
func (c *Client) GetMongoCollectionByTableName(
	tableName string,
) *mongo.Collection {
	return c.options.mongoDB.Collection(tableName)
}

// getFieldNames will get the field names in a slice of strings
func getFieldNames(fieldResult interface{}) []string {
	if fieldResult == nil {
		return []string{}
	}

	fields := make([]string, 0)

	model := reflect.ValueOf(fieldResult)
	if model.Kind() == reflect.Ptr {
		model = model.Elem()
	}
	if model.Kind() == reflect.Slice {
		elemType := model.Type().Elem()
		fmt.Println(elemType.Kind())
		if elemType.Kind() == reflect.Ptr {
			model = reflect.New(elemType.Elem())
		} else {
			model = reflect.New(elemType)
		}
	}
	if model.Kind() == reflect.Ptr {
		model = model.Elem()
	}

	for i := 0; i < model.Type().NumField(); i++ {
		field := model.Type().Field(i)
		fields = append(fields, field.Tag.Get(bsonTagName))
	}

	return fields
}

// setPrefix will automatically append the table prefix if found
func setPrefix(prefix, collection string) string {
	if len(prefix) > 0 {
		return prefix + "_" + collection
	}
	return collection
}

// getMongoQueryConditions will build the Mongo query conditions
// this functions tries to mimic the way gorm generates a where clause (naively)
func getMongoQueryConditions(
	model interface{},
	conditions map[string]interface{},
	customProcessor func(conditions *map[string]interface{}),
) map[string]interface{} {
	if conditions == nil {
		conditions = map[string]interface{}{}
	} else {
		// check for id field
		_, ok := conditions[sqlIDField]
		if ok {
			conditions[mongoIDField] = conditions[sqlIDField]
			delete(conditions, sqlIDField)
		}

		processMongoConditions(&conditions, customProcessor)
	}

	// add model ID to the query conditions, if set on the model
	id := GetModelStringAttribute(model, sqlIDFieldProper)
	if id != nil && *id != "" {
		conditions[mongoIDField] = *id
	}

	return conditions
}

// processMongoConditions will process all conditions for Mongo, including custom processing
func processMongoConditions(conditions *map[string]interface{},
	customProcessor func(conditions *map[string]interface{}),
) *map[string]interface{} {
	// Transform the id field to mongo _id field
	_, ok := (*conditions)[sqlIDField]
	if ok {
		(*conditions)[mongoIDField] = (*conditions)[sqlIDField]
		delete(*conditions, sqlIDField)
	}

	// Transform the map of metadata to key / value query
	_, ok = (*conditions)[metadataField]
	if ok {
		processMetadataConditions(conditions)
	}

	// Do we have a custom processor?
	if customProcessor != nil {
		customProcessor(conditions)
	}

	// Handle all conditions post-processing
	for key, condition := range *conditions {
		if key == conditionAnd || key == conditionOr {
			var slice []map[string]interface{}
			a, _ := json.Marshal(condition) //nolint:errchkjson // this check might break the current code
			_ = json.Unmarshal(a, &slice)
			var newConditions []map[string]interface{}
			for _, c := range slice {
				newConditions = append(newConditions, *processMongoConditions(&c, customProcessor)) //nolint:scopelint,gosec // ignore for now
			}
			(*conditions)[key] = newConditions
		}
	}

	return conditions
}

// processMetadataConditions will process metadata conditions
func processMetadataConditions(conditions *map[string]interface{}) {
	// marshal / unmarshal into standard map[string]interface{}
	m, _ := json.Marshal((*conditions)[metadataField]) //nolint:errchkjson // this check might break the current code
	var r map[string]interface{}
	_ = json.Unmarshal(m, &r)

	// Loop and create the key associations
	metadata := make([]map[string]interface{}, 0)
	for key, value := range r {
		metadata = append(metadata, map[string]interface{}{
			metadataField + ".k": key,
			metadataField + ".v": value,
		})
	}

	// Found some metadata
	if len(metadata) > 0 {
		_, ok := (*conditions)[conditionAnd]
		if ok {
			and := (*conditions)[conditionAnd].([]map[string]interface{})
			and = append(and, metadata...)
			(*conditions)[conditionAnd] = and
		} else {
			(*conditions)[conditionAnd] = metadata
		}
	}

	// Remove the field from conditions
	delete(*conditions, metadataField)
}

// openMongoDatabase will open a new database or use an existing connection
func openMongoDatabase(ctx context.Context, config *MongoDBConfig) (*mongo.Database, error) {
	// Use an existing connection
	if config.ExistingConnection != nil {
		return config.ExistingConnection, nil
	}

	// Create the new client
	nrMon := nrmongo.NewCommandMonitor(nil)
	client, err := mongo.Connect(
		ctx,
		options.Client().SetMonitor(nrMon),
		options.Client().ApplyURI(config.URI),
	)
	if err != nil {
		return nil, err
	}

	// Check the connection
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	// Return the client
	return client.Database(
		config.DatabaseName,
	), nil
}

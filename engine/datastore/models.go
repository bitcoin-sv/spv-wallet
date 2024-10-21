package datastore

import (
	"context"
	"errors"
	"fmt"
	"math"
	"reflect"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
)

// SaveModel will take care of creating or updating a model (primary key based) (abstracting the database)
//
// value is a pointer to the model
func (c *Client) SaveModel(
	_ context.Context,
	model interface{},
	tx *Transaction,
	newRecord, commitTx bool,
) error {
	if !IsSQLEngine(c.Engine()) {
		return ErrUnsupportedEngine
	}

	// Capture any panics
	defer func() {
		if r := recover(); r != nil {
			c.DebugLog(context.Background(), fmt.Sprintf("panic recovered: %v", r))
			_ = tx.Rollback()
		}
	}()
	if err := tx.sqlTx.Error; err != nil {
		return err
	}

	// Create vs Update
	if newRecord {
		if err := tx.sqlTx.Omit(clause.Associations).Create(model).Error; err != nil {
			_ = tx.Rollback()
			// todo add duplicate key check for Postgres and SQLite
			return err
		}
	} else {
		if err := tx.sqlTx.Omit(clause.Associations).Save(model).Error; err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	// Commit & check for errors
	if commitTx {
		if err := tx.Commit(); err != nil {
			return err
		}
	}

	// Return the tx
	return nil
}

// IncrementModel will increment the given field atomically in the database and return the new value
func (c *Client) IncrementModel(
	_ context.Context,
	model interface{},
	fieldName string,
	increment int64,
) (newValue int64, err error) {
	if !IsSQLEngine(c.Engine()) {
		return 0, ErrUnsupportedEngine
	}

	// Create a new transaction
	if err = c.options.db.Transaction(func(tx *gorm.DB) error {
		// Get the id of the model
		id := GetModelStringAttribute(model, sqlIDFieldProper)
		if id == nil {
			return spverrors.Newf("model is missing an %s field", sqlIDFieldProper)
		}

		// Get model if exist
		var result map[string]interface{}
		if err = tx.Model(&model).Clauses(clause.Locking{Strength: "UPDATE"}).Where(sqlIDField+" = ?", id).First(&result).Error; err != nil {
			return err
		}

		if result == nil {
			newValue = increment
			return nil
		}

		// Increment Counter
		newValue = convertToInt64(result[fieldName]) + increment
		return tx.Model(&model).Where(sqlIDField+" = ?", id).Update(fieldName, newValue).Error
	}); err != nil {
		return
	}

	return
}

// CreateInBatches create all the models given in batches
func (c *Client) CreateInBatches(
	_ context.Context,
	models interface{},
	batchSize int,
) error {
	tx := c.options.db.CreateInBatches(models, batchSize)
	return tx.Error
}

// convertToInt64 will convert an interface to an int64
func convertToInt64(i interface{}) int64 {
	switch v := i.(type) {
	case int:
		return int64(v)
	case int32:
		return int64(v)
	case uint32:
		return int64(v)
	case uint64:
		// Clamp values that are larger than MaxInt64
		if v > math.MaxInt64 {
			return math.MaxInt64
		}
		return int64(v)
	case int64:
		return v
	default:
		return 0
	}
}

// GetModel will get a model from the datastore
func (c *Client) GetModel(
	ctx context.Context,
	model interface{},
	conditions map[string]interface{},
	timeout time.Duration,
	forceWriteDB bool,
) error {
	if !IsSQLEngine(c.Engine()) {
		return ErrUnsupportedEngine
	}

	// Create a new context, and new db tx
	ctxDB, cancel := createCtx(ctx, c.options.db, timeout, c.IsDebug(), c.options.loggerDB)
	defer cancel()

	tx := ctxDB.Model(model)

	if forceWriteDB && c.Engine() == PostgreSQL {
		tx = ctxDB.Clauses(dbresolver.Write)
	}

	tx = tx.Select("*") // todo: optimize by specific fields

	if len(conditions) > 0 {
		var err error
		if tx, err = ApplyCustomWhere(c, tx, conditions, model); err != nil {
			return err
		}
	}

	return checkResult(tx.Find(model))
}

// GetModels will return a slice of models based on the given conditions
func (c *Client) GetModels(
	ctx context.Context,
	models interface{},
	conditions map[string]interface{},
	queryParams *QueryParams,
	fieldResults interface{},
	timeout time.Duration,
) error {
	if queryParams == nil {
		// init a new empty object for the default queryParams
		queryParams = &QueryParams{}
	}
	// Set default page size
	if queryParams.Page > 0 && queryParams.PageSize < 1 {
		queryParams.PageSize = defaultPageSize
	}

	// lower case the sort direction (asc / desc)
	queryParams.SortDirection = strings.ToLower(queryParams.SortDirection)

	// Switch on the datastore engines
	if !IsSQLEngine(c.Engine()) {
		return ErrUnsupportedEngine
	}
	return c.find(ctx, models, conditions, queryParams, fieldResults, timeout)
}

// GetModelCount will return a count of the model matching conditions
func (c *Client) GetModelCount(
	ctx context.Context,
	model interface{},
	conditions map[string]interface{},
	timeout time.Duration,
) (int64, error) {
	// Switch on the datastore engines
	if !IsSQLEngine(c.Engine()) {
		return 0, ErrUnsupportedEngine
	}

	return c.count(ctx, model, conditions, timeout)
}

// GetModelsAggregate will return an aggregate count of the model matching conditions
func (c *Client) GetModelsAggregate(ctx context.Context, models interface{},
	conditions map[string]interface{}, aggregateColumn string, timeout time.Duration,
) (map[string]interface{}, error) {
	// Switch on the datastore engines
	if !IsSQLEngine(c.Engine()) {
		return nil, ErrUnsupportedEngine
	}

	return c.aggregate(ctx, models, conditions, aggregateColumn, timeout)
}

// find will get records and return
func (c *Client) find(ctx context.Context, result interface{}, conditions map[string]interface{},
	queryParams *QueryParams, fieldResults interface{}, timeout time.Duration,
) error {
	// Find the type
	if reflect.TypeOf(result).Elem().Kind() != reflect.Slice {
		return spverrors.Newf("field: result is not a slice, found: %s", reflect.TypeOf(result).Kind().String())
	}

	// Create a new context, and new db tx
	ctxDB, cancel := createCtx(ctx, c.options.db, timeout, c.IsDebug(), c.options.loggerDB)
	defer cancel()

	tx := ctxDB.Model(result)

	// Create the offset
	offset := (queryParams.Page - 1) * queryParams.PageSize

	// Use the limit and offset
	if queryParams.Page > 0 && queryParams.PageSize > 0 {
		tx = tx.Limit(queryParams.PageSize).Offset(offset)
	}

	// Use an order field/sort
	if len(queryParams.OrderByField) > 0 {
		tx = tx.Order(clause.OrderByColumn{
			Column: clause.Column{
				Name: queryParams.OrderByField,
			},
			Desc: strings.ToLower(queryParams.SortDirection) == SortDesc,
		})
	}

	if len(conditions) > 0 {
		var err error
		if tx, err = ApplyCustomWhere(c, tx, conditions, result); err != nil {
			return err
		}
	}

	// Skip the conditions
	if fieldResults != nil {
		return checkResult(tx.Find(fieldResults))
	}
	return checkResult(tx.Find(result))
}

// find will get records and return
func (c *Client) count(ctx context.Context, model interface{}, conditions map[string]interface{},
	timeout time.Duration,
) (int64, error) {
	// Create a new context, and new db tx
	ctxDB, cancel := createCtx(ctx, c.options.db, timeout, c.IsDebug(), c.options.loggerDB)
	defer cancel()

	tx := ctxDB.Model(model)

	// Check for errors or no records found
	if len(conditions) > 0 {
		var err error
		if tx, err = ApplyCustomWhere(c, tx, conditions, model); err != nil {
			return 0, err
		}
	}
	var count int64
	err := checkResult(tx.Count(&count))

	return count, err
}

// find will get records and return
func (c *Client) aggregate(ctx context.Context, model interface{}, conditions map[string]interface{},
	aggregateColumn string, timeout time.Duration,
) (map[string]interface{}, error) {
	// Find the type
	if reflect.TypeOf(model).Elem().Kind() != reflect.Slice {
		return nil, spverrors.Newf("field: result is not a slice, found: %s", reflect.TypeOf(model).Kind().String())
	}

	// Create a new context, and new db tx
	ctxDB, cancel := createCtx(ctx, c.options.db, timeout, c.IsDebug(), c.options.loggerDB)
	defer cancel()

	// Get the tx
	tx := ctxDB.Model(model)

	// Check for errors or no records found
	var aggregate []map[string]interface{}
	if len(conditions) > 0 {
		var err error
		if tx, err = ApplyCustomWhere(c, tx, conditions, model); err != nil {
			return nil, err
		}
		err = checkResult(tx.Group(aggregateColumn).Scan(&aggregate))
		if err != nil {
			return nil, err
		}
	} else {
		aggregateCol := aggregateColumn

		// Check for a known date field
		if StringInSlice(aggregateCol, DateFields) {
			if c.Engine() == PostgreSQL {
				aggregateCol = "to_char(" + aggregateCol + ", 'YYYYMMDD')"
			} else {
				aggregateCol = "strftime('%Y%m%d', " + aggregateCol + ")"
			}
		}
		err := checkResult(tx.Select(aggregateCol + " as id, COUNT(id) AS count").Group(aggregateCol).Scan(&aggregate))
		if err != nil {
			return nil, err
		}
	}

	// Create the result
	aggregateResult := make(map[string]interface{})
	for _, item := range aggregate {
		key := item[sqlIDField].(string)
		aggregateResult[key] = item[accumulationCountField]
	}

	return aggregateResult, nil
}

// Execute a SQL query
func (c *Client) Execute(query string) *gorm.DB {
	if IsSQLEngine(c.Engine()) {
		return c.options.db.Exec(query)
	}

	return nil
}

// Raw a raw SQL query
func (c *Client) Raw(query string) *gorm.DB {
	if IsSQLEngine(c.Engine()) {
		return c.options.db.Raw(query)
	}

	return nil
}

// checkResult will check for records or error
func checkResult(result *gorm.DB) error {
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return ErrNoResults
		}
		return result.Error
	}

	// We should actually have some rows according to GORM
	if result.RowsAffected == 0 {
		return ErrNoResults
	}
	return nil
}

// createCtx will make a new DB context
func createCtx(ctx context.Context, db *gorm.DB, timeout time.Duration, debug bool,
	optionalLogger logger.Interface,
) (*gorm.DB, context.CancelFunc) {
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, timeout)
	return db.Session(getGormSessionConfig(db.PrepareStmt, debug, optionalLogger)).WithContext(ctx), cancel
}

package engine

import (
	"context"
	"errors"
	"time"

	"github.com/mrz1836/go-cachestore"
	"github.com/mrz1836/go-datastore"
)

// Get will retrieve a model from the Cachestore or Datastore using the provided conditions
//
// use bypassCache to skip checking the Cachestore for the record
func Get(
	ctx context.Context,
	model ModelInterface,
	conditions map[string]interface{},
	_ bool,
	timeout time.Duration,
	forceWriteDB bool,
) error {
	if timeout == 0 {
		timeout = defaultDatabaseReadTimeout
	}

	/*
		// todo: add cache support here for basic model lookups
	*/

	// Attempt to Get the model (by model fields & given conditions)
	return model.Client().Datastore().GetModel(ctx, model, conditions, timeout, forceWriteDB)
}

// getModels will retrieve model(s) from the Cachestore or Datastore using the provided conditions
//
// use bypassCache to skip checking the Cachestore for the record
func getModels(
	ctx context.Context,
	datastore datastore.ClientInterface,
	models interface{},
	conditions map[string]interface{},
	queryParams *datastore.QueryParams,
	timeout time.Duration,
) error {
	// Attempt to Get the model (by model fields & given conditions)
	return datastore.GetModels(ctx, models, conditions, queryParams, nil, timeout)
}

// getModelsAggregate will retrieve a count of the model(s) from the Cachestore or Datastore using the provided conditions
func getModelsAggregate(
	ctx context.Context,
	datastore datastore.ClientInterface,
	models interface{},
	conditions map[string]interface{},
	aggregateColumn string,
	timeout time.Duration,
) (map[string]interface{}, error) {
	// Attempt to Get the model (by model fields & given conditions)
	return datastore.GetModelsAggregate(ctx, models, conditions, aggregateColumn, timeout)
}

// getModelCount will retrieve a count of the model from the Cachestore or Datastore using the provided conditions
func getModelCount(
	ctx context.Context,
	datastore datastore.ClientInterface,
	model interface{},
	conditions map[string]interface{},
	timeout time.Duration, //nolint:nolintlint,unparam // default timeout is passed most of the time
) (int64, error) {
	// Attempt to Get the model (by model fields & given conditions)
	return datastore.GetModelCount(ctx, model, conditions, timeout)
}

// getModelsByConditions will get models by given conditions
func getModelsByConditions(ctx context.Context, modelName ModelName, modelItems interface{},
	metadata *Metadata, conditions *map[string]interface{}, queryParams *datastore.QueryParams,
	opts ...ModelOps,
) error {
	dbConditions := map[string]interface{}{}

	if metadata != nil {
		dbConditions[metadataField] = metadata
	}

	if conditions != nil && len(*conditions) > 0 {
		and := make([]map[string]interface{}, 0)
		if _, ok := dbConditions["$and"]; ok {
			and = dbConditions["$and"].([]map[string]interface{})
		}
		and = append(and, *conditions)
		dbConditions["$and"] = and
	}

	// Get the records
	if err := getModels(
		ctx, NewBaseModel(modelName, opts...).Client().Datastore(),
		modelItems, dbConditions, queryParams, defaultDatabaseReadTimeout,
	); err != nil {
		if errors.Is(err, datastore.ErrNoResults) {
			return nil
		}
		return err
	}

	return nil
}

// getModelsAggregateByConditions will get aggregates of models by given conditions
func getModelsAggregateByConditions(ctx context.Context, modelName ModelName, models interface{},
	metadata *Metadata, conditions *map[string]interface{}, aggregateColumn string,
	opts ...ModelOps,
) (map[string]interface{}, error) {
	dbConditions := map[string]interface{}{}

	if metadata != nil {
		dbConditions[metadataField] = metadata
	}

	if conditions != nil && len(*conditions) > 0 {
		and := make([]map[string]interface{}, 0)
		if _, ok := dbConditions["$and"]; ok {
			and = dbConditions["$and"].([]map[string]interface{})
		}
		and = append(and, *conditions)
		dbConditions["$and"] = and
	}

	// Get the records
	results, err := getModelsAggregate(
		ctx, NewBaseModel(modelName, opts...).Client().Datastore(),
		models, dbConditions, aggregateColumn, defaultDatabaseReadTimeout,
	)
	if err != nil {
		if errors.Is(err, datastore.ErrNoResults) {
			return nil, nil
		}
		return nil, err
	}

	return results, nil
}

// getModelCountByConditions will get model counts (sums) from given conditions
func getModelCountByConditions(ctx context.Context, modelName ModelName, model interface{},
	metadata *Metadata, conditions *map[string]interface{}, opts ...ModelOps,
) (int64, error) {
	dbConditions := map[string]interface{}{}

	if metadata != nil {
		dbConditions[metadataField] = metadata
	}

	if conditions != nil && len(*conditions) > 0 {
		and := make([]map[string]interface{}, 0)
		if _, ok := dbConditions["$and"]; ok {
			and = dbConditions["$and"].([]map[string]interface{})
		}
		and = append(and, *conditions)
		dbConditions["$and"] = and
	}

	// Get the records
	count, err := getModelCount(
		ctx, NewBaseModel(modelName, opts...).Client().Datastore(),
		model, dbConditions, defaultDatabaseReadTimeout,
	)
	if err != nil {
		if errors.Is(err, datastore.ErrNoResults) {
			return 0, nil
		}
		return 0, err
	}

	return count, nil
}

// getModelFromCache will attempt to get a model from cache
func getModelFromCache(ctx context.Context, cacheClient cachestore.ClientInterface,
	key string, model ModelInterface,
) (bool, error) { // Success if the key was found
	if err := cacheClient.GetModel(ctx, key, model); err != nil {
		if errors.Is(err, cachestore.ErrKeyNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

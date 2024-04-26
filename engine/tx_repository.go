package engine

import (
	"context"
	"errors"

	"github.com/bitcoin-sv/spv-wallet/engine/datastore"
	"github.com/libsv/go-bt"
)

// getTransactionByID will get the model from a given transaction ID
func getTransactionByID(ctx context.Context, xPubID, txID string, opts ...ModelOps) (*Transaction, error) {
	// Construct an empty tx
	tx := emptyTx(opts...)
	tx.ID = txID
	tx.XPubID = xPubID

	// Get the record
	if err := Get(ctx, tx, nil, false, defaultDatabaseReadTimeout, false); err != nil {
		if errors.Is(err, datastore.ErrNoResults) {
			return nil, nil
		}
		return nil, err
	}

	return tx, nil
}

// getTransactions will get all the transactions with the given conditions
func getTransactions(ctx context.Context, metadata *Metadata, conditions map[string]interface{},
	queryParams *datastore.QueryParams, opts ...ModelOps,
) ([]*Transaction, error) {
	modelItems := make([]*Transaction, 0)
	if err := getModelsByConditions(ctx, ModelTransaction, &modelItems, metadata, conditions, queryParams, opts...); err != nil {
		return nil, err
	}

	return modelItems, nil
}

// getTransactionsAggregate will get a count of all transactions per aggregate column with the given conditions
func getTransactionsAggregate(ctx context.Context, metadata *Metadata, conditions map[string]interface{},
	aggregateColumn string, opts ...ModelOps,
) (map[string]interface{}, error) {
	modelItems := make([]*Transaction, 0)
	results, err := getModelsAggregateByConditions(
		ctx, ModelTransaction, &modelItems, metadata, conditions, aggregateColumn, opts...,
	)
	if err != nil {
		return nil, err
	}

	return results, nil
}

// getTransactionsCount will get a count of all the transactions with the given conditions
func getTransactionsCount(ctx context.Context, metadata *Metadata, conditions map[string]interface{},
	opts ...ModelOps,
) (int64, error) {
	return getModelCountByConditions(ctx, ModelTransaction, Transaction{}, metadata, conditions, opts...)
}

// getTransactionsCountByXpubID will get the count of all the models for a given xpub ID
func getTransactionsCountByXpubID(ctx context.Context, xPubID string, metadata *Metadata,
	conditions map[string]interface{}, opts ...ModelOps,
) (int64, error) {
	dbConditions := processDBConditions(xPubID, conditions, metadata)

	return getTransactionsCountInternal(ctx, dbConditions, opts...)
}

// getTransactionsByXpubID will get all the models for a given xpub ID
func getTransactionsByXpubID(ctx context.Context, xPubID string,
	metadata *Metadata, conditions map[string]interface{},
	queryParams *datastore.QueryParams, opts ...ModelOps,
) ([]*Transaction, error) {
	dbConditions := processDBConditions(xPubID, conditions, metadata)

	return getTransactionsInternal(ctx, dbConditions, xPubID, queryParams, opts...)
}

func processDBConditions(xPubID string, conditions map[string]interface{},
	metadata *Metadata,
) map[string]interface{} {
	dbConditions := map[string]interface{}{
		"$or": []map[string]interface{}{{
			"xpub_in_ids": xPubID,
		}, {
			"xpub_out_ids": xPubID,
		}},
	}

	if metadata != nil && len(*metadata) > 0 {
		and := make([]map[string]interface{}, 0)
		if _, ok := dbConditions["$and"]; ok {
			and = dbConditions["$and"].([]map[string]interface{})
		}
		for key, value := range *metadata {
			condition := map[string]interface{}{
				"$or": []map[string]interface{}{{
					metadataField: Metadata{
						key: value,
					},
				}, {
					xPubMetadataField: XpubMetadata{
						xPubID: Metadata{
							key: value,
						},
					},
				}},
			}
			and = append(and, condition)
		}
		dbConditions["$and"] = and
	}

	if len(conditions) > 0 {
		and := make([]map[string]interface{}, 0)
		if _, ok := dbConditions["$and"]; ok {
			and = dbConditions["$and"].([]map[string]interface{})
		}
		and = append(and, conditions)
		dbConditions["$and"] = and
	}

	return dbConditions
}

// getTransactionsInternal get all transactions for the given conditions
// NOTE: this function should only be used internally, it allows to query the whole transaction table
func getTransactionsInternal(ctx context.Context, conditions map[string]interface{}, xPubID string,
	queryParams *datastore.QueryParams, opts ...ModelOps,
) ([]*Transaction, error) {
	var models []Transaction
	if err := getModels(
		ctx,
		NewBaseModel(ModelTransaction, opts...).Client().Datastore(),
		&models,
		conditions,
		queryParams,
		defaultDatabaseReadTimeout,
	); err != nil {
		if errors.Is(err, datastore.ErrNoResults) {
			return nil, nil
		}
		return nil, err
	}

	// Loop and enrich
	transactions := make([]*Transaction, 0)
	for index := range models {
		models[index].enrich(ModelTransaction, opts...)
		models[index].XPubID = xPubID
		tx := &models[index]
		transactions = append(transactions, tx)
	}

	return transactions, nil
}

// getTransactionsCountInternal get a count of all transactions for the given conditions
func getTransactionsCountInternal(ctx context.Context, conditions map[string]interface{},
	opts ...ModelOps,
) (int64, error) {
	count, err := getModelCount(
		ctx,
		NewBaseModel(ModelNameEmpty, opts...).Client().Datastore(),
		Transaction{},
		conditions,
		defaultDatabaseReadTimeout,
	)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func getTransactionByHex(ctx context.Context, hex string, opts ...ModelOps) (*Transaction, error) {
	btTx, err := bt.NewTxFromString(hex)
	if err != nil {
		return nil, err
	}

	return getTransactionByID(ctx, "", btTx.GetTxID(), opts...)
}

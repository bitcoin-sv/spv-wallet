package engine

import (
	"context"
	"errors"

	"github.com/bitcoin-sv/spv-wallet/engine/datastore"
)

/*** exported funcs ***/

// GetSyncTransactionByID will get a sync transaction
func GetSyncTransactionByID(ctx context.Context, id string, opts ...ModelOps) (*SyncTransaction, error) {
	// Get the records by status
	txs, err := _getSyncTransactionsByConditions(ctx,
		map[string]interface{}{
			idField: id,
		},
		nil, opts...,
	)
	if err != nil {
		return nil, err
	}
	if len(txs) != 1 {
		return nil, nil
	}

	return txs[0], nil
}

// GetSyncTransactionByTxID will get a sync transaction by it's transaction id.
func GetSyncTransactionByTxID(ctx context.Context, txID string, opts ...ModelOps) (*SyncTransaction, error) {
	// Get the records by status
	txs, err := _getSyncTransactionsByConditions(ctx,
		map[string]interface{}{
			idField: txID,
		},
		nil, opts...,
	)
	if err != nil {
		return nil, err
	}
	if len(txs) != 1 {
		return nil, nil
	}

	return txs[0], nil
}

/*** /exported funcs ***/

/*** public unexported funcs ***/

// getTransactionsToSync will get the sync transactions to sync
func getTransactionsToSync(ctx context.Context, queryParams *datastore.QueryParams,
	opts ...ModelOps,
) ([]*SyncTransaction, error) {
	// Get the records by status
	txs, err := _getSyncTransactionsByConditions(
		ctx,
		map[string]interface{}{
			syncStatusField: SyncStatusReady.String(),
		},
		queryParams, opts...,
	)
	if err != nil {
		return nil, err
	}
	return txs, nil
}

/*** /public unexported funcs ***/

// getTransactionsToSync will get the sync transactions to sync
func _getSyncTransactionsByConditions(ctx context.Context, conditions map[string]interface{},
	queryParams *datastore.QueryParams, opts ...ModelOps,
) ([]*SyncTransaction, error) {
	if queryParams == nil {
		queryParams = &datastore.QueryParams{
			OrderByField:  createdAtField,
			SortDirection: datastore.SortAsc,
		}
	} else if queryParams.OrderByField == "" || queryParams.SortDirection == "" {
		queryParams.OrderByField = createdAtField
		queryParams.SortDirection = datastore.SortAsc
	}

	// Get the records
	var models []SyncTransaction
	if err := getModels(
		ctx, NewBaseModel(ModelNameEmpty, opts...).Client().Datastore(),
		&models, conditions, queryParams, defaultDatabaseReadTimeout,
	); err != nil {
		if errors.Is(err, datastore.ErrNoResults) {
			return nil, nil
		}
		return nil, err
	}

	// Loop and enrich
	txs := make([]*SyncTransaction, 0)
	for index := range models {
		models[index].enrich(ModelSyncTransaction, opts...)
		txs = append(txs, &models[index])
	}

	return txs, nil
}

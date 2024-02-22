package engine

import (
	"context"
	"encoding/hex"
	"errors"

	"github.com/libsv/go-bt/v2"
	"github.com/mrz1836/go-datastore"
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

// getTransactionsToBroadcast will get the sync transactions to broadcast
func getTransactionsToBroadcast(ctx context.Context, queryParams *datastore.QueryParams,
	opts ...ModelOps,
) ([]*SyncTransaction, error) {
	// Get the records by status
	scTxs, err := _getSyncTransactionsByConditions(
		ctx,
		map[string]interface{}{
			broadcastStatusField: SyncStatusReady.String(),
		},
		queryParams, opts...,
	)
	if err != nil {
		return nil, err
	} else if len(scTxs) == 0 {
		return nil, nil
	}

	// hydrate and see if it's ready to sync
	res := make([]*SyncTransaction, 0, len(scTxs))

	for _, sTx := range scTxs {
		// hydrate
		sTx.transaction, err = getTransactionByID(
			ctx, "", sTx.ID, opts...,
		)
		if err != nil {
			return nil, err
		} else if sTx.transaction == nil {
			return nil, ErrMissingTransaction
		}

		parentsBroadcast, err := _areParentsBroadcasted(ctx, sTx.transaction, opts...)
		if err != nil {
			return nil, err
		}

		if !parentsBroadcast {
			// if all parents are not broadcast, then we cannot broadcast this tx
			continue
		}

		res = append(res, sTx)
	}

	return res, nil
}

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

func _areParentsBroadcasted(ctx context.Context, tx *Transaction, opts ...ModelOps) (bool, error) {
	// get the sync transaction of all inputs
	btTx, err := bt.NewTxFromString(tx.Hex)
	if err != nil {
		return false, err
	}

	// check that all inputs we handled have been broadcast, or are not handled by SPV Wallet Engine
	parentsBroadcasted := true
	for _, input := range btTx.Inputs {
		var parentTx *SyncTransaction
		previousTxID := hex.EncodeToString(bt.ReverseBytes(input.PreviousTxID()))
		parentTx, err = GetSyncTransactionByID(ctx, previousTxID, opts...)
		if err != nil {
			return false, err
		}
		// if we have a sync transaction, and it is not complete, then we cannot broadcast
		if parentTx != nil && parentTx.BroadcastStatus != SyncStatusComplete {
			parentsBroadcasted = false
		}
	}

	return parentsBroadcasted, nil
}

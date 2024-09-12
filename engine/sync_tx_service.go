package engine

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/bitcoin-sv/spv-wallet/engine/chainstate"
	"github.com/bitcoin-sv/spv-wallet/engine/datastore"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
)

// processSyncTransactions will process sync transaction records
func processSyncTransactions(ctx context.Context, maxTransactions int, opts ...ModelOps) error {
	queryParams := &datastore.QueryParams{
		Page:          1,
		PageSize:      maxTransactions,
		OrderByField:  "created_at",
		SortDirection: "desc",
	}

	// Get x records
	records, err := getTransactionsToSync(
		ctx, queryParams, opts...,
	)
	if err != nil {
		return err
	} else if len(records) == 0 {
		return nil
	}

	for index := range records {
		if err = _syncTxDataFromChain(
			ctx, records[index], nil,
		); err != nil {
			return err
		}
	}

	return nil
}

// broadcastTxAndUpdateSync will broadcast transaction and and SyncStatus in syncTx
// It most probably will be deleted after syncTX removal
func broadcastTxAndUpdateSync(ctx context.Context, tx *Transaction) error {
	syncTx := tx.syncTransaction
	err := broadcastTransaction(ctx, tx)
	if err != nil {
		return err
	}

	// Update sync status to be ready now
	if syncTx.SyncStatus == SyncStatusPending {
		syncTx.SyncStatus = SyncStatusReady
	}

	return syncTx.Save(ctx)
}

func broadcastTransaction(ctx context.Context, tx *Transaction) error {
	client := tx.Client()
	chainstateSrv := client.Chainstate()

	// Successfully capture any panics, convert to readable string and log the error
	defer recoverAndLog(tx.Client().Logger())

	// Create the lock and set the release for after the function completes
	unlock, err := newWriteLock(
		ctx, fmt.Sprintf(lockKeyProcessBroadcastTx, tx.GetID()), client.Cachestore(),
	)
	defer unlock()
	if err != nil {
		return err
	}

	// Broadcast
	txHex, hexFormat := _getTxHexInFormat(ctx, tx, chainstateSrv.SupportedBroadcastFormats(), client)
	br := chainstateSrv.Broadcast(ctx, tx.ID, txHex, hexFormat, defaultBroadcastTimeout)

	if br.Failure != nil { // broadcast failed
		return br.Failure.Error
	}

	return nil
}

// ///////////////

func _getTxHexInFormat(ctx context.Context, tx *Transaction, prefferedFormat chainstate.HexFormatFlag, store TransactionGetter) (txHex string, actualFormat chainstate.HexFormatFlag) {
	if prefferedFormat.Contains(chainstate.Ef) {
		efHex, ok := ToEfHex(ctx, tx, store)

		if ok {
			txHex = efHex
			actualFormat = chainstate.Ef
			return
		}
	}

	// return rawtx hex
	txHex = tx.Hex
	actualFormat = chainstate.RawTx

	return
}

// _syncTxDataFromChain will process the sync transaction record, or save the failure
func _syncTxDataFromChain(ctx context.Context, syncTx *SyncTransaction, transaction *Transaction) error {
	// Successfully capture any panics, convert to readable string and log the error
	defer recoverAndLog(syncTx.Client().Logger())

	var err error

	if transaction == nil {
		if transaction, err = _getTransaction(ctx, syncTx.ID, syncTx.GetOptions(false)); err != nil {
			return spverrors.ErrCouldNotFindTransaction
		}
	}

	// Find on-chain
	var txInfo *chainstate.TransactionInfo
	if txInfo, err = syncTx.Client().Chainstate().QueryTransaction(
		ctx, syncTx.ID, chainstate.RequiredOnChain, defaultQueryTxTimeout,
	); err != nil {
		if errors.Is(err, spverrors.ErrCouldNotFindTransaction) {
			syncTx.Client().Logger().Info().
				Str("txID", syncTx.ID).
				Msgf("Transaction not found on-chain, will try again later")

			syncTx.SyncStatus = SyncStatusReady
			return nil
		}
		return spverrors.Wrapf(err, "could not query transaction")
	}
	return processSyncTxSave(ctx, txInfo, syncTx, transaction)
}

func _getTransaction(ctx context.Context, id string, opts []ModelOps) (*Transaction, error) {
	transaction, err := getTransactionByID(ctx, "", id, opts...)
	if err != nil {
		return nil, err
	}

	if transaction == nil {
		return nil, spverrors.ErrCouldNotFindTransaction
	}

	return transaction, nil
}

func processSyncTxSave(ctx context.Context, txInfo *chainstate.TransactionInfo, syncTx *SyncTransaction, transaction *Transaction) error {
	if !txInfo.Valid() {
		syncTx.Client().Logger().Warn().
			Str("txID", syncTx.ID).
			Msgf("txInfo is invalid, will try again later")

		if syncTx.Client().IsDebug() {
			txInfoJSON, _ := json.Marshal(txInfo)
			syncTx.Client().Logger().Debug().
				Str("txID", syncTx.ID).
				Msgf("txInfo: %s", string(txInfoJSON))
		}
		return nil
	}

	transaction.setChainInfo(txInfo)
	if err := transaction.Save(ctx); err != nil {
		return err
	}

	syncTx.SyncStatus = SyncStatusComplete

	if err := syncTx.Save(ctx); err != nil {
		return err
	}

	syncTx.Client().Logger().Info().
		Str("txID", syncTx.ID).
		Msgf("Transaction processed successfully")
	return nil
}

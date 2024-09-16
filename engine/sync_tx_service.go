package engine

import (
	"context"
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
		if err = _syncTxDataFromChain(ctx, records[index]); err != nil {
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
func _syncTxDataFromChain(ctx context.Context, syncTx *SyncTransaction) error {
	logger := syncTx.Client().Logger()
	defer recoverAndLog(logger)

	tx, err := _getTransaction(ctx, syncTx.ID, syncTx.GetOptions(false))
	if err != nil {
		return spverrors.ErrCouldNotFindTransaction
	}

	chainstateService := syncTx.Client().Chainstate()
	txInfo, err := chainstateService.QueryTransaction(
		ctx, syncTx.ID, chainstate.RequiredOnChain, defaultQueryTxTimeout,
	)
	if err != nil {
		if errors.Is(err, spverrors.ErrCouldNotFindTransaction) {
			/* DEPRECATED block of code with syncTx - will be removed soon */
			syncTx.SyncStatus = SyncStatusReady
			/* DEPRECATED END */
			return nil
		}
		return spverrors.Wrapf(err, "could not query transaction")
	}

	tx.BlockHash = txInfo.BlockHash
	tx.BlockHeight = uint64(txInfo.BlockHeight)
	tx.TxStatus = string(txInfo.TxStatus)
	tx.SetBUMP(txInfo.BUMP)

	if !tx.IsOnChain() {
		logger.Warn().Interface("txInfo", txInfo).Msgf("TransactionInfo is invalid")
		return nil
	}
	if err := tx.Save(ctx); err != nil {
		return err
	}

	/* DEPRECATED block of code with syncTx - will be removed soon */
	syncTx.SyncStatus = SyncStatusComplete
	if err := syncTx.Save(ctx); err != nil {
		return err
	}
	/* DEPRECATED END */

	return nil
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

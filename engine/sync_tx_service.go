package engine

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/bitcoin-sv/spv-wallet/engine/chainstate"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
)

// processSyncTransactions will process sync transaction records
func processSyncTransactions(ctx context.Context, client *Client) {
	logger := client.Logger()
	db := client.Datastore().DB()
	chainstateService := client.Chainstate()

	delayedForBroadcasted := time.Now().Add(-1 * time.Hour)
	delayedForNotBroadcasted := time.Now().Add(-10 * time.Minute)
	var txIDsToSync []struct {
		ID string
	}
	err := db.
		Model(&Transaction{}).
		Where("state = ? AND created_at < ?", TxStatusBroadcasted, delayedForBroadcasted).
		Or("state IN (?) AND created_at < ?", []TxStatus{TxStatusCreated, TxStatusSent}, delayedForNotBroadcasted).
		Find(&txIDsToSync).Error
	if err != nil {
		logger.Error().Err(err).Msg("Cannot fetch transactions to sync")
		return
	}

	var tx Transaction
	for _, record := range txIDsToSync {
		txID := record.ID
		err := db.First(&tx).Error
		if err != nil {
			logger.Warn().Str("txID", txID).Msg("Cannot get transaction by ID even though the ID was returned from DB")
			continue
		}
		tx.client = client // hydrate the client
		txInfo, err := chainstateService.QueryTransaction(ctx, txID, chainstate.RequiredOnChain, defaultQueryTxTimeout)
		if err != nil {
			if errors.Is(err, spverrors.ErrBroadcastUnreachable) {
				logger.Warn().Msgf("%s", err.Error())
				// checking subsequent transactions is pointless if the broadcast server (ARC) is unreachable
				// will try again in the next cycle
				return
			}

			if errors.Is(err, spverrors.ErrBroadcastRejectedTransaction) {
				tx.TxStatus = string(TxStatusProblematic)
				if err := tx.Save(ctx); err != nil {
					logger.Error().Err(err).Str("txID", txID).Msg("Cannot update transaction status to problematic")
					continue
				}
			}

			if errors.Is(err, spverrors.ErrCouldNotFindTransaction) {
				if tx.CreatedAt.Before(time.Now().Add(-24 * time.Hour)) {
					tx.TxStatus = string(TxStatusProblematic)
					if err := tx.Save(ctx); err != nil {
						logger.Error().Err(err).Str("txID", txID).Msg("Cannot update transaction status to problematic")
						continue
					}
				}

				if tx.TxStatus == string(TxStatusCreated) || tx.TxStatus == string(TxStatusSent) {
					// TODO try to broadcast again
					continue
				}
			}

			logger.Error().Err(err).Str("txID", txID).Msg("Cannot query transaction")
		}

		tx.BlockHash = txInfo.BlockHash
		tx.BlockHeight = uint64(txInfo.BlockHeight)
		tx.SetBUMP(txInfo.BUMP)
		tx.UpdateFromBroadcastStatus(txInfo.TxStatus)

		tx.TxStatus = string(TxStatusProblematic)
		if err := tx.Save(ctx); err != nil {
			logger.Error().Err(err).Str("txID", txID).Msg("Cannot update transaction status to problematic")
			continue
		}

	}
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

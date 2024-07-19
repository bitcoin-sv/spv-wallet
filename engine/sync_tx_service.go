package engine

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/bitcoin-sv/go-paymail"
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

// processBroadcastTransactions will process sync transaction records
func processBroadcastTransactions(ctx context.Context, maxTransactions int, opts ...ModelOps) error {
	queryParams := &datastore.QueryParams{
		Page:          1,
		PageSize:      maxTransactions,
		OrderByField:  createdAtField,
		SortDirection: datastore.SortAsc,
	}

	// Get maxTransactions records, grouped by xpub
	snTxs, err := getTransactionsToBroadcast(ctx, queryParams, opts...)
	if err != nil {
		return err
	} else if len(snTxs) == 0 {
		return nil
	}

	// Process the transactions per xpub, in parallel
	txsByXpub := _groupByXpub(snTxs)

	// we limit the number of concurrent broadcasts to the number of cpus*2, since there is lots of IO wait
	limit := make(chan bool, runtime.NumCPU()*2)
	wg := new(sync.WaitGroup)

	for xPubID := range txsByXpub {
		limit <- true // limit the number of routines running at the same time
		wg.Add(1)
		go func(xPubID string) {
			defer wg.Done()
			defer func() { <-limit }()

			for _, tx := range txsByXpub[xPubID] {
				if err = broadcastSyncTransaction(
					ctx, tx,
				); err != nil {
					tx.Client().Logger().Error().
						Str("txID", tx.ID).
						Str("xpubID", xPubID).
						Msgf("error running broadcast tx: %s", err.Error())
					return // stop processing transactions for this xpub if we found an error
				}
			}
		}(xPubID)
	}
	wg.Wait()

	return nil
}

// broadcastSyncTransaction will broadcast transaction related to syncTx record
func broadcastSyncTransaction(ctx context.Context, syncTx *SyncTransaction) error {
	// Successfully capture any panics, convert to readable string and log the error
	defer recoverAndLog(syncTx.Client().Logger())

	// Create the lock and set the release for after the function completes
	unlock, err := newWriteLock(
		ctx, fmt.Sprintf(lockKeyProcessBroadcastTx, syncTx.GetID()), syncTx.Client().Cachestore(),
	)
	defer unlock()
	if err != nil {
		return err
	}

	client := syncTx.Client()
	chainstateSrv := client.Chainstate()

	// Get the transaction HEX
	tx := syncTx.transaction
	if tx == nil || tx.Hex == "" {
		if tx, err = _getTransaction(ctx, syncTx.ID, syncTx.GetOptions(false)); err != nil {
			return nil
		}
	}

	// Broadcast
	txHex, hexFormat := _getTxHexInFormat(ctx, tx, chainstateSrv.SupportedBroadcastFormats(), client)
	br := chainstateSrv.Broadcast(ctx, syncTx.ID, txHex, hexFormat, defaultBroadcastTimeout)

	if br.Failure != nil { // broadcast failed
		if br.Failure.InvalidTx {
			syncTx.BroadcastStatus = SyncStatusError // invalid transaction, won't be broadcasted anymore
		} else {
			syncTx.BroadcastStatus = SyncStatusReady // client error, try again later
		}

		_addSyncResult(ctx, syncTx, syncActionBroadcast, br.Provider, br.Failure.Error.Error())
		return br.Failure.Error
	}

	// Update the sync information
	syncTx.BroadcastStatus = SyncStatusComplete
	// Update sync status to be ready now
	if syncTx.SyncStatus == SyncStatusPending {
		syncTx.SyncStatus = SyncStatusReady
	}

	syncTx.Results.Results = append(syncTx.Results.Results, &SyncResult{
		Action:        syncActionBroadcast,
		ExecutedAt:    time.Now().UTC(),
		Provider:      br.Provider,
		StatusMessage: "broadcast success",
	})

	// Update the sync transaction record
	if err = syncTx.Save(ctx); err != nil {
		_addSyncResult(ctx, syncTx, syncActionBroadcast, "internal", err.Error())
		return err
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
			_addSyncResult(ctx, syncTx, syncActionSync, "all", "transaction not found on-chain")
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
		_addSyncResult(ctx, syncTx, syncActionSync, "internal", err.Error())
		return err
	}

	syncTx.SyncStatus = SyncStatusComplete
	syncTx.Results.Results = append(syncTx.Results.Results, &SyncResult{
		Action:        syncActionSync,
		ExecutedAt:    time.Now().UTC(),
		Provider:      chainstate.ProviderBroadcastClient,
		StatusMessage: "transaction was found on-chain by " + chainstate.ProviderBroadcastClient,
	})

	if err := syncTx.Save(ctx); err != nil {
		_addSyncResult(ctx, syncTx, syncActionSync, "internal", err.Error())
		return err
	}

	syncTx.Client().Logger().Info().
		Str("txID", syncTx.ID).
		Msgf("Transaction processed successfully")
	return nil
}

// processP2PTransaction will process the sync transaction record, or save the failure
func processP2PTransaction(ctx context.Context, tx *Transaction) error {
	// Successfully capture any panics, convert to readable string and log the error
	defer recoverAndLog(tx.Client().Logger())

	syncTx := tx.syncTransaction
	// Create the lock and set the release for after the function completes
	unlock, err := newWriteLock(
		ctx, fmt.Sprintf(lockKeyProcessP2PTx, syncTx.GetID()), syncTx.Client().Cachestore(),
	)
	defer unlock()
	if err != nil {
		return err
	}

	// No draft?
	if len(tx.DraftID) == 0 {
		syncTx.P2PStatus = SyncStatusError
		_addSyncResult(ctx, syncTx, syncActionP2P, "all", "no draft found, cannot complete p2p")

		return nil
	}

	// Notify any P2P paymail providers associated to the transaction
	var results []*SyncResult
	if results, err = _notifyPaymailProviders(ctx, tx); err != nil {
		syncTx.P2PStatus = SyncStatusReady
		_addSyncResult(ctx, syncTx, syncActionP2P, "", err.Error())
		return err
	}

	// Update if we have some results
	if len(results) > 0 {
		syncTx.Results.Results = append(syncTx.Results.Results, results...)
	}

	// Save the record
	syncTx.P2PStatus = SyncStatusComplete

	// Update sync status to be ready now
	if syncTx.SyncStatus == SyncStatusPending {
		syncTx.SyncStatus = SyncStatusReady
	}

	if err = syncTx.Save(ctx); err != nil {
		syncTx.P2PStatus = SyncStatusError
		_addSyncResult(ctx, syncTx, syncActionP2P, "internal", err.Error())
		return err
	}

	// Done!
	return nil
}

// _notifyPaymailProviders will notify any associated Paymail providers
func _notifyPaymailProviders(ctx context.Context, transaction *Transaction) ([]*SyncResult, error) {
	pm := transaction.Client().PaymailClient()
	outputs := transaction.draftTransaction.Configuration.Outputs

	notifiedReceivers := make([]string, 0)
	results := make([]*SyncResult, len(outputs))

	var payload *paymail.P2PTransactionPayload
	var err error

	for _, out := range outputs {
		p4 := out.PaymailP4

		if p4 == nil || p4.ResolutionType != ResolutionTypeP2P {
			continue
		}

		receiver := fmt.Sprintf("%s@%s", p4.Alias, p4.Domain)
		if contains(notifiedReceivers, func(x string) bool { return x == receiver }) {
			continue // no need to send the same transaction to the same receiver second time
		}

		if payload, err = finalizeP2PTransaction(
			ctx,
			pm,
			p4,
			transaction,
		); err != nil {
			return nil, err
		}

		notifiedReceivers = append(notifiedReceivers, receiver)
		results = append(results, &SyncResult{
			Action:        syncActionP2P,
			ExecutedAt:    time.Now().UTC(),
			Provider:      p4.ReceiveEndpoint,
			StatusMessage: "success: " + payload.TxID,
		})

	}
	return results, nil
}

// utils

func _groupByXpub(scTxs []*SyncTransaction) map[string][]*SyncTransaction {
	txsByXpub := make(map[string][]*SyncTransaction)

	// group transactions by xpub and return including the tx itself
	for _, tx := range scTxs {
		xPubID := "" // fallback if we have no input xpubs
		if len(tx.transaction.XpubInIDs) > 0 {
			// use the first xpub for the grouping
			// in most cases when we are broadcasting, there should be only 1 xpub in
			xPubID = tx.transaction.XpubInIDs[0]
		}

		if txsByXpub[xPubID] == nil {
			txsByXpub[xPubID] = make([]*SyncTransaction, 0)
		}
		txsByXpub[xPubID] = append(txsByXpub[xPubID], tx)
	}

	return txsByXpub
}

// _addSyncResult will save the error message for a sync tx
func _addSyncResult(ctx context.Context, syncTx *SyncTransaction,
	action, provider, message string,
) {
	syncTx.Results.Results = append(syncTx.Results.Results, &SyncResult{
		Action:        action,
		ExecutedAt:    time.Now().UTC(),
		Provider:      provider,
		StatusMessage: message,
	})

	if syncTx.IsNew() {
		return // do not save if new record! caller should decide if want to save new record
	}

	_ = syncTx.Save(ctx)
}

package engine

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/bitcoin-sv/go-paymail"
	"github.com/bitcoin-sv/spv-wallet/engine/chainstate"
	"github.com/bitcoin-sv/spv-wallet/engine/datastore"
	customTypes "github.com/bitcoin-sv/spv-wallet/engine/datastore/customtypes"
	"github.com/bitcoin-sv/spv-wallet/engine/notifications"
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

	txHex, hexFormat := _getTxHexInFormat(ctx, tx, chainstateSrv.SupportedBroadcastFormats(), client)

	// Broadcast
	var provider string
	if provider, err = chainstateSrv.Broadcast(
		ctx, syncTx.ID, txHex, hexFormat, defaultBroadcastTimeout,
	); err != nil {
		_bailAndSaveSyncTransaction(ctx, syncTx, SyncStatusReady, syncActionBroadcast, provider, err.Error())
		return err
	}

	// Update the sync information
	statusMsg := "broadcast success"

	syncTx.BroadcastStatus = SyncStatusComplete
	syncTx.Results.LastMessage = statusMsg
	syncTx.LastAttempt = customTypes.NullTime{
		NullTime: sql.NullTime{
			Time:  time.Now().UTC(),
			Valid: true,
		},
	}

	syncTx.Results.Results = append(syncTx.Results.Results, &SyncResult{
		Action:        syncActionBroadcast,
		ExecutedAt:    time.Now().UTC(),
		Provider:      provider,
		StatusMessage: statusMsg,
	})

	// Update sync status to be ready now
	if syncTx.SyncStatus == SyncStatusPending {
		syncTx.SyncStatus = SyncStatusReady
	}

	// Update the sync transaction record
	if err = syncTx.Save(ctx); err != nil {
		_bailAndSaveSyncTransaction(
			ctx, syncTx, SyncStatusError, syncActionBroadcast, "internal", err.Error(),
		)
		return err
	}

	// Fire a notification
	notify(notifications.EventTypeBroadcast, syncTx)

	return nil
}

/////////////////

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
			return ErrMissingTransaction
		}
	}

	// Find on-chain
	var txInfo *chainstate.TransactionInfo
	// only mAPI currently provides merkle proof, so QueryTransaction should be used here
	if txInfo, err = syncTx.Client().Chainstate().QueryTransaction(
		ctx, syncTx.ID, chainstate.RequiredOnChain, defaultQueryTxTimeout,
	); err != nil {
		if errors.Is(err, chainstate.ErrTransactionNotFound) {
			syncTx.Client().Logger().Info().
				Str("txID", syncTx.ID).
				Msgf("Transaction not found on-chain, will try again later")

			_bailAndSaveSyncTransaction(
				ctx, syncTx, SyncStatusReady, syncActionSync, "all", "transaction not found on-chain",
			)
			return nil
		}
		return err
	}
	return processSyncTxSave(ctx, txInfo, syncTx, transaction)
}

func _getTransaction(ctx context.Context, id string, opts []ModelOps) (*Transaction, error) {
	transaction, err := getTransactionByID(ctx, "", id, opts...)
	if err != nil {
		return nil, err
	}

	if transaction == nil {
		return nil, ErrMissingTransaction
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

	message := "transaction was found on-chain by " + chainstate.ProviderBroadcastClient

	if err := transaction.Save(ctx); err != nil {
		_bailAndSaveSyncTransaction(
			ctx, syncTx, SyncStatusError, syncActionSync, "internal", err.Error(),
		)
		return err
	}

	syncTx.SyncStatus = SyncStatusComplete
	syncTx.Results.LastMessage = message
	syncTx.Results.Results = append(syncTx.Results.Results, &SyncResult{
		Action:        syncActionSync,
		ExecutedAt:    time.Now().UTC(),
		Provider:      chainstate.ProviderBroadcastClient,
		StatusMessage: message,
	})

	if err := syncTx.Save(ctx); err != nil {
		_bailAndSaveSyncTransaction(ctx, syncTx, SyncStatusError, syncActionSync, "internal", err.Error())
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
		_bailAndSaveSyncTransaction(
			ctx, syncTx, SyncStatusComplete, syncActionP2P, "all", "no draft found, cannot complete p2p",
		)
		return nil
	}

	// Notify any P2P paymail providers associated to the transaction
	var results []*SyncResult
	if results, err = _notifyPaymailProviders(ctx, tx); err != nil {
		_bailAndSaveSyncTransaction(
			ctx, syncTx, SyncStatusReady, syncActionP2P, "", err.Error(),
		)
		return err
	}

	// Update if we have some results
	if len(results) > 0 {
		syncTx.Results.Results = append(syncTx.Results.Results, results...)
		syncTx.Results.LastMessage = fmt.Sprintf("notified %d paymail provider(s)", len(results))
	}

	// Save the record
	syncTx.P2PStatus = SyncStatusComplete

	// Update sync status to be ready now
	if syncTx.SyncStatus == SyncStatusPending {
		syncTx.SyncStatus = SyncStatusReady
	}

	if err = syncTx.Save(ctx); err != nil {
		_bailAndSaveSyncTransaction(
			ctx, syncTx, SyncStatusError, syncActionP2P, "internal", err.Error(),
		)
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

// _bailAndSaveSyncTransaction will save the error message for a sync tx
func _bailAndSaveSyncTransaction(ctx context.Context, syncTx *SyncTransaction, status SyncStatus,
	action, provider, message string,
) {
	if action == syncActionSync {
		syncTx.SyncStatus = status
	} else if action == syncActionP2P {
		syncTx.P2PStatus = status
	} else if action == syncActionBroadcast {
		syncTx.BroadcastStatus = status
	}
	syncTx.LastAttempt = customTypes.NullTime{
		NullTime: sql.NullTime{
			Time:  time.Now().UTC(),
			Valid: true,
		},
	}
	syncTx.Results.LastMessage = message
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

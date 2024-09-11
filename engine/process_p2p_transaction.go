package engine

import (
	"context"
	"fmt"
	"time"

	"github.com/bitcoin-sv/go-paymail"
)

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
		syncTx.addSyncResult(ctx, syncActionP2P, "all", "no draft found, cannot complete p2p")

		return nil // TODO: why nil here??
	}

	// Notify any P2P paymail providers associated to the transaction
	var results []*SyncResult
	if results, err = _notifyPaymailProviders(ctx, tx); err != nil {
		syncTx.addSyncResult(ctx, syncActionP2P, "", err.Error())
		return err
	}

	// Update if we have some results
	if len(results) > 0 {
		syncTx.Results.Results = append(syncTx.Results.Results, results...)
	}

	// Update sync status to be ready now
	if syncTx.SyncStatus == SyncStatusPending {
		syncTx.SyncStatus = SyncStatusReady
	}

	if err = syncTx.Save(ctx); err != nil {
		syncTx.addSyncResult(ctx, syncActionP2P, "internal", err.Error())
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

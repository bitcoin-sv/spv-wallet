package engine

import (
	"context"
	"fmt"
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
		return nil // TODO: why nil here??
	}

	// Notify any P2P paymail providers associated to the transaction
	if err = _notifyPaymailProviders(ctx, tx); err != nil {
		return err
	}

	// Update sync status to be ready now
	if syncTx.SyncStatus == SyncStatusPending {
		syncTx.SyncStatus = SyncStatusReady
	}

	if err = syncTx.Save(ctx); err != nil {
		return err
	}

	// Done!
	return nil
}

// _notifyPaymailProviders will notify any associated Paymail providers
func _notifyPaymailProviders(ctx context.Context, transaction *Transaction) error {
	pm := transaction.Client().PaymailClient()
	outputs := transaction.draftTransaction.Configuration.Outputs

	notifiedReceivers := make([]string, 0)
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

		if err = finalizeP2PTransaction(
			ctx,
			pm,
			p4,
			transaction,
		); err != nil {
			return err
		}

		notifiedReceivers = append(notifiedReceivers, receiver)
	}
	return nil
}

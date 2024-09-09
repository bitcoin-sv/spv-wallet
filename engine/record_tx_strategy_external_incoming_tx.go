package engine

import (
	"context"
	"fmt"

	"github.com/libsv/go-bt/v2"
)

type externalIncomingTx struct {
	BtTx         *bt.Tx
	broadcastNow bool // e.g. BEEF must be broadcasted now
	txID         string
}

func (strategy *externalIncomingTx) Name() string {
	return "external_incoming_tx"
}

func (strategy *externalIncomingTx) Execute(ctx context.Context, c ClientInterface, opts []ModelOps) (*Transaction, error) {
	logger := c.Logger()
	transaction, err := _createExternalTxToRecord(ctx, strategy, c, opts)
	if err != nil {
		logger.Error().Msgf("creation of external incoming tx failed. Reason: %v", err)
		return nil, err
	}

	logger.Info().
		Str("txID", transaction.ID).
		Msg("start without ITC")

	if strategy.broadcastNow || transaction.syncTransaction.BroadcastStatus == SyncStatusReady {
		if err := broadcastSyncTransaction(ctx, transaction.syncTransaction); err != nil {
			logger.Warn().Str("txID", transaction.ID).Msgf("broadcasting failed. Reason: %s", err)
			return nil, err
		}
	}

	// record
	if err = transaction.Save(ctx); err != nil {
		logger.Error().
			Str("txID", transaction.ID).
			Msgf("saving of Transaction failed. Reason: %v", err)
		return nil, err
	}

	logger.Info().
		Str("txID", transaction.ID).
		Msg("External incoming tx execute complete")
	return transaction, nil
}

func (strategy *externalIncomingTx) Validate() error {
	if strategy.BtTx == nil {
		return ErrMissingFieldHex
	}

	return nil // is valid
}

func (strategy *externalIncomingTx) TxID() string {
	if strategy.txID == "" {
		strategy.txID = strategy.BtTx.TxID()
	}
	return strategy.txID
}

func (strategy *externalIncomingTx) LockKey() string {
	return fmt.Sprintf("incoming-%s", strategy.TxID())
}

func (strategy *externalIncomingTx) ForceBroadcast(force bool) {
	strategy.broadcastNow = force
}

func _createExternalTxToRecord(ctx context.Context, eTx *externalIncomingTx, c ClientInterface, opts []ModelOps) (*Transaction, error) {
	// Create NEW tx model
	tx := txFromBtTx(eTx.BtTx, c.DefaultModelOptions(append(opts, New())...)...)
	_hydrateExternalWithSync(tx)

	if !tx.TransactionBase.hasOneKnownDestination(ctx, c) {
		return nil, ErrNoMatchingOutputs
	}

	if err := tx.processUtxos(ctx); err != nil {
		return nil, err
	}

	return tx, nil
}

func _hydrateExternalWithSync(tx *Transaction) {
	sync := newSyncTransaction(
		tx.ID,
		tx.Client().DefaultSyncConfig(),
		tx.GetOptions(true)...,
	)

	// to simplify: broadcast every external incoming txs
	sync.BroadcastStatus = SyncStatusReady

	sync.P2PStatus = SyncStatusSkipped  // the sender of the Tx should have already notified us
	sync.SyncStatus = SyncStatusPending // wait until transactions will be broadcasted

	// Use the same metadata
	sync.Metadata = tx.Metadata
	sync.transaction = tx
	tx.syncTransaction = sync
}

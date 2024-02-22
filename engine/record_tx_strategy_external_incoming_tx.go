package engine

import (
	"context"
	"fmt"

	"github.com/libsv/go-bt/v2"
	"github.com/rs/zerolog"
)

type externalIncomingTx struct {
	BtTx                 *bt.Tx
	broadcastNow         bool // e.g. BEEF must be broadcasted now
	allowBroadcastErrors bool // only BEEF cannot allow for broadcast errors
	txID                 string
}

func (strategy *externalIncomingTx) Name() string {
	return "external_incoming_tx"
}

func (strategy *externalIncomingTx) Execute(ctx context.Context, c ClientInterface, opts []ModelOps) (*Transaction, error) {
	logger := c.Logger()
	transaction, err := _createExternalTxToRecord(ctx, strategy, c, opts)
	if err != nil {
		return nil, fmt.Errorf("creation of external incoming tx failed. Reason: %w", err)
	}

	logger.Info().
		Str("txID", transaction.ID).
		Msg("start without ITC")

	if strategy.broadcastNow || transaction.syncTransaction.BroadcastStatus == SyncStatusReady {

		err = _externalIncomingBroadcast(ctx, logger, transaction, strategy.allowBroadcastErrors)
		if err != nil {
			logger.Error().
				Str("txID", transaction.ID).
				Msgf("broadcasting failed, transaction rejected! Reason: %s", err)

			return nil, fmt.Errorf("broadcasting failed, transaction rejected! Reason: %w", err)
		}
	}

	// record
	if err = transaction.Save(ctx); err != nil {
		return nil, fmt.Errorf("saving of Transaction failed. Reason: %w", err)
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

func (strategy *externalIncomingTx) FailOnBroadcastError(forceFail bool) {
	strategy.allowBroadcastErrors = !forceFail
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

func _externalIncomingBroadcast(ctx context.Context, logger *zerolog.Logger, tx *Transaction, allowErrors bool) error {
	logger.Info().Str("txID", tx.ID).Msg("start broadcast")

	err := broadcastSyncTransaction(ctx, tx.syncTransaction)

	if err == nil {
		logger.Info().
			Str("txID", tx.ID).
			Msg("broadcast complete")

		return nil
	}

	if allowErrors {
		logger.Warn().
			Str("txID", tx.ID).
			Msgf("broadcasting failed, next try will be handled by task manager. Reason: %s", err)

		// ignore error, transaction will be broadcasted in a cron task
		return nil
	}

	return err
}

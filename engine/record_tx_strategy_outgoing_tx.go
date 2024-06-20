package engine

import (
	"context"
	"errors"
	"fmt"
	"github.com/bitcoin-sv/spv-wallet/spverrors"

	"github.com/libsv/go-bt/v2"
	"github.com/rs/zerolog"
)

type outgoingTx struct {
	BtTx           *bt.Tx
	RelatedDraftID string
	XPubKey        string
	txID           string
}

func (strategy *outgoingTx) Name() string {
	return "outgoing_tx"
}

func (strategy *outgoingTx) Execute(ctx context.Context, c ClientInterface, opts []ModelOps) (*Transaction, error) {
	logger := c.Logger()
	logger.Info().
		Str("txID", strategy.TxID()).
		Msg("start recording outgoing transaction")

	// create
	transaction, err := _createOutgoingTxToRecord(ctx, strategy, c, opts)
	if err != nil {
		logger.Error().
			Str("txID", strategy.TxID()).
			Msgf("creation of outgoing tx failed. Reason: %v", err)
		return nil, spverrors.ErrCreateOutgoingTxFailed
	}

	if err = transaction.Save(ctx); err != nil {
		logger.Error().
			Str("txID", strategy.TxID()).
			Msgf("saving of Transaction failed. Reason: %v", err)
		return nil, spverrors.ErrDuringSaveTx
	}

	// process
	if transaction.syncTransaction.P2PStatus == SyncStatusReady {
		if err = _outgoingNotifyP2p(ctx, logger, transaction); err != nil {
			// reject transaction if P2P notification failed
			logger.Error().
				Str("txID", transaction.ID).
				Msgf("transaction rejected by P2P provider, try to revert transaction. Reason: %s", err)

			if revertErr := c.RevertTransaction(ctx, transaction.ID); revertErr != nil {
				logger.Error().
					Str("txID", transaction.ID).
					Msgf("FATAL! Reverting transaction after failed P2P notification failed. Reason: %s", revertErr)
			}

			return nil, err
		}
	}

	// get newest syncTx from DB - if it's an internal tx it could be broadcasted by us already
	syncTx, err := GetSyncTransactionByID(ctx, transaction.ID, transaction.GetOptions(false)...)
	if err != nil || syncTx == nil {
		logger.Error().
			Str("txID", transaction.ID).
			Msgf("getting syncTx failed. Reason: %v", err)
		return nil, spverrors.ErrCouldNotFindSyncTx
	}

	if syncTx.BroadcastStatus == SyncStatusReady {
		transaction.syncTransaction = syncTx
		_outgoingBroadcast(ctx, logger, transaction) // ignore error
	}

	logger.Info().
		Str("txID", transaction.ID).
		Msgf("complete recording outgoing transaction")
	return transaction, nil
}

func (strategy *outgoingTx) Validate() error {
	if strategy.BtTx == nil {
		return ErrMissingFieldHex
	}

	if strategy.RelatedDraftID == "" {
		return errors.New("empty RelatedDraftID")
	}

	if strategy.XPubKey == "" {
		return errors.New("empty xPubKey")
	}

	return nil // is valid
}

func (strategy *outgoingTx) TxID() string {
	if strategy.txID == "" {
		strategy.txID = strategy.BtTx.TxID()
	}
	return strategy.txID
}

func (strategy *outgoingTx) LockKey() string {
	return fmt.Sprintf("outgoing-%s", strategy.TxID())
}

func _createOutgoingTxToRecord(ctx context.Context, oTx *outgoingTx, c ClientInterface, opts []ModelOps) (*Transaction, error) {
	// Create NEW transaction model
	newOpts := c.DefaultModelOptions(append(opts, WithXPub(oTx.XPubKey), New())...)
	tx := txFromBtTx(oTx.BtTx, newOpts...)
	tx.DraftID = oTx.RelatedDraftID

	// hydrate
	if err := _hydrateOutgoingWithDraft(ctx, tx); err != nil {
		return nil, err
	}

	_hydrateOutgoingWithSync(tx)

	if err := tx.processUtxos(ctx); err != nil {
		return nil, err
	}

	return tx, nil
}

func _hydrateOutgoingWithDraft(ctx context.Context, tx *Transaction) error {
	draft, err := getDraftTransactionID(ctx, tx.XPubID, tx.DraftID, tx.GetOptions(false)...)
	if err != nil {
		return err
	}

	if draft == nil {
		return spverrors.ErrCouldNotFindDraftTx
	}

	if len(draft.Configuration.Outputs) == 0 {
		return spverrors.ErrDraftTxHasNoOutputs
	}

	if draft.Configuration.Sync == nil {
		draft.Configuration.Sync = tx.Client().DefaultSyncConfig()
	}

	tx.draftTransaction = draft

	return nil // success
}

func _hydrateOutgoingWithSync(tx *Transaction) {
	sync := newSyncTransaction(tx.ID, tx.draftTransaction.Configuration.Sync, tx.GetOptions(true)...)

	// setup synchronization
	sync.BroadcastStatus = _getBroadcastSyncStatus(tx)
	sync.P2PStatus = _getP2pSyncStatus(tx)
	sync.SyncStatus = SyncStatusPending // wait until transaction is broadcasted or P2P provider is notified

	sync.Metadata = tx.Metadata

	sync.transaction = tx
	tx.syncTransaction = sync
}

func _getBroadcastSyncStatus(tx *Transaction) SyncStatus {
	// immediately broadcast if is not BEEF
	broadcast := SyncStatusReady // broadcast immediately

	outputs := tx.draftTransaction.Configuration.Outputs

	for _, o := range outputs {
		if o.PaymailP4 != nil {
			if o.PaymailP4.Format == BeefPaymailPayloadFormat {
				broadcast = SyncStatusSkipped // postpone broadcasting if tx contains outputs in BEEF

				break
			}
		}
	}

	return broadcast
}

func _getP2pSyncStatus(tx *Transaction) SyncStatus {
	p2pStatus := SyncStatusSkipped

	outputs := tx.draftTransaction.Configuration.Outputs
	for _, o := range outputs {
		if o.PaymailP4 != nil && o.PaymailP4.ResolutionType == ResolutionTypeP2P {
			p2pStatus = SyncStatusReady // notify p2p immediately

			break
		}
	}

	return p2pStatus
}

func _outgoingNotifyP2p(ctx context.Context, logger *zerolog.Logger, tx *Transaction) error {
	logger.Info().
		Str("txID", tx.ID).
		Msg("start p2p")

	if err := processP2PTransaction(ctx, tx); err != nil {
		logger.Error().
			Str("txID", tx.ID).
			Msgf("processP2PTransaction failed. Reason: %s", err)

		return spverrors.ErrProcessP2PTx
	}

	logger.Info().
		Str("txID", tx.ID).
		Msg("p2p complete")
	return nil
}

func _outgoingBroadcast(ctx context.Context, logger *zerolog.Logger, tx *Transaction) {
	logger.Info().
		Str("txID", tx.ID).
		Msg("start broadcast")

	if err := broadcastSyncTransaction(ctx, tx.syncTransaction); err != nil {
		// ignore error, transaction will be broadcasted by cron task
		logger.Warn().
			Str("txID", tx.ID).
			Msgf("broadcasting failed, next try will be handled by task manager. Reason: %s", err)
	} else {
		logger.Info().
			Str("txID", tx.ID).
			Msg("broadcast complete")
	}
}

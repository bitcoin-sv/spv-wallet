package engine

import (
	"context"
	"fmt"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/libsv/go-bt/v2"
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

	transaction, err := _createOutgoingTxToRecord(ctx, strategy, c, opts)
	transaction.TxStatus = TxStatusCreated
	if err != nil {
		return nil, spverrors.ErrCreateOutgoingTxFailed
	}

	if err = transaction.Save(ctx); err != nil {
		return nil, spverrors.ErrDuringSaveTx
	}

	notifyP2P := _shouldNotifyP2P(transaction)

	if notifyP2P {
		if err := processP2PTransaction(ctx, transaction); err != nil {
			if revertErr := c.RevertTransaction(ctx, transaction.ID); revertErr != nil {
				return nil, fmt.Errorf("reverting transaction failed %w; after P2P notification failed: %w", revertErr, err)
			}
			logger.Warn().Str("txID", transaction.ID).Msgf("processP2PTransaction failed. Reason: %v", err)
			return nil, spverrors.ErrProcessP2PTx
		}
	}

	// transaction can be updated by internal_incoming_tx
	transaction, err = _getTransaction(ctx, transaction.ID, WithClient(c))
	if err != nil {
		return nil, err
	}

	if transaction.TxStatus == TxStatusBroadcasted {
		// no need to broadcast twice
		return transaction, nil
	}

	if notifyP2P {
		transaction.TxStatus = TxStatusSent
	}

	if err := broadcastTransaction(ctx, transaction); err != nil {
		logger.Warn().Str("txID", transaction.ID).Msgf("broadcasting failed in outgoingTx strategy")
		// ignore error, transaction most likely is successfully broadcasted by payment receiver
		// TODO: return a Warning to a client
	} else {
		transaction.TxStatus = TxStatusBroadcasted
	}

	if err := transaction.Save(ctx); err != nil {
		c.Logger().Error().Str("txID", transaction.ID).Err(err).Msgf("Outgoing transaction has been %v but failed save to db", transaction.TxStatus)
	}

	return transaction, nil
}

func (strategy *outgoingTx) Validate() error {
	if strategy.BtTx == nil {
		return ErrMissingFieldHex
	}

	if strategy.RelatedDraftID == "" {
		return spverrors.ErrEmptyRelatedDraftID
	}

	if strategy.XPubKey == "" {
		return spverrors.ErrEmptyXpubKey
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
	sync.SyncStatus = SyncStatusPending // wait until transaction is broadcasted or P2P provider is notified

	sync.Metadata = tx.Metadata

	sync.transaction = tx
	tx.syncTransaction = sync
}

func _shouldNotifyP2P(tx *Transaction) bool {
	for _, o := range tx.draftTransaction.Configuration.Outputs {
		if o.PaymailP4 != nil && o.PaymailP4.ResolutionType == ResolutionTypeP2P {
			return true
		}
	}
	return false
}

func _getTransaction(ctx context.Context, id string, opts ...ModelOps) (*Transaction, error) {
	transaction, err := getTransactionByID(ctx, "", id, opts...)
	if err != nil {
		return nil, err
	}

	if transaction == nil {
		return nil, spverrors.ErrCouldNotFindTransaction
	}

	return transaction, nil
}

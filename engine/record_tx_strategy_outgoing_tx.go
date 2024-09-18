package engine

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/bitcoin-sv/spv-wallet/engine/chainstate"
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

	var transaction *Transaction
	var err error

	if transaction, err = strategy.createOutgoingTxToRecord(ctx, c, opts); err != nil {
		return nil, spverrors.ErrCreateOutgoingTxFailed
	}

	if err = transaction.processUtxos(ctx); err != nil {
		return nil, err
	}

	if err = transaction.Save(ctx); err != nil {
		return nil, spverrors.ErrDuringSaveTx
	}

	if _shouldNotifyP2P(transaction) {
		if err = processP2PTransaction(ctx, transaction); err != nil {
			return nil, _handleNotifyP2PError(ctx, c, transaction)
		}
	}

	// transaction can be updated by internal_incoming_tx
	transaction, err = getTransactionByID(ctx, "", transaction.ID, WithClient(c))
	if transaction == nil || err != nil {
		logger.Error().Msg("Cannot find transaction even though it was saved a moment ago")
		return nil, spverrors.ErrInternal
	}

	if transaction.TxStatus == TxStatusBroadcasted {
		// no need to broadcast twice
		return transaction, nil
	}

	if err = broadcastTransaction(ctx, transaction); err != nil {
		logger.Warn().Str("txID", transaction.ID).Msgf("broadcasting failed in outgoingTx strategy")
		// ignore error, transaction most likely is successfully broadcasted by payment receiver
		// TODO: return a Warning to a client
	} else {
		transaction.TxStatus = TxStatusBroadcasted
	}

	if err = transaction.Save(ctx); err != nil {
		logger.Error().Str("txID", transaction.ID).Err(err).Msg("Outgoing transaction has been processed but failed save to db")
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

func (strategy *outgoingTx) createOutgoingTxToRecord(ctx context.Context, c ClientInterface, opts []ModelOps) (*Transaction, error) {
	// Create NEW transaction model
	newOpts := c.DefaultModelOptions(append(opts, WithXPub(strategy.XPubKey), New())...)
	tx := txFromBtTx(strategy.BtTx, newOpts...)
	tx.DraftID = strategy.RelatedDraftID

	if err := _hydrateOutgoingWithDraft(ctx, tx); err != nil {
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

func _shouldNotifyP2P(tx *Transaction) bool {
	return slices.ContainsFunc(tx.draftTransaction.Configuration.Outputs, func(o *TransactionOutput) bool {
		return o.PaymailP4 != nil && o.PaymailP4.ResolutionType == ResolutionTypeP2P
	})
}

func _handleNotifyP2PError(ctx context.Context, c ClientInterface, transaction *Transaction) error {
	logger := c.Logger()
	chainstateService := c.Chainstate()

	saveAsProblematic := func() {
		transaction.TxStatus = TxStatusProblematic
		if err := transaction.Save(ctx); err != nil {
			logger.Error().Str("txID", transaction.ID).Err(err).Msg("Error saving transaction after notifyP2P failed")
		}
	}

	_, err := chainstateService.QueryTransaction(ctx, transaction.ID, chainstate.RequiredInMempool, defaultQueryTxTimeout)
	if errors.Is(err, spverrors.ErrCouldNotFindTransaction) {
		if err := c.RevertTransaction(ctx, transaction.ID); err != nil {
			saveAsProblematic()
			return spverrors.ErrProcessP2PTx.Wrap(err)
		}

		// RevertTransaction saves the transaction itself as REVERTED
		return spverrors.ErrProcessP2PTx
	}
	saveAsProblematic()
	return spverrors.ErrProcessP2PTx
}

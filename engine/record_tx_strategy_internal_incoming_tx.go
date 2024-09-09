package engine

import (
	"context"
	"fmt"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/libsv/go-bt/v2"
)

type internalIncomingTx struct {
	Tx                   *Transaction
	broadcastNow         bool // e.g. BEEF must be broadcasted now
	allowBroadcastErrors bool // only BEEF cannot allow for broadcast errors
}

func (strategy *internalIncomingTx) Name() string {
	return "internal_incoming_tx"
}

func (strategy *internalIncomingTx) Execute(ctx context.Context, c ClientInterface, _ []ModelOps) (*Transaction, error) {
	logger := c.Logger()
	logger.Info().
		Str("txID", strategy.Tx.ID).
		Msg("start recording internal incoming transaction")

	// process
	transaction := strategy.Tx
	syncTx, err := GetSyncTransactionByID(ctx, transaction.ID, transaction.GetOptions(false)...)
	if err != nil {
		return nil, err
	} else if syncTx == nil {
		return nil, spverrors.ErrCouldNotFindSyncTx
	}

	if strategy.broadcastNow || syncTx.BroadcastStatus == SyncStatusReady {
		syncTx.transaction = transaction
		transaction.syncTransaction = syncTx

		if err := broadcastSyncTransaction(ctx, syncTx); err != nil {
			logger.Warn().Str("txID", transaction.ID).Bool("ignored", !strategy.allowBroadcastErrors).Msgf("broadcasting failed. Reason: %s", err)
			if !strategy.allowBroadcastErrors {
				return nil, err
			}
		}
	}

	logger.Info().
		Str("txID", transaction.ID).
		Msg("complete recording internal incoming transaction")
	return transaction, nil
}

func (strategy *internalIncomingTx) Validate() error {
	if strategy.Tx == nil {
		return spverrors.ErrEmptyTx
	}

	if _, err := bt.NewTxFromString(strategy.Tx.Hex); err != nil {
		return spverrors.ErrInvalidHex
	}

	return nil // is valid
}

func (strategy *internalIncomingTx) TxID() string {
	return strategy.Tx.ID
}

func (strategy *internalIncomingTx) LockKey() string {
	return fmt.Sprintf("incoming-%s", strategy.Tx.ID)
}

func (strategy *internalIncomingTx) ForceBroadcast(force bool) {
	strategy.broadcastNow = force
}

func (strategy *internalIncomingTx) FailOnBroadcastError(forceFail bool) {
	strategy.allowBroadcastErrors = !forceFail
}

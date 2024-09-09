package engine

import (
	"context"
	"fmt"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/libsv/go-bt/v2"
)

type internalIncomingTx struct {
	Tx *Transaction
}

func (strategy *internalIncomingTx) Name() string {
	return "internal_incoming_tx"
}

func (strategy *internalIncomingTx) Execute(ctx context.Context, _ ClientInterface, _ []ModelOps) (*Transaction, error) {
	transaction := strategy.Tx
	syncTx, err := GetSyncTransactionByID(ctx, transaction.ID, transaction.GetOptions(false)...)
	if err != nil {
		return nil, err
	} else if syncTx == nil {
		return nil, spverrors.ErrCouldNotFindSyncTx
	}

	syncTx.transaction = transaction
	transaction.syncTransaction = syncTx

	if err := broadcastSyncTransaction(ctx, syncTx); err != nil {
		return nil, err
	}

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

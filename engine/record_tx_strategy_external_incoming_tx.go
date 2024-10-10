package engine

import (
	"context"
	"fmt"

	trx "github.com/bitcoin-sv/go-sdk/transaction"
)

type externalIncomingTx struct {
	SDKTx *trx.Transaction
	txID  string
}

func (strategy *externalIncomingTx) Name() string {
	return "external_incoming_tx"
}

func (strategy *externalIncomingTx) Execute(ctx context.Context, c ClientInterface, opts []ModelOps) (*Transaction, error) {
	transaction, err := _createExternalTxToRecord(ctx, strategy, c, opts)
	if err != nil {
		return nil, err
	}
	if err = broadcastTransaction(ctx, transaction); err != nil {
		return nil, err
	}
	transaction.TxStatus = TxStatusBroadcasted

	if err := transaction.Save(ctx); err != nil {
		c.Logger().Error().Str("txID", transaction.ID).Err(err).Msg("Incoming external transaction has been broadcasted but failed save to db")
	}

	return transaction, nil
}

func (strategy *externalIncomingTx) Validate() error {
	if strategy.SDKTx == nil {
		return ErrMissingFieldHex
	}

	return nil // is valid
}

func (strategy *externalIncomingTx) TxID() string {
	if strategy.txID == "" {
		strategy.txID = strategy.SDKTx.TxID().String()
	}
	return strategy.txID
}

func (strategy *externalIncomingTx) LockKey() string {
	return fmt.Sprintf("incoming-%s", strategy.TxID())
}

func _createExternalTxToRecord(ctx context.Context, eTx *externalIncomingTx, c ClientInterface, opts []ModelOps) (*Transaction, error) {
	// Create NEW tx model
	tx := txFromSDKTx(eTx.SDKTx, c.DefaultModelOptions(append(opts, New())...)...)

	if !tx.TransactionBase.hasOneKnownDestination(ctx, c) {
		return nil, ErrNoMatchingOutputs
	}

	if err := tx.processUtxos(ctx); err != nil {
		return nil, err
	}

	return tx, nil
}

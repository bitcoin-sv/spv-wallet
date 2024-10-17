package engine

import (
	"context"
	"fmt"

	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
)

func broadcastTransaction(ctx context.Context, txModel *Transaction) error {
	client := txModel.Client()

	defer recoverAndLog(txModel.Client().Logger())

	unlock, err := newWriteLock(
		ctx, fmt.Sprintf(lockKeyProcessBroadcastTx, txModel.GetID()), client.Cachestore(),
	)
	defer unlock()
	if err != nil {
		return err
	}

	tx, err := sdk.NewTransactionFromHex(txModel.Hex)
	if err != nil {
		return spverrors.ErrParseTransactionFromHex.Wrap(err)
	}

	_, err = client.Chain().Broadcast(ctx, tx)
	if err != nil {
		return spverrors.ErrBroadcast.Wrap(err)
	}
	return nil
}

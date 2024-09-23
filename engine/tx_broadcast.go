package engine

import (
	"context"
	"fmt"

	"github.com/bitcoin-sv/spv-wallet/engine/chainstate"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
)

func broadcastTransaction(ctx context.Context, tx *Transaction) error {
	client := tx.Client()
	chainstateSrv := client.Chainstate()

	defer recoverAndLog(tx.Client().Logger())

	unlock, err := newWriteLock(
		ctx, fmt.Sprintf(lockKeyProcessBroadcastTx, tx.GetID()), client.Cachestore(),
	)
	defer unlock()
	if err != nil {
		return err
	}

	txHex := tx.Hex
	hexFormat := chainstate.RawTx
	if chainstateSrv.SupportedBroadcastFormats().Contains(chainstate.Ef) {
		// try to convert to EF, with rawRx as a fallback
		efHex, ok := ToEfHex(ctx, tx, client)
		if ok {
			txHex = efHex
			hexFormat = chainstate.Ef
		}
	}

	br := chainstateSrv.Broadcast(ctx, tx.ID, txHex, hexFormat, defaultBroadcastTimeout)

	if br.Failure != nil {
		if br.Failure.InvalidTx {
			return spverrors.ErrBroadcastRejectedTransaction.Wrap((br.Failure.Error))
		}
		return br.Failure.Error
	}

	return nil
}

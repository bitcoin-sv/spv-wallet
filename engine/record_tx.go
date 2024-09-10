package engine

import (
	"context"
	"fmt"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/libsv/go-bt/v2"
)

type recordTxStrategy interface {
	Name() string
	TxID() string
	LockKey() string
	Validate() error
	Execute(ctx context.Context, c ClientInterface, opts []ModelOps) (*Transaction, error)
}

type recordIncomingTxStrategy interface {
	recordTxStrategy
}

func recordTransaction(ctx context.Context, c ClientInterface, strategy recordTxStrategy, opts ...ModelOps) (transaction *Transaction, err error) {
	if metrics, enabled := c.Metrics(); enabled {
		end := metrics.TrackRecordTransaction(strategy.Name())
		defer func() {
			success := err == nil
			end(success)
		}()
	}

	unlock, err := newWriteLock(ctx, fmt.Sprintf(lockKeyRecordTx, strategy.LockKey()), c.Cachestore())
	defer unlock()
	if err != nil {
		return nil, spverrors.ErrInternal.Wrap(err)
	}

	logger := c.Logger()
	logger.Debug().Str("strategy", strategy.Name()).Str("txID", strategy.TxID()).Msg("Start executing recordTx strategy.")

	transaction, err = strategy.Execute(ctx, c, opts)
	if err != nil {
		logger.Warn().Str("strategy", strategy.Name()).Str("txID", strategy.TxID()).Err(err).Msg("Failed to execure recordTx strategy.")
	}
	return
}

func getOutgoingTxRecordStrategy(xPubKey string, btTx *bt.Tx, draftID string) (recordTxStrategy, error) {
	rts := &outgoingTx{
		BtTx:           btTx,
		RelatedDraftID: draftID,
		XPubKey:        xPubKey,
	}

	if err := rts.Validate(); err != nil {
		return nil, err
	}

	return rts, nil
}

func getIncomingTxRecordStrategy(ctx context.Context, c ClientInterface, btTx *bt.Tx) (recordIncomingTxStrategy, error) {
	tx, err := getTransactionByHex(ctx, btTx.String(), c.DefaultModelOptions()...)
	if err != nil {
		return nil, err
	}

	var rts recordIncomingTxStrategy

	if tx != nil {
		tx.parsedTx = btTx
		rts = &internalIncomingTx{
			Tx: tx,
		}
	} else {
		rts = &externalIncomingTx{
			BtTx: btTx,
		}
	}

	if err := rts.Validate(); err != nil {
		return nil, err //nolint:wrapcheck // wrapped by our code below
	}

	return rts, nil
}

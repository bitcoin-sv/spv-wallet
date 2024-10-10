package engine

import (
	"context"
	"fmt"

	trx "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
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
		logger.Warn().Str("strategy", strategy.Name()).Str("txID", strategy.TxID()).Err(err).Msg("Failed to execute recordTx strategy.")
	}
	return
}

func getOutgoingTxRecordStrategy(xPubKey string, sdkTx *trx.Transaction, draftID string) (recordTxStrategy, error) {
	rts := &outgoingTx{
		SDKTx:          sdkTx,
		RelatedDraftID: draftID,
		XPubKey:        xPubKey,
	}

	if err := rts.Validate(); err != nil {
		return nil, err
	}

	return rts, nil
}

func getIncomingTxRecordStrategy(ctx context.Context, c ClientInterface, sdkTx *trx.Transaction) (recordIncomingTxStrategy, error) {
	tx, err := getTransactionByHex(ctx, sdkTx.String(), c.DefaultModelOptions()...)
	if err != nil {
		return nil, err
	}

	var rts recordIncomingTxStrategy

	if tx != nil {
		tx.parsedTx = sdkTx
		rts = &internalIncomingTx{
			Tx: tx,
		}
	} else {
		rts = &externalIncomingTx{
			SDKTx: sdkTx,
		}
	}

	if err := rts.Validate(); err != nil {
		return nil, err //nolint:wrapcheck // wrapped by our code below
	}

	return rts, nil
}

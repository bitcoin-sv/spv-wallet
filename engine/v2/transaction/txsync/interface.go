package txsync

import (
	"context"

	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/txmodels"
)

// TransactionsRepo is an interface for transactions repository.
type TransactionsRepo interface {
	UpdateTransaction(ctx context.Context, trackedTx *txmodels.TrackedTransaction) error
	GetTransaction(ctx context.Context, txID string) (transaction *txmodels.TrackedTransaction, err error)
}

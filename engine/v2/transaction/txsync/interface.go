package txsync

import (
	"context"

	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/txmodels"
)

// TransactionsRepo is an interface for transactions repository.
type TransactionsRepo interface {
	SetStatus(ctx context.Context, txID string, status txmodels.TxStatus) error
}

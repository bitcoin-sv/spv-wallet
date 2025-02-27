package txsync

import (
	"context"

	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/txmodels"
)

// TransactionsRepo is an interface for transactions repository.
type TransactionsRepo interface {
	SetStatus(ctx context.Context, txID string, status txmodels.TxStatus) error
	SetAsMined(ctx context.Context, txID string, blockHash string, blockHeight int64, beefHex string) error
	GetTransactionHex(ctx context.Context, txID string) (hex string, isBEEF bool, err error)
}

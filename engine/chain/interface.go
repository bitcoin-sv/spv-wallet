package chain

import (
	"context"

	"github.com/bitcoin-sv/spv-wallet/engine/chain/models"
)

// QueryService for querying transactions.
type QueryService interface {
	Query(ctx context.Context, txID string) (*chainmodels.TXInfo, error)
}

// Service related to the chain.
type Service interface {
	QueryService
}

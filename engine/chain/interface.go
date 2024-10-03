package chain

import (
	"context"

	"github.com/bitcoin-sv/spv-wallet/engine/chain/models"
)

// QueryService for querying transactions.
type QueryService interface {
	QueryTransaction(ctx context.Context, txID string) (*chainmodels.TXInfo, error)
}

// PolicyService for querying policy.
type PolicyService interface {
	GetPolicy(ctx context.Context) (*chainmodels.Policy, error)
}

// Service related to the chain.
type Service interface {
	QueryService
	PolicyService
}

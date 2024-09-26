package chain

import (
	"context"
	"github.com/bitcoin-sv/spv-wallet/engine/chain/models"
)

type QueryService interface {
	Query(ctx context.Context, txID string) (*chainmodels.TXInfo, chainmodels.QueryTXOutcome, error)
}

type Service interface {
	QueryService
}

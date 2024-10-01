package chain

import (
	"github.com/bitcoin-sv/spv-wallet/engine/chain/internal/query"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog"
)

type chainService struct {
	QueryService
}

// NewChainService creates a new chain service.
func NewChainService(logger zerolog.Logger, httpClient *resty.Client, arcURL, arcToken, deploymentID string) Service {
	return &chainService{
		query.NewQueryService(logger.With().Str("chain", "query").Logger(), httpClient, arcURL, arcToken, deploymentID),
	}
}

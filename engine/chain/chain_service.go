package chain

import (
	"github.com/bitcoin-sv/spv-wallet/engine/chain/internal/query"
	"github.com/rs/zerolog"
)

type chainService struct {
	QueryService
}

// NewChainService creates a new chain service.
func NewChainService(logger zerolog.Logger) Service {
	return &chainService{
		query.NewQueryService(logger.With().Str("chain", "query").Logger()),
	}
}

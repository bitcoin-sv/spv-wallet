package chain

import (
	"github.com/bitcoin-sv/spv-wallet/engine/chain/internal/policy"
	"github.com/bitcoin-sv/spv-wallet/engine/chain/internal/query"
	"github.com/bitcoin-sv/spv-wallet/engine/chain/models"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog"
)

type chainService struct {
	QueryService
	PolicyService
}

// NewChainService creates a new chain service.
func NewChainService(logger zerolog.Logger, httpClient *resty.Client, arcCfg chainmodels.ARCConfig) Service {
	return &chainService{
		query.NewQueryService(logger.With().Str("chain", "query").Logger(), httpClient, arcCfg),
		policy.NewPolicyService(logger.With().Str("chain", "policy").Logger(), httpClient, arcCfg),
	}
}

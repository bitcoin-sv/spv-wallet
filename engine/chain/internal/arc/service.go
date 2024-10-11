package arc

import (
	"github.com/bitcoin-sv/spv-wallet/engine/chain/models"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog"
)

// Service for arc requests.
type Service struct {
	logger     zerolog.Logger
	httpClient *resty.Client
	arcCfg     chainmodels.ARCConfig
}

// NewARCService creates a new arc service.
func NewARCService(logger zerolog.Logger, httpClient *resty.Client, arcCfg chainmodels.ARCConfig) *Service {
	return &Service{
		logger:     logger,
		httpClient: httpClient,
		arcCfg:     arcCfg,
	}
}

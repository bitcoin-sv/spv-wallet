package bhs

import (
	"github.com/bitcoin-sv/spv-wallet/engine/chain/models"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog"
)

// Service for BHS requests.
type Service struct {
	logger     zerolog.Logger
	httpClient *resty.Client
	bhsCfg     chainmodels.BHSConfig
}

// NewBHSService creates a new instance of BHS service.
func NewBHSService(logger zerolog.Logger, httpClient *resty.Client, bhsCfg chainmodels.BHSConfig) *Service {
	return &Service{
		logger:     logger,
		httpClient: httpClient,
		bhsCfg:     bhsCfg,
	}
}

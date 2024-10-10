package arc

import (
	"github.com/bitcoin-sv/spv-wallet/engine/chain/internal/junglebus"
	chainmodels "github.com/bitcoin-sv/spv-wallet/engine/chain/models"
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
	if arcCfg.UseJunglebus && arcCfg.TxsGetter != nil {
		arcCfg.TxsGetter = newCombinedTxsGetter(
			arcCfg.TxsGetter,
			junglebus.NewJunglebusService(logger.With().Str("service", "junglebus").Logger(), httpClient),
		)
	}
	return &Service{
		logger:     logger,
		httpClient: httpClient,
		arcCfg:     arcCfg,
	}
}

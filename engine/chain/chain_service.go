package chain

import (
	"github.com/bitcoin-sv/spv-wallet/engine/chain/internal"
	"github.com/bitcoin-sv/spv-wallet/engine/chain/internal/arc"
	"github.com/bitcoin-sv/spv-wallet/engine/chain/internal/bhs"
	"github.com/bitcoin-sv/spv-wallet/engine/chain/internal/junglebus"
	"github.com/bitcoin-sv/spv-wallet/engine/chain/models"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog"
)

type arcService = arc.Service

type chainService struct {
	*arcService
	BHSService
}

// NewChainService creates a new chain service.
func NewChainService(logger zerolog.Logger, httpClient *resty.Client, arcCfg chainmodels.ARCConfig, bhsConf chainmodels.BHSConfig) Service {
	if httpClient == nil {
		panic("httpClient is required")
	}

	if arcCfg.UseJunglebus {
		addJunglebusTxsGetter(logger, httpClient, &arcCfg)
	}

	return &chainService{
		arc.NewARCService(logger.With().Str("chain", "arc").Logger(), httpClient, arcCfg),
		bhs.NewBHSService(logger.With().Str("chain", "bhs").Logger(), httpClient, bhsConf),
	}
}

func addJunglebusTxsGetter(logger zerolog.Logger, httpClient *resty.Client, arcCfg *chainmodels.ARCConfig) {
	junglebusTxsGetter := junglebus.NewJunglebusService(logger.With().Str("service", "junglebus").Logger(), httpClient)

	arcCfg.TxsGetter = internal.CombineTxsGetters(
		arcCfg.TxsGetter,
		junglebusTxsGetter,
	)
}

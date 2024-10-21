package arc

import (
	"context"

	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/chain/errors"
	"github.com/bitcoin-sv/spv-wallet/engine/chain/internal/ef"
	"github.com/bitcoin-sv/spv-wallet/engine/chain/models"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog"
)

// Service for arc requests.
type Service struct {
	logger      zerolog.Logger
	httpClient  *resty.Client
	arcCfg      chainmodels.ARCConfig
	efConverter interface {
		Convert(ctx context.Context, tx *sdk.Transaction) (string, error)
	}
}

// NewARCService creates a new arc service.
func NewARCService(logger zerolog.Logger, httpClient *resty.Client, arcCfg chainmodels.ARCConfig) *Service {
	service := &Service{
		logger:     logger,
		httpClient: httpClient,
		arcCfg:     arcCfg,
	}

	if arcCfg.TxsGetter == nil {
		logger.Warn().Msg("No transactions getter provided. Unsourced transactions will be broadcasted as raw hex.")
		service.efConverter = &noTxsGettersEFConverter{}
	} else {
		service.efConverter = ef.NewConverter(arcCfg.TxsGetter)
	}

	return service
}

type noTxsGettersEFConverter struct{}

func (n *noTxsGettersEFConverter) Convert(_ context.Context, tx *sdk.Transaction) (string, error) {
	efHex, err := tx.EFHex()
	if err != nil {
		return "", chainerrors.ErrEFConversion.Wrap(err)
	}
	return efHex, nil
}

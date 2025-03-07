package fee

import (
	"context"

	"github.com/bitcoin-sv/spv-wallet/config"
	"github.com/bitcoin-sv/spv-wallet/conv"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/utils/must"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"github.com/bitcoin-sv/spv-wallet/models/optional"
	"github.com/rs/zerolog"
	"github.com/samber/lo"
)

// Provider is an interface that provides fee units from miners.
type Provider interface {
	GetFeeUnit(ctx context.Context) (*bsv.FeeUnit, error)
}

// Service is a fee service that provides fee units for transactions.
type Service struct {
	logger      zerolog.Logger
	feeUnit     optional.Param[bsv.FeeUnit]
	feeProvider Provider
}

// NewService creates a new fee service.
func NewService(cfg *config.AppConfig, feeProvider Provider, logger zerolog.Logger) *Service {
	must.BeTrue(cfg != nil, "config is required")
	must.BeTrue(feeProvider != nil, "feeProvider is required")
	logger = logger.With().Str("service", "fee").Logger()

	return &Service{
		feeUnit: lo.IfF(cfg.CustomFeeUnit != nil,
			func() optional.Param[bsv.FeeUnit] {
				satoshis, err := conv.IntToUint64(cfg.CustomFeeUnit.Satoshis)
				must.HaveNoErrorf(err, "error converting custom fee unit %d satoshis", cfg.CustomFeeUnit.Satoshis)
				logger.Log().
					Int("satoshis", cfg.CustomFeeUnit.Satoshis).
					Int("bytes", cfg.CustomFeeUnit.Bytes).
					Msg("Fee unit found in configuration, using custom fee unit")

				return &bsv.FeeUnit{
					Satoshis: bsv.Satoshis(satoshis),
					Bytes:    cfg.CustomFeeUnit.Bytes,
				}
			}).
			Else(nil),
		feeProvider: feeProvider,
	}
}

// GetFeeUnit returns the fee unit that should be used for transactions.
func (s *Service) GetFeeUnit(ctx context.Context) (bsv.FeeUnit, error) {
	if s.feeUnit != nil {
		return *s.feeUnit, nil
	}

	s.logger.Debug().Msg("Fee unit not found in config, will try to receive it from miners")

	feeUnit, err := s.feeProvider.GetFeeUnit(ctx)
	if err != nil {
		return bsv.FeeUnit{}, spverrors.Wrapf(err, "failed to get fee unit from miners")
	}
	if feeUnit == nil {
		return bsv.FeeUnit{}, spverrors.Newf("received empty fee unit from miners")
	}
	s.feeUnit = feeUnit

	s.logger.Log().Any("satoshis", feeUnit.Satoshis).Int("bytes", feeUnit.Bytes).
		Msg("Received fee unit from miners, will use it from now on")

	return *feeUnit, nil
}

package query

import (
	"context"
	"github.com/bitcoin-sv/spv-wallet/engine/chain/models"
	"github.com/rs/zerolog"
)

type Service struct {
	logger zerolog.Logger
}

func (s *Service) Query(ctx context.Context, txID string) (*chainmodels.TXInfo, chainmodels.QueryTXOutcome, error) {
	// TODO implement
	return nil, chainmodels.QueryTxOutcomeFailed, nil
}

func NewQueryService(logger zerolog.Logger) *Service {
	return &Service{
		logger: logger,
	}
}

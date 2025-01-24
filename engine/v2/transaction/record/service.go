package record

import (
	"context"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction/txmodels"
	"github.com/rs/zerolog"
	"iter"
)

// Service for recording transactions
type Service struct {
	addresses  AddressesService
	outputs    OutputsRepo
	operations OperationsRepo

	broadcaster Broadcaster
	logger      zerolog.Logger
}

// NewService creates a new service for transactions
func NewService(
	logger zerolog.Logger,
	addressesRepo AddressesService,
	outputsRepo OutputsRepo,
	operationsRepo OperationsRepo,
	broadcaster Broadcaster,
) *Service {
	return &Service{
		addresses:   addressesRepo,
		outputs:     outputsRepo,
		operations:  operationsRepo,
		broadcaster: broadcaster,
		logger:      logger,
	}
}

func (s *Service) SaveOperations(ctx context.Context, opRows iter.Seq[*txmodels.NewOperation]) error {
	err := s.operations.SaveAll(ctx, opRows)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to save operations")
	}
	return err
}

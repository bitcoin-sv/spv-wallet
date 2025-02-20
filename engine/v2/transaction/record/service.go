package record

import (
	"context"
	"iter"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/txmodels"
	"github.com/rs/zerolog"
)

// Service for recording transactions
type Service struct {
	addresses    AddressesService
	outputs      OutputsRepo
	operations   OperationsRepo
	transactions TransactionsRepo

	broadcaster     Broadcaster
	paymailNotifier PaymailNotifier
	logger          zerolog.Logger
}

// NewService creates a new service for transactions
func NewService(
	logger zerolog.Logger,
	addressesRepo AddressesService,
	outputsRepo OutputsRepo,
	operationsRepo OperationsRepo,
	transactionsRepo TransactionsRepo,
	broadcaster Broadcaster,
	paymailNotifier PaymailNotifier,
) *Service {
	return &Service{
		addresses:       addressesRepo,
		outputs:         outputsRepo,
		operations:      operationsRepo,
		broadcaster:     broadcaster,
		transactions:    transactionsRepo,
		logger:          logger,
		paymailNotifier: paymailNotifier,
	}
}

// SaveOperations saves all operations along with their transactions
// NOTE: This is crucial for transaction recording
func (s *Service) SaveOperations(ctx context.Context, opRows iter.Seq[*txmodels.NewOperation]) error {
	err := s.operations.SaveAll(ctx, opRows)
	if err != nil {
		return spverrors.Wrapf(err, "failed to save operations")
	}
	return nil
}

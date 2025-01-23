package record

import "github.com/rs/zerolog"

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

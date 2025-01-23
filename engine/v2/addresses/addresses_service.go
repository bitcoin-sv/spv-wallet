package addresses

import (
	"context"
	"iter"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/addresses/addressesmodels"
)

// Service for (P2PKH) addresses
type Service struct {
	addressesRepo AddressRepo
}

// NewService creates a new addresses service
func NewService(addresses AddressRepo) *Service {
	return &Service{
		addressesRepo: addresses,
	}
}

// Create creates a new address
func (s *Service) Create(ctx context.Context, newAddress *addressesmodels.NewAddress) error {
	err := s.addressesRepo.Create(ctx, newAddress)
	if err != nil {
		return spverrors.Wrapf(err, "failed to create address")
	}
	return nil
}

// FindByStringAddresses finds addresses by their string representation
func (s *Service) FindByStringAddresses(ctx context.Context, addresses iter.Seq[string]) ([]addressesmodels.Address, error) {
	results, err := s.addressesRepo.FindByStringAddresses(ctx, addresses)
	if err != nil {
		return nil, spverrors.Wrapf(err, "failed to find addresses")
	}
	return results, nil
}

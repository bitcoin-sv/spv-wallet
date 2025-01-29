package addresses

import (
	"context"
	"iter"

	"github.com/bitcoin-sv/spv-wallet/engine/v2/addresses/addressesmodels"
)

// AddressRepo is an interface for addresses repository.
type AddressRepo interface {
	Create(ctx context.Context, newAddress *addressesmodels.NewAddress) error
	FindByStringAddresses(ctx context.Context, addresses iter.Seq[string]) ([]addressesmodels.Address, error)
}

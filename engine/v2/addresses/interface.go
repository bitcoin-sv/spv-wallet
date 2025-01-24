package addresses

import (
	"context"

	"github.com/bitcoin-sv/spv-wallet/engine/v2/addresses/addressesmodels"
)

// AddressRepo is an interface for addresses repository.
type AddressRepo interface {
	Create(ctx context.Context, newAddress *addressesmodels.NewAddress) error
}

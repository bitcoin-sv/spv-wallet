package addressesmodels

import "github.com/bitcoin-sv/spv-wallet/models/bsv"

// NewAddress is a data for creating a new address.
type NewAddress struct {
	UserID             string
	Address            string
	CustomInstructions bsv.CustomInstructions
}

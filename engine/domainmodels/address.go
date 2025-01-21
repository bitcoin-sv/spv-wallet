package domainmodels

import "github.com/bitcoin-sv/spv-wallet/models/bsv"

// NewAddress represents data for creating a new address
type NewAddress struct {
	UserID             string
	Address            string
	CustomInstructions bsv.CustomInstructions
}

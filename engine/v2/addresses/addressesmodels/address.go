package addressesmodels

import (
	"time"

	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"gorm.io/datatypes"
)

// NewAddress is a data for creating a new address.
type NewAddress struct {
	UserID             string
	Address            string
	CustomInstructions bsv.CustomInstructions
}

// Address represents domain model for P2PKH address.
type Address struct {
	Address string

	CreatedAt time.Time
	UpdatedAt time.Time

	CustomInstructions datatypes.JSONSlice[bsv.CustomInstruction]

	UserID string
}

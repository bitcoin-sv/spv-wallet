package utxo

import (
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/outlines"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/outlines/utxo/internal/sql"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"gorm.io/gorm"
)

// NewSelector creates a new instance of UTXOSelector.
func NewSelector(db *gorm.DB, feeUnit bsv.FeeUnit) outlines.UTXOSelector {
	if db == nil {
		panic("db is required")
	}

	if !feeUnit.IsValid() {
		panic("valid fee unit is required")
	}

	return sql.NewUTXOSelector(db, feeUnit)
}

package utxo

import (
	"context"

	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/outlines"
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

	return &fixMe_ReplaceME{}
}

// FIXME
type fixMe_ReplaceME struct {
}

func (f *fixMe_ReplaceME) Select(ctx context.Context, tx *sdk.Transaction, userID string) ([]*bsv.Outpoint, error) {
	return []*bsv.Outpoint{
		{
			TxID: "a0000000001e1b81dd2c9c0c6cd67f9bdf832e9c2bb12a1d57f30cb6ebbe78d9",
			Vout: 0,
		},
	}, nil
}

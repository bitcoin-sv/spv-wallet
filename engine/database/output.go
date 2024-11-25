package database

import "github.com/bitcoin-sv/spv-wallet/models/bsv"

// Output represents an output of a transaction.
type Output struct {
	TxID       string `gorm:"primaryKey"`
	Vout       uint32 `gorm:"primaryKey"`
	SpendingTX string `gorm:"type:char(64)"`
}

// IsSpent returns true if the output is spent.
func (o *Output) IsSpent() bool {
	return o.SpendingTX != ""
}

// Outpoint returns bsv.Outpoint object which identifies the output.
func (o *Output) Outpoint() *bsv.Outpoint {
	return &bsv.Outpoint{
		TxID: o.TxID,
		Vout: o.Vout,
	}
}

package database

import (
	"time"

	"github.com/bitcoin-sv/spv-wallet/models/bsv"
)

// TrackedOutput represents an output of a transaction.
type TrackedOutput struct {
	TxID       string `gorm:"primaryKey"`
	Vout       uint32 `gorm:"primaryKey"`
	SpendingTX string `gorm:"type:char(64)"`

	UserID string
	User   *User `gorm:"foreignKey:UserID"`

	Satoshis bsv.Satoshis

	CreatedAt time.Time
	UpdatedAt time.Time
}

// IsSpent returns true if the output is spent.
func (o *TrackedOutput) IsSpent() bool {
	return o.SpendingTX != ""
}

// Outpoint returns bsv.Outpoint object which identifies the output.
func (o *TrackedOutput) Outpoint() *bsv.Outpoint {
	return &bsv.Outpoint{
		TxID: o.TxID,
		Vout: o.Vout,
	}
}

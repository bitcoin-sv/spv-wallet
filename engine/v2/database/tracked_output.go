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

	Satoshis bsv.Satoshis

	CreatedAt time.Time
	UpdatedAt time.Time
}

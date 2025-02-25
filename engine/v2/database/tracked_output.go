package database

import (
	"gorm.io/gorm"
	"strings"
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

// AfterFind is a GORM hook that is called after retrieving the record from the database.
func (o *TrackedOutput) AfterFind(_ *gorm.DB) (err error) {
	// Trim left spaces from SpendingTX
	// Because the field is char(64) in the database, it will be padded with spaces.
	// For some reason the value is padded only on postgres,
	// so if changing, make sure to check it on all databases.
	o.SpendingTX = strings.TrimLeft(o.SpendingTX, " ")
	return nil
}

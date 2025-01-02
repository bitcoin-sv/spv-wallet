package database

import (
	dberrors "github.com/bitcoin-sv/spv-wallet/engine/database/errors"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"gorm.io/gorm"
)

// Output represents an output of a transaction.
type TrackedOutput struct {
	TxID       string `gorm:"primaryKey"`
	Vout       uint32 `gorm:"primaryKey"`
	SpendingTX string `gorm:"type:char(64)"`

	UserID                       string `gorm:"-"`
	Bucket                       string `gorm:"-"`
	Satoshis                     uint64 `gorm:"-"`
	UnlockingScriptEstimatedSize uint64 `gorm:"-"`
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

func (o *TrackedOutput) AfterSave(tx *gorm.DB) error {
	if !o.IsSpent() {
		utxo := &UserUtxos{
			TxID: o.TxID,
			Vout: o.Vout,

			UserID:                       o.UserID,
			Bucket:                       o.Bucket,
			Satoshis:                     o.Satoshis,
			UnlockingScriptEstimatedSize: o.UnlockingScriptEstimatedSize,
		}
		if err := tx.Create(utxo).Error; err != nil {
			return dberrors.ErrDBFailed.Wrap(err)
		}
	}

	return nil
}

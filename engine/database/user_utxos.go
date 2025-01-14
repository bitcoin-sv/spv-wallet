package database

import (
	"time"

	"gorm.io/datatypes"
)

// UserUtxos is a table holding user's Unspent Transaction Outputs (UTXOs).
// TODO: It should be renamed to UserUTXO.
type UserUtxos struct {
	UserID                       string `gorm:"primaryKey;uniqueIndex:idx_window,sort:asc,priority:1"`
	TxID                         string `gorm:"primaryKey;uniqueIndex:idx_window,sort:asc,priority:4"`
	Vout                         uint32 `gorm:"primaryKey;uniqueIndex:idx_window,sort:asc,priority:5"`
	Satoshis                     uint64
	UnlockingScriptEstimatedSize uint64
	Bucket                       string    `gorm:"check:chk_not_data_bucket,bucket <> 'data'"`
	CreatedAt                    time.Time `gorm:"uniqueIndex:idx_window,sort:asc,priority:3"`
	TouchedAt                    time.Time `gorm:"uniqueIndex:idx_window,sort:asc,priority:2"`
	CustomInstructions           datatypes.JSONSlice[CustomInstruction]
}

// NewP2PKHUserUTXO creates a new UserUtxos instance for a P2PKH output based on the given output and custom instructions.
func NewP2PKHUserUTXO(output *Output, customInstructions datatypes.JSONSlice[CustomInstruction]) *UserUtxos {
	return &UserUtxos{
		UserID:                       output.UserID,
		TxID:                         output.TxID,
		Vout:                         output.Vout,
		Satoshis:                     uint64(output.Satoshis),
		UnlockingScriptEstimatedSize: 106,
		Bucket:                       "bsv",
		CustomInstructions:           customInstructions,
	}
}

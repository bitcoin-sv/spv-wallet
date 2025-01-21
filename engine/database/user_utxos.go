package database

import (
	"time"

	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"gorm.io/datatypes"
)

// EstimatedInputSizeForP2PKH is the estimated size increase when adding and unlocking P2PKH input to transaction.
// 32 bytes txID
// + 4 bytes vout index
// + 1 byte script length
// + 107 bytes script pub key
// + 4 bytes nSequence
const EstimatedInputSizeForP2PKH = 148

// UserUTXO is a table holding user's Unspent Transaction Outputs (UTXOs).
type UserUTXO struct {
	UserID   string `gorm:"primaryKey;uniqueIndex:idx_window,sort:asc,priority:1"`
	TxID     string `gorm:"primaryKey;uniqueIndex:idx_window,sort:asc,priority:4"`
	Vout     uint32 `gorm:"primaryKey;uniqueIndex:idx_window,sort:asc,priority:5"`
	Satoshis uint64
	// EstimatedInputSize is the estimated size increase when adding and unlocking this UTXO to a transaction.
	EstimatedInputSize uint64
	Bucket             string    `gorm:"check:chk_not_data_bucket,bucket <> 'data'"`
	CreatedAt          time.Time `gorm:"uniqueIndex:idx_window,sort:asc,priority:3"`
	// TouchedAt is the time when the UTXO was last touched (selected for preparing transaction outline) - used for prioritizing UTXO selection.
	TouchedAt time.Time `gorm:"uniqueIndex:idx_window,sort:asc,priority:2"`
	// CustomInstructions is the list of instructions for unlocking given UTXO (it should be understood by client).
	CustomInstructions datatypes.JSONSlice[CustomInstruction]
}

// NewP2PKHUserUTXO creates a new UserUTXO instance for a P2PKH output based on the given output and custom instructions.
func NewP2PKHUserUTXO(output *TrackedOutput, customInstructions datatypes.JSONSlice[CustomInstruction]) *UserUTXO {
	return &UserUTXO{
		UserID:             output.UserID,
		TxID:               output.TxID,
		Vout:               output.Vout,
		Satoshis:           uint64(output.Satoshis),
		EstimatedInputSize: EstimatedInputSizeForP2PKH,
		Bucket:             "bsv",
		CustomInstructions: customInstructions,
	}
}

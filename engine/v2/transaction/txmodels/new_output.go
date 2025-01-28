package txmodels

import (
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
)

// EstimatedInputSizeForP2PKH is the estimated size increase when adding and unlocking P2PKH input to transaction.
// 32 bytes txID
// + 4 bytes vout index
// + 1 byte script length
// + 107 bytes script pub key
// + 4 bytes nSequence
const EstimatedInputSizeForP2PKH = 148

// SpendableUTXO holds the data for spendable outputs.
type SpendableUTXO struct {
	// EstimatedInputSize is the estimated size increase when adding and unlocking this UTXO to a transaction.
	EstimatedInputSize uint64

	// CustomInstructions is the list of instructions for unlocking given UTXO (it should be understood by client).
	CustomInstructions bsv.CustomInstructions
}

// NewOutput holds the data for creating a new output.
type NewOutput struct {
	UserID   string
	TxID     string
	Vout     uint32
	Satoshis bsv.Satoshis
	Bucket   string

	UTXO *SpendableUTXO

	Data []byte
}

// NewOutputForP2PKH creates a new output for P2PKH address.
func NewOutputForP2PKH(outpoint bsv.Outpoint, userID string, satoshis bsv.Satoshis, customInstructions bsv.CustomInstructions) NewOutput {
	return NewOutput{
		UserID:   userID,
		TxID:     outpoint.TxID,
		Vout:     outpoint.Vout,
		Satoshis: satoshis,
		Bucket:   "bsv",
		UTXO: &SpendableUTXO{
			EstimatedInputSize: EstimatedInputSizeForP2PKH,
			CustomInstructions: customInstructions,
		},
	}
}

// NewOutputForData creates a new output for data.
func NewOutputForData(outpoint bsv.Outpoint, userID string, data []byte) NewOutput {
	return NewOutput{
		UserID: userID,
		TxID:   outpoint.TxID,
		Vout:   outpoint.Vout,
		Data:   data,
	}
}

package models

import (
	"time"

	"github.com/bitcoin-sv/spv-wallet/models/common"
)

// UtxoPointer is a pointer model that represents a utxo.
type UtxoPointer struct {
	// TransactionID is a transaction id that utxo points to.
	TransactionID string `json:"transaction_id"`
	// OutputIndex is a output index that utxo points to.
	OutputIndex uint32 `json:"output_index"`
}

// Utxo is a model that represents a utxo.
type Utxo struct {
	// Model is a common model that contains common fields for all models.
	common.Model
	// UtxoPointer is a pointer to a utxo object.
	UtxoPointer `json:",inline"`

	// ID is a utxo id.
	ID string `json:"id"`
	// XpubID is a utxo related xpub id.
	XpubID string `json:"xpub_id"`
	// Satoshis is a utxo satoshis amount.
	Satoshis uint64 `json:"satoshis"`
	// ScriptPubKey is a utxo script pub key.
	ScriptPubKey string `json:"script_pub_key"`
	// Type is a utxo type.
	Type string `json:"type"`
	// DraftID is a utxo transaction related draft id.
	DraftID string `json:"draft_id"`
	// ReservedAt is a time utxo was reserved at.
	ReservedAt time.Time `json:"reserved_at"`
	// SpendingTxID is a spending transaction id - null if not spent yet.
	SpendingTxID string `json:"spending_tx_id"`
	// Transaction is a transaction pointer that utxo points to.
	Transaction *Transaction `json:"transaction"`
}

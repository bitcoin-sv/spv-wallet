package models

import (
	"time"

	"github.com/bitcoin-sv/spv-wallet/models/common"
)

// UtxoPointer is a pointer model that represents a utxo.
type UtxoPointer struct {
	// TransactionID is a transaction id that utxo points to.
	TransactionID string `json:"transaction_id" example:"01d0d0067652f684c6acb3683763f353fce55f6496521c7d99e71e1d27e53f5c"`
	// OutputIndex is a output index that utxo points to.
	OutputIndex uint32 `json:"output_index" example:"0"`
}

// Utxo is a model that represents a utxo.
type Utxo struct {
	// Model is a common model that contains common fields for all models.
	common.OldModel
	// UtxoPointer is a pointer to a utxo object.
	UtxoPointer `json:",inline"`

	// ID is a utxo id which is a hash from transaction id and output index.
	ID string `json:"id" example:"c706a448748d398d542cf4dfad797c9a4b123ebb72dbfb8b27f3d0f1dda99b58"`
	// XpubID is a utxo related xpub id.
	XpubID string `json:"xpub_id" example:"bb8593f85ef8056a77026ad415f02128f3768906de53e9e8bf8749fe2d66cf50"`
	// Satoshis is a utxo satoshis amount.
	Satoshis uint64 `json:"satoshis" example:"100"`
	// ScriptPubKey is a utxo script pub key.
	ScriptPubKey string `json:"script_pub_key" example:"76a91433ba3607a902bc022164bcb6e993f27bd040241c88ac"`
	// Type is a utxo type.
	Type string `json:"type" example:"pubkeyhash"`
	// DraftID is a utxo transaction related draft id.
	DraftID string `json:"draft_id" example:"b356f7fa00cd3f20cce6c21d704cd13e871d28d714a5ebd0532f5a0e0cde63f7"`
	// ReservedAt is a time utxo was reserved at.
	ReservedAt time.Time `json:"reserved_at"  example:"2024-02-26T11:00:28.069911Z"`
	// SpendingTxID is a spending transaction id - null if not spent yet.
	SpendingTxID string `json:"spending_tx_id" example:"01d0d0067652f684c6acb3683763f353fce55f6496521c7d99e71e1d27e53f5c"`
	// Transaction is a transaction pointer that utxo points to.
	Transaction *Transaction `json:"transaction"`
}

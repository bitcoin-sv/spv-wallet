package models

import "github.com/bitcoin-sv/spv-wallet/models/common"

// Transaction is a model that represents a transaction.
type Transaction struct {
	// Model is a common model that contains common fields for all models.
	common.Model
	// ID is a transaction id.
	ID string `json:"id"`
	// Hex is a transaction hex.
	Hex string `json:"hex"`
	// XpubInIDs is a slice of xpub input ids.
	XpubInIDs []string `json:"xpub_in_ids"`
	// XpubOutIDs is a slice of xpub output ids.
	XpubOutIDs []string `json:"xpub_out_ids"`
	// BlockHash is a block hash that transaction is in.
	BlockHash string `json:"block_hash"`
	// BlockHeight is a block height that transaction is in.
	BlockHeight uint64 `json:"block_height"`
	// Fee is a transaction fee.
	Fee uint64 `json:"fee"`
	// NumberOfInputs is a number of transaction inputs.
	NumberOfInputs uint32 `json:"number_of_inputs"`
	// NumberOfOutputs is a number of transaction outputs.
	NumberOfOutputs uint32 `json:"number_of_outputs"`
	// DraftID is a transaction related draft id.
	DraftID string `json:"draft_id"`
	// TotalValue is a total input value.
	TotalValue uint64 `json:"total_value"`
	// OutputValue is a total output value.
	OutputValue int64 `json:"output_value,omitempty"`
	// Outputs represents all spv-wallet-transaction outputs. Will be shown only for admin.
	Outputs map[string]int64 `json:"outputs,omitempty"`
	// Status is a transaction status.
	Status string `json:"status"`
	// TransactionDirection is a transaction direction (inbound/outbound).
	TransactionDirection string `json:"direction"`
}

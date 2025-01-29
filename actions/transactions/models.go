package transactions

import (
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/models/response"
)

// UpdateTransactionRequest is the model used for updating a transaction
type UpdateTransactionRequest struct {
	// Accepts a JSON object for embedding custom metadata, enabling arbitrary additional information to be associated with the resource
	Metadata engine.Metadata `json:"metadata" swaggertype:"object,string" example:"key:value,key2:value2"`
}

// RecordTransaction is the model for recording a transaction
type RecordTransaction struct {
	// Accepts a JSON object for embedding custom metadata, enabling arbitrary additional information to be associated with the resource
	Metadata engine.Metadata `json:"metadata" swaggertype:"object,string" example:"key:value,key2:value2"`
	// Hex of the transaction
	Hex string `json:"hex" example:"0100000002..."`
	// ReferenceID which is a ID of the draft transaction
	ReferenceID string `json:"referenceId" example:"b356f7fa00cd3f20cce6c21d704cd13e871d28d714a5ebd0532f5a0e0cde63f7"`
}

// NewDraftTransaction is the model for creating a new transaction
type NewDraftTransaction struct {
	// Configuration of the transaction
	Config response.TransactionConfig `json:"config"`
	// Accepts a JSON object for embedding custom metadata, enabling arbitrary additional information to be associated with the resource
	Metadata engine.Metadata `json:"metadata" swaggertype:"object,string" example:"key:value,key2:value2"`
}

package transactions

import (
	"github.com/bitcoin-sv/spv-wallet/actions/common"
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
)

// UpdateTransaction is the model for updating a transaction
type UpdateTransaction struct {
	// Id of the transaction which is a hash of the transaction
	ID string `json:"id" example:"01d0d0067652f684c6acb3683763f353fce55f6496521c7d99e71e1d27e53f5c"`
	// Accepts a JSON object for embedding custom metadata, enabling arbitrary additional information to be associated with the resource
	Metadata engine.Metadata `json:"metadata" swaggertype:"object,string" example:"key:value,key2:value2"`
}

// RecordTransaction is the model for recording a transaction
type RecordTransaction struct {
	// Hex of the transaction
	Hex string `json:"hex" example:"0100000002..."`
	// ReferenceID which is a ID of the draft transaction
	ReferenceID string `json:"reference_id" example:"b356f7fa00cd3f20cce6c21d704cd13e871d28d714a5ebd0532f5a0e0cde63f7"`
	// Accepts a JSON object for embedding custom metadata, enabling arbitrary additional information to be associated with the resource
	Metadata engine.Metadata `json:"metadata" swaggertype:"object,string" example:"key:value,key2:value2"`
}

// NewTransaction is the model for creating a new transaction
type NewTransaction struct {
	// Configuration of the transaction
	Config models.TransactionConfig `json:"config"`
	// Accepts a JSON object for embedding custom metadata, enabling arbitrary additional information to be associated with the resource
	Metadata engine.Metadata `json:"metadata" swaggertype:"object,string" example:"key:value,key2:value2"`
}

// SearchTransactions is a model for handling searching with filters and metadata
type SearchTransactions = common.SearchModel[filter.TransactionFilter]

// CountTransactions is a model for handling counting filtered transactions
type CountTransactions = common.ConditionsModel[filter.TransactionFilter]

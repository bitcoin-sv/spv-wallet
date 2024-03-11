package transactions

import (
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/models"
)

// UpdateTransaction is the model for updating a transaction
type UpdateTransaction struct {
	// Id of the transaction which is a hash of the transaction
	ID string `json:"id"`
	// Accepts a JSON object for embedding custom metadata, enabling arbitrary additional information to be associated with the resource
	Metadata engine.Metadata `json:"metadata"`
}

// RecordTransaction is the model for recording a transaction
type RecordTransaction struct {
	// Hex of the transaction
	Hex string `json:"hex"`
	// ReferenceID which is a ID of the draft transaction
	ReferenceID string `json:"reference_id"`
	// Accepts a JSON object for embedding custom metadata, enabling arbitrary additional information to be associated with the resource
	Metadata engine.Metadata `json:"metadata"`
}

// NewTransaction is the model for creating a new transaction
type NewTransaction struct {
	// Configuration of the transaction
	Config models.TransactionConfig `json:"config"`
	// Accepts a JSON object for embedding custom metadata, enabling arbitrary additional information to be associated with the resource
	Metadata engine.Metadata `json:"metadata"`
}

package transactions

import (
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/models"
)

// UpdateTransaction is the model for updating a transaction
type UpdateTransaction struct {
	ID       string          `json:"id"`
	Metadata engine.Metadata `json:"metadata"`
}

// RecordTransaction is the model for recording a transaction
type RecordTransaction struct {
	Hex         string          `json:"hex"`
	ReferenceID string          `json:"reference_id"`
	Metadata    engine.Metadata `json:"metadata"`
}

// NewTransaction is the model for creating a new transaction
type NewTransaction struct {
	Config   models.TransactionConfig `json:"config"`
	Metadata engine.Metadata          `json:"metadata"`
}

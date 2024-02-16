package transactions

import (
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/models"
)

type UpdateTransaction struct {
	Id       string          `json:"id"`
	Metadata engine.Metadata `json:"metadata"`
}

type RecordTransaction struct {
	Hex         string          `json:"hex"`
	ReferenceId string          `json:"reference_id"`
	Metadata    engine.Metadata `json:"metadata"`
}

type NewTransaction struct {
	Config   models.TransactionConfig `json:"config"`
	Metadata engine.Metadata          `json:"metadata"`
}

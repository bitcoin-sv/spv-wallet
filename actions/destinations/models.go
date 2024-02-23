package destinations

import (
	"github.com/bitcoin-sv/spv-wallet/engine"
)

// CreateDestination is the model for creating a destination
type CreateDestination struct {
	Metadata engine.Metadata `json:"metadata"`
}

// UpdateDestination is the model for updating a destination
type UpdateDestination struct {
	ID            string          `json:"id"`
	Address       string          `json:"address"`
	LockingScript string          `json:"locking_script"`
	Metadata      engine.Metadata `json:"metadata"`
}

package destinations

import (
	"github.com/bitcoin-sv/spv-wallet/engine"
)

// CreateDestination is the model for creating a destination
type CreateDestination struct {
	// Accepts a JSON object for embedding custom metadata, enabling arbitrary additional information to be associated with the resource
	Metadata engine.Metadata `json:"metadata"`
}

// UpdateDestination is the model for updating a destination
type UpdateDestination struct {
	// ID of the destination which is the hash of the LockingScript
	ID string `json:"id"`
	// Address of the destination
	Address string `json:"address"`
	// LockingScript of the destination
	LockingScript string `json:"locking_script"`
	// Accepts a JSON object for embedding custom metadata, enabling arbitrary additional information to be associated with the resource
	Metadata engine.Metadata `json:"metadata"`
}

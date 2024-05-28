package destinations

import (
	"github.com/bitcoin-sv/spv-wallet/engine"
)

// CreateDestination is the model for creating a destination
type CreateDestination struct {
	// Accepts a JSON object for embedding custom metadata, enabling arbitrary additional information to be associated with the resource
	Metadata engine.Metadata `json:"metadata" swaggertype:"object,string" example:"key:value,key2:value2"`
}

// UpdateDestination is the model for updating a destination
type UpdateDestination struct {
	// Accepts a JSON object for embedding custom metadata, enabling arbitrary additional information to be associated with the resource
	Metadata engine.Metadata `json:"metadata" swaggertype:"object,string" example:"key:value,key2:value2"`
	// ID of the destination which is the hash of the LockingScript
	ID string `json:"id" example:"82a5d848f997819a478b05fb713208d7f3aa66da5ba00953b9845fb1701f9b98"`
	// Address of the destination
	Address string `json:"address" example:"1CDUf7CKu8ocTTkhcYUbq75t14Ft168K65"`
	// LockingScript of the destination
	LockingScript string `json:"locking_script" example:"76a9147b05764a97f3b4b981471492aa703b188e45979b88ac"`
}

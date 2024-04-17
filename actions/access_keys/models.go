package accesskeys

import "github.com/bitcoin-sv/spv-wallet/engine"

// CreateAccessKey is the model for creating an access key
type CreateAccessKey struct {
	// Accepts a JSON object for embedding custom metadata, enabling arbitrary additional information to be associated with the resource
	Metadata engine.Metadata `json:"metadata" swaggertype:"object,string" example:"key:value,key2:value2"`
}

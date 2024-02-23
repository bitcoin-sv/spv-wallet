package accesskeys

import "github.com/bitcoin-sv/spv-wallet/engine"

// CreateAccessKey is the model for creating an access key
type CreateAccessKey struct {
	Metadata engine.Metadata `json:"metadata"`
}

package contacts

import "github.com/bitcoin-sv/spv-wallet/engine"

// CreateContact is the model for creating a contact
type CreateContact struct {
	Metadata engine.Metadata `json:"metadata"`
}

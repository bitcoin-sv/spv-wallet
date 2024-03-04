package contacts

import "github.com/bitcoin-sv/spv-wallet/engine"

// CreateContact is the model for creating a contact
type CreateContact struct {
	FullName string          `json:"fullName"`
	Paymail  string          `json:"paymail"`
	Metadata engine.Metadata `json:"metadata"`

	RequesterFullName string `json:"requesterFullName"`
	RequesterPaymail  string `json:"requesterPaymail"`
}

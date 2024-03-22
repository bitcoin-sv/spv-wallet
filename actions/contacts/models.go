package contacts

import (
	"errors"

	"github.com/bitcoin-sv/spv-wallet/engine"
)

// UpsertContact is the model for creating a contact
type UpsertContact struct {
	FullName string          `json:"fullName"`
	Metadata engine.Metadata `json:"metadata" swaggertype:"object,string" example:"key:value,key2:value2"`

	// optional
	RequesterPaymail string `json:"requesterPaymail"`
}

func (p *UpsertContact) validate() error {
	if p.FullName == "" {
		return errors.New("fullName is required")
	}

	return nil
}

// UpdateContact is the model for updating a contact
type UpdateContact struct {
	XPubID   string          `json:"xpub_id"`
	FullName string          `json:"full_name"`
	Paymail  string          `json:"paymail"`
	PubKey   string          `json:"pubKey"`
	Status   string          `json:"status"`
	Metadata engine.Metadata `json:"metadata"`
}

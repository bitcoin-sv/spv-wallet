package contacts

import (
	"errors"

	"github.com/bitcoin-sv/spv-wallet/engine"
)

// CreateContact is the model for creating a contact
type CreateContact struct {
	FullName string          `json:"fullName"`
	Paymail  string          `json:"paymail"`
	Metadata engine.Metadata `json:"metadata"`

	RequesterFullName string `json:"requesterFullName"`
	RequesterPaymail  string `json:"requesterPaymail"`
}

func (p *CreateContact) validate() error {
	if p.FullName == "" {
		return errors.New("fullName is required")
	}

	if p.Paymail == "" {
		return errors.New("paymail is required")
	}

	if p.RequesterFullName == "" {
		return errors.New("requesterFullName is required")
	}

	if p.RequesterPaymail == "" {
		return errors.New("requesterPaymail is required")
	}

	return nil
}

package contacts

import (
	"errors"

	"github.com/bitcoin-sv/spv-wallet/engine"
)

// UpsertContact is the model for creating a contact
type UpsertContact struct {
	FullName string          `json:"fullName"`
	Metadata engine.Metadata `json:"metadata"`
}

func (p *UpsertContact) validate() error {
	if p.FullName == "" {
		return errors.New("fullName is required")
	}

	return nil
}

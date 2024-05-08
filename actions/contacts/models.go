package contacts

import (
	"errors"

	"github.com/bitcoin-sv/spv-wallet/actions/common"
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
)

// UpsertContact is the model for creating a contact
type UpsertContact struct {
	// The complete name of the contact, including first name, middle name (if applicable), and last name.
	FullName string `json:"fullName"`
	// Accepts a JSON object for embedding custom metadata, enabling arbitrary additional information to be associated with the resource
	Metadata engine.Metadata `json:"metadata" swaggertype:"object,string" example:"key:value,key2:value2"`
	// Optional paymail address owned by the user to bind the contact to. It is required in case if user has multiple paymail addresses
	RequesterPaymail string `json:"requesterPaymail"`
}

func (p *UpsertContact) validate() error {
	if p.FullName == "" {
		return errors.New("fullName is required")
	}

	return nil
}

// SearchContacts is a model for handling searching with filters and metadata
type SearchContacts = common.SearchModel[filter.ContactFilter]

// CountContacts is a model for handling counting filtered contacts
type CountContacts = common.ConditionsModel[filter.ContactFilter]

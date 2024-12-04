package contacts

import (
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
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

// ContactData is a type for contact data
type ContactData struct {
	// OwnerXpubID is a xpub id related to contact.
	OwnerXpubID string `json:"xpub_id"`
	// FullName is name which could be shown instead of whole paymail address.
	FullName string `json:"full_name"`
	// Paymail is a paymail address related to contact.
	Paymail string `json:"paymail"`
}

func (p *UpsertContact) validate() error {
	if p.FullName == "" {
		return spverrors.ErrMissingContactFullName
	}

	return nil
}

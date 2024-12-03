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
	OwnerXpubID string `json:"xpub_id" toml:"xpub_id" yaml:"xpub_id" gorm:"column:xpub_id;<-:create;type:char(64);foreignKey:XpubID;reference:ID;index;comment:This is the related xPub"`
	// FullName is name which could be shown instead of whole paymail address.
	FullName string `json:"full_name" toml:"full_name" yaml:"full_name" gorm:"<-create;comment:This is the contact's full name"`
	// Paymail is a paymail address related to contact.
	Paymail string `json:"paymail" toml:"paymail" yaml:"paymail" gorm:"<-create;comment:This is the paymail address alias@domain.com"`
}

func (p *UpsertContact) validate() error {
	if p.FullName == "" {
		return spverrors.ErrMissingContactFullName
	}

	return nil
}

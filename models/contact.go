package models

import (
	"github.com/bitcoin-sv/spv-wallet/models/common"
	"github.com/bitcoin-sv/spv-wallet/models/response"
)

// CreateContactResponse is a model for response on contact creation.
type CreateContactResponse struct {
	Contact        *Contact          `json:"contact"`
	AdditionalInfo map[string]string `json:"additionalInfo"`
}

// Contact is a model for contact.
type Contact struct {
	common.Model

	// ID is a unique identifier of contact.
	ID string `json:"id" example:"68af358bde7d8641621c7dd3de1a276c9a62cfa9e2d0740494519f1ba61e2f4a"`
	// FullName is name which could be shown instead of whole paymail address.
	FullName string `json:"fullName" example:"Test User"`
	// Paymail is a paymail address related to contact.
	Paymail string `json:"paymail" example:"test@spv-wallet.com"`
	// PubKey is a public key related to contact (receiver).
	PubKey string `json:"pubKey" example:"xpub661MyMwAqRbcGpZVrSHU..."`
	// Status is a contact's current status.
	Status response.ContactStatus `json:"status" example:"unconfirmed"`
}

// AdminConfirmContactPair is a model for request to confirm contact pair.
type AdminConfirmContactPair struct {
	PaymailA string `json:"paymailA"`
	PaymailB string `json:"paymailB"`
}

// AddAdditionalInfo adds additional information (as key-value map) to the response.
func (m *CreateContactResponse) AddAdditionalInfo(k, v string) {
	if m.AdditionalInfo == nil {
		m.AdditionalInfo = make(map[string]string)
	}

	m.AdditionalInfo[k] = v
}

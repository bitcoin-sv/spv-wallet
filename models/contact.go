package models

import "github.com/bitcoin-sv/spv-wallet/models/common"

type CreateContactResponse struct {
	Contact        *Contact          `json:"contact"`
	AdditionalInfo map[string]string `json:"additionalInfo"`
}

type Contact struct {
	common.Model

	ID string `json:"id"`

	XpubID string `json:"xpubID"`

	FullName string `json:"fullName"`

	Paymail string `json:"paymail"`

	PubKey string `json:"pubKey"`

	Status string `json:"status"`
}

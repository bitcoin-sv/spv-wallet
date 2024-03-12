package models

import "github.com/bitcoin-sv/spv-wallet/models/common"

type ContactStatus string

type Contact struct {
	common.Model

	// ID is the hash of the xpub and paymail
	ID string `json:"id" example:"68af358bde7d8641621c7dd3de1a276c9a62cfa9e2d0740494519f1ba61e2f4a"`
	// XpubID is the contact's xpub related id used to register contact.
	XpubID string `json:"xpubID" example:"bb8593f85ef8056a77026ad415f02128f3768906de53e9e8bf8749fe2d66cf50""`
	// FullName is name which could be shown instead of whole paymail address.
	FullName string `json:"fullName" example:"Test User"`
	// Paymail is a paymail address related to contact.
	Paymail string `json:"paymail" example:"test@spv-wallet.com"`
	// PubKey is a public key related to contact (receiver).
	PubKey string `json:"pubKey" example:"xpub661MyMwAqRbcGpZVrSHU..."`
	// Status is a contact's current status.
	Status ContactStatus `json:"status" example:"confirmed"`
}

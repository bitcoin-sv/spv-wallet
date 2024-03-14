package models

import "github.com/bitcoin-sv/spv-wallet/models/common"

// PaymailAddress is a model that represents a paymail address.
type PaymailAddress struct {
	// Model is a common model that contains common fields for all models.
	common.Model

	// ID is a paymail address id.
	ID string `json:"id"`
	// XpubID is a paymail address's xpub related id used to register paymail address.
	XpubID string `json:"xpub_id"`
	// Alias is a paymail address's alias (first part of paymail).
	Alias string `json:"alias"`
	// Domain is a paymail address's domain (second part of paymail).
	Domain string `json:"domain"`
	// PublicName is a paymail address's public name.
	PublicName string `json:"public_name"`
	// Avatar is a paymail address's avatar.
	Avatar string `json:"avatar"`
}

package models

import "github.com/bitcoin-sv/spv-wallet/models/common"

// PaymailAddress is a model that represents a paymail address.
type PaymailAddress struct {
	// Model is a common model that contains common fields for all models.
	common.OldModel

	// ID is a paymail address id.
	ID string `json:"id" example:"c0ba4a52c89279268476a141be7569200cff2ca4892512b07ca75c25a95c16cd"`
	// XpubID is a paymail address's xpub related id used to register paymail address.
	XpubID string `json:"xpub_id" example:"bb8593f85ef8056a77026ad415f02128f3768906de53e9e8bf8749fe2d66cf50"`
	// Alias is a paymail address's alias (first part of paymail).
	Alias string `json:"alias" example:"test"`
	// Domain is a paymail address's domain (second part of paymail).
	Domain string `json:"domain" example:"spvwallet.com"`
	// PublicName is a paymail address's public name.
	PublicName string `json:"public_name" example:"Test User"`
	// Avatar is a paymail address's avatar.
	Avatar string `json:"avatar" example:"https://spvwallet.com/avatar.png"`
}

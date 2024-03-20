package admin

import (
	"github.com/bitcoin-sv/spv-wallet/engine"
)

// CreatePaymail is the model for creating a paymail
type CreatePaymail struct {
	// The xpub with which the paymail is associated
	Key string `json:"key" example:"xpub661MyMwAqRbcGpZVrSHU..."`
	// The paymail address
	Address string `json:"address" example:"test@spv-wallet.com"`
	// The public name of the paymail
	PublicName string `json:"public_name" example:"Test"`
	// The avatar of the paymail (url address)
	Avatar string `json:"avatar" example:"https://example.com/avatar.png"`
	// Accepts a JSON object for embedding custom metadata, enabling arbitrary additional information to be associated with the resource
	Metadata engine.Metadata `json:"metadata" swaggertype:"object,string" example:"key:value,key2:value2"`
}

// PaymailAddress is the model containing only paymail address used for getting and deleting paymail address
type PaymailAddress struct {
	// The paymail address
	Address string `json:"address" example:"test@spv-wallet.com"`
}

// RecordTransaction is the model for recording a transaction
type RecordTransaction struct {
	// The transaction hex
	Hex string `json:"hex" example:"0100000002..."`
}

// CreateXpub is the model for creating an xpub
type CreateXpub struct {
	// The xpub key
	Key string `json:"key" example:"xpub661MyMwAqRbcGpZVrSHU..."`
	// Accepts a JSON object for embedding custom metadata, enabling arbitrary additional information to be associated with the resource
	Metadata engine.Metadata `json:"metadata" swaggertype:"object,string" example:"key:value,key2:value2"`
}

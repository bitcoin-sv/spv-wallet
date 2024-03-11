package admin

import (
	"github.com/bitcoin-sv/spv-wallet/engine"
)

// CreatePaymail is the model for creating a paymail
type CreatePaymail struct {
	// The xpub with which the paymail is associated
	XpubID string `json:"xpub_id"`
	// The paymail address, example: example@spv-wallet.com
	Address string `json:"address"`
	// The public name of the paymail
	PublicName string `json:"public_name"`
	// The avatar of the paymail (url address)
	Avatar string `json:"avatar"`
	// Accepts a JSON object for embedding custom metadata, enabling arbitrary additional information to be associated with the resource
	Metadata engine.Metadata `json:"metadata"`
}

// PaymailAddress is the model containing only paymail address used for getting and deleting paymail address
type PaymailAddress struct {
	// The paymail address example: example@spv-wallet.com
	Address string `json:"address"`
}

// RecordTransaction is the model for recording a transaction
type RecordTransaction struct {
	// The transaction hex
	Hex string `json:"hex"`
}

// CreateXpub is the model for creating an xpub
type CreateXpub struct {
	// The xpub key
	Key string `json:"key"`
	// Accepts a JSON object for embedding custom metadata, enabling arbitrary additional information to be associated with the resource
	Metadata engine.Metadata `json:"metadata"`
}

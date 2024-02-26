package admin

import (
	"github.com/bitcoin-sv/spv-wallet/engine"
)

// CreatePaymail is the model for creating a paymail
type CreatePaymail struct {
	XpubID     string          `json:"xpub_id"`
	Address    string          `json:"address"`
	PublicName string          `json:"public_name"`
	Avatar     string          `json:"avatar"`
	Metadata   engine.Metadata `json:"metadata"`
}

// PaymailAddress is the model containing only paymail address used for getting and deleting paymail address
type PaymailAddress struct {
	Address string `json:"address"`
}

// RecordTransaction is the model for recording a transaction
type RecordTransaction struct {
	Hex string `json:"hex"`
}

// CreateXpub is the model for creating an xpub
type CreateXpub struct {
	Key      string          `json:"key"`
	Metadata engine.Metadata `json:"metadata"`
}

package admin

import (
	"github.com/bitcoin-sv/spv-wallet/engine"
)

// AdminCreatePaymail is the model for creating a paymail
type AdminCreatePaymail struct {
	XpubID     string          `json:"xpub_id"`
	Address    string          `json:"address"`
	PublicName string          `json:"public_name"`
	Avatar     string          `json:"avatar"`
	Metadata   engine.Metadata `json:"metadata"`
}

// AdminPaymailAddress is the model containing only paymail address used for getting and deleting paymail address
type AdminPaymailAddress struct {
	Address string `json:"address"`
}

// AdminRecordTransaction is the model for recording a transaction
type AdminRecordTransaction struct {
	Hex string `json:"hex"`
}

// AdminCreateXpub is the model for creating an xpub
type AdminCreateXpub struct {
	Key      string          `json:"key"`
	Metadata engine.Metadata `json:"metadata"`
}

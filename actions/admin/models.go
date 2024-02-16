package admin

import (
	"github.com/bitcoin-sv/spv-wallet/engine"
)

type AdminCreatePaymail struct {
	XpubID     string          `json:"xpub_id"`
	Address    string          `json:"address"`
	PublicName string          `json:"public_name"`
	Avatar     string          `json:"avatar"`
	Metadata   engine.Metadata `json:"metadata"`
}

type AdminPaymailAddress struct {
	Address string `json:"address"`
}

type AdminRecordTransaction struct {
	Hex string `json:"hex"`
}

type AdminCreateXpub struct {
	Key      string          `json:"key"`
	Metadata engine.Metadata `json:"metadata"`
}

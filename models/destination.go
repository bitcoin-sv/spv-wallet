package models

import (
	"time"

	"github.com/bitcoin-sv/spv-wallet/models/common"
)

// Destination is a model that represents a destination - registered in a spv-wallet with xpub.
type Destination struct {
	// Model is a common model that contains common fields for all models.
	common.Model

	// ID is a destination id.
	ID string `json:"id"`
	// XpubID is a destination's xpub related id used to register destination.
	XpubID string `json:"xpub_id"`
	// LockingScript is a destination's locking script.
	LockingScript string `json:"locking_script"`
	// Type is a destination's type.
	Type string `json:"type"`

	// Chain is the (chain)/num/ext_derivation_num location of the address related to the xPub.
	Chain uint32 `json:"chain"`
	// Num is the chain/(num)/ext_derivation_num location of the address related to the xPub.
	Num uint32 `json:"num"`
	//PaymailExternalDerivationNum is the chain/num/(ext_derivation_num) location of the address related to the xPub.
	PaymailExternalDerivationNum *uint32 `json:"paymail_external_derivation_num"`

	// Address is a destination's address.
	Address string `json:"address"`
	// DraftID is a destination's draft id.
	DraftID string `json:"draft_id"`
	// Monitor is a time when destination was monitored.
	Monitor time.Time `json:"monitor"`
}

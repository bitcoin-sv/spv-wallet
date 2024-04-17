package models

import (
	"github.com/bitcoin-sv/spv-wallet/models/common"
)

// Destination is a model that represents a destination - registered in a spv-wallet with xpub.
type Destination struct {
	// Model is a common model that contains common fields for all models.
	common.Model

	// ID is a destination id which is the hash of the LockingScript.
	ID string `json:"id" example:"82a5d848f997819a478b05fb713208d7f3aa66da5ba00953b9845fb1701f9b98"`
	// XpubID is a destination's xpub related id used to register destination.
	XpubID string `json:"xpub_id" example:"bb8593f85ef8056a77026ad415f02128f3768906de53e9e8bf8749fe2d66cf50"`
	// LockingScript is a destination's locking script.
	LockingScript string `json:"locking_script" example:"76a9147b05764a97f3b4b981471492aa703b188e45979b88ac"`
	// Type is a destination's type.
	Type string `json:"type" example:"pubkeyhash"`
	// Chain is a destination's chain representation.
	Chain uint32 `json:"chain" example:"0"`
	// Num is a destination's num representation.
	Num uint32 `json:"num" example:"0"`
	// PaymailExternalDerivationNum is the chain/num/(ext_derivation_num) location of the address related to the xPub.
	PaymailExternalDerivationNum *uint32 `json:"paymail_external_derivation_num" example:"0"`
	// Address is a destination's address.
	Address string `json:"address" example:"1CDUf7CKu8ocTTkhcYUbq75t14Ft168K65"`
	// DraftID is a destination's draft id.
	DraftID string `json:"draft_id" example:"b356f7fa00cd3f20cce6c21d704cd13e871d28d714a5ebd0532f5a0e0cde63f7"`
}

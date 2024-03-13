package models

import "github.com/bitcoin-sv/spv-wallet/models/common"

// Xpub is a model that represents a xpub.
type Xpub struct {
	// Model is a common model that contains common fields for all models.
	common.Model

	// ID is a xpub id.
	ID string `json:"id"`
	// CurrentBalance is a xpub's current balance.
	CurrentBalance uint64 `json:"current_balance"`
	// NextInternalNum is the index derivation number use to generate NEXT internal xPub (internal xPub are used for change destinations).
	NextInternalNum uint32 `json:"next_internal_num"`
	// NextExternalNum is the index derivation number use to generate NEXT external xPub (external xPub are used for address destinations).
	NextExternalNum uint32 `json:"next_external_num"`
}

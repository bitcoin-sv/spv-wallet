package models

import "github.com/bitcoin-sv/spv-wallet/models/common"

// Xpub is a model that represents a xpub.
type Xpub struct {
	// Model is a common model that contains common fields for all models.
	common.Model

	// ID is a hash of the xpub.
	ID string `json:"id" example:"bb8593f85ef8056a77026ad415f02128f3768906de53e9e8bf8749fe2d66cf50"`
	// CurrentBalance is a xpub's current balance.
	CurrentBalance uint64 `json:"current_balance" example:"1234"`
	// NextInternalNum is a next internal num.
	NextInternalNum uint32 `json:"next_internal_num" example:"0"`
	// NextExternalNum is a next external num.
	NextExternalNum uint32 `json:"next_external_num" example:"0"`
}

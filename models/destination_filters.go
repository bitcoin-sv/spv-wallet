package models

import (
	"github.com/bitcoin-sv/spv-wallet/models/common"
)

type DestinationFilters struct {
	// LockingScript is a destination's locking script.
	LockingScript *string `json:"locking_script,omitempty" example:"76a9147b05764a97f3b4b981471492aa703b188e45979b88ac"`
	// Address is a destination's address.
	Address *string `json:"address,omitempty" example:"1CDUf7CKu8ocTTkhcYUbq75t14Ft168K65"`
	// DraftID is a destination's draft id.
	DraftID *string `json:"draft_id,omitempty" example:"b356f7fa00cd3f20cce6c21d704cd13e871d28d714a5ebd0532f5a0e0cde63f7"`
	// IncludeDeleted is a flag which includes deleted destinations.
	IncludeDeleted *bool `json:"include_deleted,omitempty" example:"true"`
	// Metadata is a metadata map of outer model.
	Metadata *map[string]interface{} `json:"metadata,omitempty" swaggertype:"object,string" example:"key:value,key2:value2"`
	// CreatedRange is a filter for destinations created within a specific time range.
	CreatedRange *common.TimeRange `json:"created_range,omitempty" example:"from:2024-02-26T11:01:28.069911,to:2025-02-26T11:01:28.069911"`
	// UpdatedRange is a filter for destinations updated within a specific time range.
	UpdatedRange *common.TimeRange `json:"updated_range,omitempty" example:"from:2024-02-26T11:01:28.069911,to:2025-02-26T11:01:28.069911"`
}

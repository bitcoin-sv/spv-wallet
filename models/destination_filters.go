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
	CreatedRange *common.TimeRange `json:"created_range,omitempty" swaggertype:"object,string" example:"from:2024-02-26T11:01:28.069911,to:2025-02-26T11:01:28.069911"`
	// UpdatedRange is a filter for destinations updated within a specific time range.
	UpdatedRange *common.TimeRange `json:"updated_range,omitempty" swaggertype:"object,string" example:"from:2024-02-26T11:01:28.069911,to:2025-02-26T11:01:28.069911"`
}

// NewDestinationFilters Constructor function to create a new instance of DestinationFilters with default values
func NewDestinationFilters() *DestinationFilters {
	includeDeleted := false
	return &DestinationFilters{
		IncludeDeleted: &includeDeleted,
	}
}

func (d *DestinationFilters) ToConditions() map[string]interface{} {
	conditions := map[string]interface{}{}

	if d.LockingScript != nil {
		addToConditions(conditions, "locking_script", d.LockingScript)
	}

	if d.Address != nil {
		addToConditions(conditions, "address", d.Address)
	}

	if d.DraftID != nil {
		addToConditions(conditions, "draft_id", d.DraftID)
	}

	if d.Metadata != nil {
		addToConditions(conditions, "metadata", d.Metadata)
	}

	if d.CreatedRange != nil {
		//addToConditions(conditions, "created_at", map[string]interface{}{
		//	"$gte": d.CreatedRange.From,
		//	"$lte": d.CreatedRange.To,
		//})
		addToConditions(conditions, "createdAt", map[string]interface{}{
			"$gte": d.CreatedRange.From,
			"$lte": d.CreatedRange.To,
		})
	}

	if d.UpdatedRange != nil {
		addToConditions(conditions, "updated_at", map[string]interface{}{
			"$gte": d.UpdatedRange.From,
			"$lte": d.UpdatedRange.To,
		})
	}

	//if d.IncludeDeleted != nil {
	//	addToConditions(conditions, "deleted_at", map[string]interface{}{
	//		"$exists": !*d.IncludeDeleted,
	//	})
	//}

	//if d.IncludeDeleted != nil {
	//	addToConditions(conditions, "deleted_at.Valid", map[string]interface{}{
	//		"$gt": "0001-01-01T00:00:00.069911Z",
	//	})
	//}

	if d.IncludeDeleted != nil {
		// TODO: need to check this
		addToConditions(conditions, "deleted_at", nil)

	}
	return conditions

}

func addToConditions(conditions map[string]interface{}, key string, value interface{}) {
	if value != nil {
		conditions[key] = value
	}
}

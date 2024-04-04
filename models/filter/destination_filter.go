package filter

// DestinationFilters is a struct for handling request parameters for destination search requests
type DestinationFilters struct {
	ModelFilter   `json:",inline"`
	LockingScript *string `json:"locking_script,omitempty" example:"76a9147b05764a97f3b4b981471492aa703b188e45979b88ac"`
	Address       *string `json:"address,omitempty" example:"1CDUf7CKu8ocTTkhcYUbq75t14Ft168K65"`
	DraftID       *string `json:"draft_id,omitempty" example:"b356f7fa00cd3f20cce6c21d704cd13e871d28d714a5ebd0532f5a0e0cde63f7"`
}

// ToDbConditions converts filter fields to the datastore conditions using gorm naming strategy
func (d *DestinationFilters) ToDbConditions() map[string]interface{} {
	conditions := d.ModelFilter.ToDbConditions()

	applyIfNotNil(conditions, "locking_script", d.LockingScript)
	applyIfNotNil(conditions, "address", d.Address)
	applyIfNotNil(conditions, "draft_id", d.DraftID)

	return conditions
}

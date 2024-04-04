package filter

// ModelFilter is a common model filter that contains common fields for all model filters.
type ModelFilter struct {
	IncludeDeleted *bool      `json:"include_deleted,omitempty" example:"true"`
	CreatedRange   *TimeRange `json:"created_range,omitempty" swaggertype:"object,string" example:"from:2024-02-26T11:01:28.069911,to:2025-02-26T11:01:28.069911"`
	UpdatedRange   *TimeRange `json:"updated_range,omitempty" swaggertype:"object,string" example:"from:2024-02-26T11:01:28.069911,to:2025-02-26T11:01:28.069911"`
}

// ToDbConditions converts filter fields to the datastore conditions using gorm naming strategy
func (mf *ModelFilter) ToDbConditions() map[string]interface{} {
	conditions := map[string]interface{}{}

	applyIfNotNilFunc(conditions, "created_at", mf.CreatedRange, func(t *TimeRange) interface{} {
		return t.ToDbConditions()
	})

	applyIfNotNilFunc(conditions, "updated_at", mf.UpdatedRange, func(t *TimeRange) interface{} {
		return t.ToDbConditions()
	})

	if mf.IncludeDeleted == nil || !*mf.IncludeDeleted {
		// if you don't want to include deleted, then add a condition to filter out deleted
		conditions["deleted_at"] = nil
	}
	return conditions
}

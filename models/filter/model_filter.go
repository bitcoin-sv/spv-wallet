package filter

// ModelFilter is a common model filter that contains common fields for all model filters.
type ModelFilter struct {
	IncludeDeleted *bool      `json:"includeDeleted,omitempty" example:"true"`
	CreatedRange   *TimeRange `json:"createdRange,omitempty" swaggertype:"object,string"`
	UpdatedRange   *TimeRange `json:"updatedRange,omitempty" swaggertype:"object,string"`
}

// ToDbConditions converts filter fields to the datastore conditions using gorm naming strategy
func (mf *ModelFilter) ToDbConditions() map[string]interface{} {
	conditions := map[string]interface{}{}

	applyConditionsIfNotNil(conditions, "created_at", mf.CreatedRange.ToDbConditions())
	applyConditionsIfNotNil(conditions, "updated_at", mf.UpdatedRange.ToDbConditions())

	if mf.IncludeDeleted == nil || !*mf.IncludeDeleted {
		// In such cases, we want to filter out deleted items, meaning we only show items
		// where 'deleted_at' is NULL (i.e., items that have not been marked as deleted).
		conditions["deleted_at"] = nil
	}
	return conditions
}

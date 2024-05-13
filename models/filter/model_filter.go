package filter

// ModelFilter is a common model filter that contains common fields for all model filters.
type ModelFilter struct {
	// IncludeDeleted is a flag whether or not to include deleted items in the search results
	IncludeDeleted *bool `json:"includeDeleted,omitempty" swaggertype:"boolean" default:"false" example:"true"`

	// CreatedRange specifies the time range when a record was created.
	CreatedRange *TimeRange `json:"createdRange,omitempty"`

	// UpdatedRange specifies the time range when a record was updated.
	UpdatedRange *TimeRange `json:"updatedRange,omitempty"`
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

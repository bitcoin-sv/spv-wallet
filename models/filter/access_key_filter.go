package filter

// AccessKeyFilter is a struct for handling request parameters for destination search requests
type AccessKeyFilter struct {
	ModelFilter `json:",inline"`
}

// ToDbConditions converts filter fields to the datastore conditions using gorm naming strategy
func (d *AccessKeyFilter) ToDbConditions() map[string]interface{} {
	conditions := d.ModelFilter.ToDbConditions()

	return conditions
}

package filter

// AccessKeyFilter is a struct for handling request parameters for destination search requests
type AccessKeyFilter struct {
	ModelFilter `json:",inline"`

	RevokedRange *TimeRange `json:"revokedRange,omitempty" swaggertype:"object,string"`
}

// ToDbConditions converts filter fields to the datastore conditions using gorm naming strategy
func (d *AccessKeyFilter) ToDbConditions() map[string]interface{} {
	conditions := d.ModelFilter.ToDbConditions()

	// Column names come from the database model, see: /engine/model_access_keys.go
	applyConditionsIfNotNil(conditions, "revoked_at", d.RevokedRange.ToDbConditions())

	return conditions
}

// AdminAccessKeyFilter wraps the AccessKeyFilter providing additional fields for admin access key search requests
type AdminAccessKeyFilter struct {
	AccessKeyFilter `json:",inline"`

	XpubID *string `json:"xpubId,omitempty"`
}

// ToDbConditions converts filter fields to the datastore conditions using gorm naming strategy
func (d *AdminAccessKeyFilter) ToDbConditions() map[string]interface{} {
	conditions := d.AccessKeyFilter.ToDbConditions()

	// Column names come from the database model, see: /engine/model_access_keys.go
	applyIfNotNil(conditions, "xpub_id", d.XpubID)

	return conditions
}

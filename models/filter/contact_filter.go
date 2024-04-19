package filter

// ContactFilter is a struct for handling request parameters for contact search requests
type ContactFilter struct {
	ModelFilter `json:",inline"`
	ID          *string `json:"id"`
	FullName    *string `json:"fullName"`
	Paymail     *string `json:"paymail"`
	PubKey      *string `json:"pubKey"`
	Status      *string `json:"status,omitempty" enums:"unconfirmed,awaiting,confirmed,rejected"`
}

// ToDbConditions converts filter fields to the datastore conditions using gorm naming strategy
func (d *ContactFilter) ToDbConditions() map[string]interface{} {
	conditions := d.ModelFilter.ToDbConditions()

	// Column names come from the database model, see: /engine/model_contact.go
	applyIfNotNil(conditions, "id", d.ID)
	applyIfNotNil(conditions, "full_name", d.FullName)
	applyIfNotNil(conditions, "paymail", d.Paymail)
	applyIfNotNil(conditions, "pub_key", d.PubKey)
	applyIfNotNil(conditions, "status", d.Status)

	return conditions
}

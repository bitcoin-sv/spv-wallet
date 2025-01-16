package filter

// ContactFilter is a struct for handling request parameters for contact search requests
type ContactFilter struct {
	// ModelFilter is a struct for handling typical request parameters for search requests
	ModelFilter `json:",inline"`
	ID          *string `json:"id" example:"ffdbe74e-0700-4710-aac5-611a1f877c7f"`
	FullName    *string `json:"fullName" example:"Alice"`
	Paymail     *string `json:"paymail" example:"alice@example.com"`
	PubKey      *string `json:"pubKey" example:"0334f01ecb971e93db179e6fb320cd1466beb0c1ec6c1c6a37aa6cb02e53d5dd1a"`
	Status      *string `json:"status,omitempty" enums:"unconfirmed,awaiting,confirmed,rejected"`
}

var validContactStatuses = getEnumValues[ContactFilter]("Status")

// ToDbConditions converts filter fields to the datastore conditions using gorm naming strategy
func (d *ContactFilter) ToDbConditions() (map[string]interface{}, error) {
	if d == nil {
		return nil, nil
	}
	conditions := d.ModelFilter.ToDbConditions()

	// Column names come from the database model, see: /engine/model_contact.go
	applyIfNotNil(conditions, "id", d.ID)
	applyIfNotNil(conditions, "full_name", d.FullName)
	applyIfNotNil(conditions, "paymail", d.Paymail)
	applyIfNotNil(conditions, "pub_key", d.PubKey)
	if err := checkAndApplyStrOption(conditions, "status", d.Status, validContactStatuses...); err != nil {
		return nil, err
	}

	return conditions, nil
}

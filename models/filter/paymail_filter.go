package filter

// AdminPaymailFilter is a struct for handling request parameters for paymail_addresses search requests
type AdminPaymailFilter struct {
	ModelFilter `json:",inline"`

	ID         *string `json:"id,omitempty"`
	XpubID     *string `json:"xpubId,omitempty"`
	Alias      *string `json:"alias,omitempty"`
	Domain     *string `json:"domain,omitempty"`
	PublicName *string `json:"publicName,omitempty"`
}

// ToDbConditions converts filter fields to the datastore conditions using gorm naming strategy
func (d *AdminPaymailFilter) ToDbConditions() map[string]interface{} {
	conditions := d.ModelFilter.ToDbConditions()

	// Column names come from the database model, see: /engine/model_paymail_addresses.go
	applyIfNotNil(conditions, "id", d.ID)
	applyIfNotNil(conditions, "xpub_id", d.XpubID)
	applyIfNotNil(conditions, "alias", d.Alias)
	applyIfNotNil(conditions, "domain", d.Domain)
	applyIfNotNil(conditions, "public_name", d.PublicName)

	return conditions
}

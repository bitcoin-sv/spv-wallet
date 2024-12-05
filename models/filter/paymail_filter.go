package filter

// PaymailFilter is a struct for handling request parameters for paymail_addresses search requests
type PaymailFilter struct {
	// ModelFilter is a struct for handling typical request parameters for search requests
	//lint:ignore SA5008 We want to reuse json tags also to mapstructure.
	ModelFilter `json:",inline,squash"`

	ID         *string `json:"id,omitempty" example:"ffb86c103d17d87c15aaf080aab6be5415c9fa885309a79b04c9910e39f2b542"`
	XpubID     *string `json:"xpubId,omitempty" example:"79f90a6bab0a44402fc64828af820e9465645658aea2d138c5205b88e6dabd00"`
	Domain     *string `json:"domain,omitempty" example:"example.com"`
	PublicName *string `json:"publicName,omitempty" example:"Alice"`
}

// ToDbConditions converts filter fields to the datastore conditions using gorm naming strategy
func (d *PaymailFilter) ToDbConditions() map[string]interface{} {
	if d == nil {
		return nil
	}
	conditions := d.ModelFilter.ToDbConditions()

	// Column names come from the database model, see: /engine/model_paymail_addresses.go
	applyIfNotNil(conditions, "id", d.ID)
	applyIfNotNil(conditions, "xpub_id", d.XpubID)
	applyIfNotNil(conditions, "domain", d.Domain)
	applyIfNotNil(conditions, "public_name", d.PublicName)

	return conditions
}

type AdminPaymailFilter struct {
	//lint:ignore SA5008 We want to reuse json tags also to mapstructure.
	PaymailFilter `json:",inline,squash"`

	Alias *string `json:"alias,omitempty" example:"alice"`
}

func (d *AdminPaymailFilter) ToDbConditions() map[string]interface{} {
	if d == nil {
		return nil
	}
	conditions := d.PaymailFilter.ToDbConditions()

	// Column names come from the database model, see: /engine/model_paymail_addresses.go
	applyIfNotNil(conditions, "alias", d.Alias)

	return conditions
}

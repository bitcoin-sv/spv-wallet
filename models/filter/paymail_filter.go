package filter

// PaymailFilter is a struct for handling request parameters for paymail_addresses search requests
type PaymailFilter struct {
	ModelFilter `json:",inline"`

	ID                 *string `json:"id,omitempty"`
	Alias              *string `json:"alias,omitempty"`
	Domain             *string `json:"domain,omitempty"`
	PublicName         *string `json:"publicName,omitempty"`
	Avatar             *string `json:"avatar,omitempty"`
	ExternalXpubKey    *string `json:"externalXpubKey,omitempty"`
	ExternalXpubKeyNum *uint32 `json:"externalXpubNum,omitempty"`
	PubKeyNum          *uint32 `json:"pubkeyNum,omitempty"`
	XpubDerivationSeq  *uint32 `json:"xpubDerivationSeq,omitempty"`
}

// ToDbConditions converts filter fields to the datastore conditions using gorm naming strategy
func (d *PaymailFilter) ToDbConditions() map[string]interface{} {
	conditions := d.ModelFilter.ToDbConditions()

	// Column names come from the database model, see: /engine/model_paymail_addresses.go
	applyIfNotNil(conditions, "id", d.ID)
	applyIfNotNil(conditions, "alias", d.Alias)
	applyIfNotNil(conditions, "domain", d.Domain)
	applyIfNotNil(conditions, "public_name", d.PublicName)
	applyIfNotNil(conditions, "avatar", d.Avatar)
	applyIfNotNil(conditions, "external_xpub_key", d.ExternalXpubKey)
	applyIfNotNil(conditions, "external_xpub_num", d.ExternalXpubKeyNum)
	applyIfNotNil(conditions, "pubkey_num", d.PubKeyNum)
	applyIfNotNil(conditions, "xpub_derivation_seq", d.XpubDerivationSeq)

	return conditions
}

package filter

// XpubFilter is a struct for handling request parameters for utxo search requests
type XpubFilter struct {
	ModelFilter `json:",inline"`

	ID              *string `json:"id,omitempty"`
	CurrentBalance  *uint64 `json:"currentBalance,omitempty"`
	NextInternalNum *uint32 `json:"nextInternalNum,omitempty"`
	NextExternalNum *uint32 `json:"nextExternalNum,omitempty"`
}

// ToDbConditions converts filter fields to the datastore conditions using gorm naming strategy
func (d *XpubFilter) ToDbConditions() map[string]interface{} {
	conditions := d.ModelFilter.ToDbConditions()

	// Column names come from the database model, see: /engine/model_xpubs.go
	applyIfNotNil(conditions, "id", d.ID)
	applyIfNotNil(conditions, "current_balance", d.CurrentBalance)
	applyIfNotNil(conditions, "next_internal_num", d.NextInternalNum)
	applyIfNotNil(conditions, "next_external_num", d.NextExternalNum)

	return conditions
}

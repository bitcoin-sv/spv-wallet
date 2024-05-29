package filter

// XpubFilter is a struct for handling request parameters for utxo search requests
type XpubFilter struct {
	ModelFilter `json:",inline"`

	ID             *string `json:"id,omitempty" example:"00b953624f78004a4c727cd28557475d5233c15f17aef545106639f4d71b712d"`
	CurrentBalance *uint64 `json:"currentBalance,omitempty" example:"1000"`
}

// ToDbConditions converts filter fields to the datastore conditions using gorm naming strategy
func (d *XpubFilter) ToDbConditions() map[string]interface{} {
	if d == nil {
		return nil
	}
	conditions := d.ModelFilter.ToDbConditions()

	// Column names come from the database model, see: /engine/model_xpubs.go
	applyIfNotNil(conditions, "id", d.ID)
	applyIfNotNil(conditions, "current_balance", d.CurrentBalance)

	return conditions
}

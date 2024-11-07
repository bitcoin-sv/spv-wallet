package filter

// AdminTransactionFilter extends TransactionFilter for admin-specific use, including xpubid filtering
type AdminTransactionFilter struct {
	//lint:ignore SA5008 We want to reuse json tags also to mapstructure.
	TransactionFilter `json:",inline,squash"`
	XPubID            *string `json:"xpubid,omitempty" example:"623bc25ce1c0fc510dea72b5ee27b2e70384c099f1f3dce9e73dd987198c3486"`
}

// ToDbConditions converts filter fields to the datastore conditions for admin-specific queries
func (f *AdminTransactionFilter) ToDbConditions() map[string]interface{} {
	conditions := f.TransactionFilter.ToDbConditions()

	return conditions
}
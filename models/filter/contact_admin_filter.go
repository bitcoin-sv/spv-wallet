package filter

// AdminContactFilter extends ContactFilter for admin-specific use, including xpubId filtering
type AdminContactFilter struct {
	//nolint:staticcheck // SA5008 We want to reuse json tags also to mapstructure.
	ContactFilter `json:",inline,squash"`
	XPubID        *string `json:"xpubId,omitempty" example:"623bc25ce1c0fc510dea72b5ee27b2e70384c099f1f3dce9e73dd987198c3486"`
}

// ToDbConditions converts filter fields to the datastore conditions for admin-specific queries
func (f *AdminContactFilter) ToDbConditions() (map[string]interface{}, error) {
	conditions, err := f.ContactFilter.ToDbConditions()
	if err != nil {
		return nil, err
	}

	applyIfNotNil(conditions, "xpub_id", f.XPubID)

	return conditions, nil
}

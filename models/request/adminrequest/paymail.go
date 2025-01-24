package adminrequest

// AddPaymail represents an admin request to add a paymail to a user.
// NOTE: Use either Alias + Domain or the whole paymail Address field
type AddPaymail struct {
	// The paymail address
	Address string `json:"address" example:"test@spv-wallet.com"`

	// Alias of the paymail (before the @)
	Alias string `json:"alias" example:"test"`
	// Domain of the paymail (after the @)
	Domain string `json:"domain" example:"spv-wallet.com"`

	// The public name of the paymail
	PublicName string `json:"publicName" example:"Test"`
	// The avatar of the paymail (url address)
	Avatar string `json:"avatar" example:"https://example.com/avatar.png"`
}

// HasAddress returns true if the paymail address is set
func (a AddPaymail) HasAddress() bool {
	return a.Address != ""
}

// HasAlias returns true if the paymail alias is set
func (a AddPaymail) HasAlias() bool {
	return a.Alias != ""
}

// HasDomain returns true if the paymail domain is set
func (a AddPaymail) HasDomain() bool {
	return a.Domain != ""
}

// AddressEqualsTo returns true if the paymail address is equal to the given string
func (a AddPaymail) AddressEqualsTo(s string) bool {
	return a.Address == s
}

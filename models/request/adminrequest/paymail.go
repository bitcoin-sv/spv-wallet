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

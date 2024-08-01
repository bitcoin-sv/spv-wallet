package response

// SharedConfig with fields which can ba shared across the application components.
// Please be aware NOT to add ANY SENSITIVE information here.
type SharedConfig struct {
	// PaymailDomains is a list of paymail domains handled by spv-wallet.
	PaymailDomains []string `json:"paymailDomains" example:"spv-wallet.com"`
	// ExperimentalFeatures is a map of experimental features handled by spv-wallet.
	ExperimentalFeatures map[string]bool `json:"experimentalFeatures" example:"pikeEnabled:true"`
}

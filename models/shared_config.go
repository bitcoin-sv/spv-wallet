package models

// SharedConfig with fields which can ba shared across the application components.
// Please be aware NOT to add ANY SENSITIVE information here.
type SharedConfig struct {
	PaymilDomains        []string        `json:"paymail_domains"`
	ExperimentalFeatures map[string]bool `json:"experimental_features"`
}

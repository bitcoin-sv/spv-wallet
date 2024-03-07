package models

// SharedConfig with fields which can ba shared across the application components.
// Please be aware NOT to add ANY SENSITIVE information here.
type SharedConfig struct {
	PaymilDomains        []string           `json:"paymail_domains"`
	ExperimentalFeatures ExperimentalConfig `json:"experimental_features"`
}

// ExperimentalConfig represents a feature flag config.
type ExperimentalConfig struct {
	PikeEnabled bool `json:"pike_enabled" mapstructure:"pike_enabled"`
}

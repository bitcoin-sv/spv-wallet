package config

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

const (
	// AuthenticationSchemeXpub is the xpub auth scheme (using xPubs as tokens)
	AuthenticationSchemeXpub = "xpub"
)

// IsAdmin will check if the key is an admin key
func (a *authenticationConfig) IsAdmin(key string) bool {
	return a.AdminKey == key
}

// Validate checks the configuration for specific rules
func (a *authenticationConfig) Validate() error {
	return validation.ValidateStruct(a,
		validation.Field(&a.AdminKey, validation.Required, validation.Length(32, 111)),
		validation.Field(&a.Scheme, validation.Required, validation.In(AuthenticationSchemeXpub)),
	)
}

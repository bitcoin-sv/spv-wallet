package config

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation"
)

const (
	// AuthenticationSchemeXpub is the xpub auth scheme (using xPubs as tokens)
	AuthenticationSchemeXpub = "xpub"
)

// IsAdmin will check if the key is an admin key
func (a *AuthenticationConfig) IsAdmin(key string) bool {
	return a.AdminKey == key
}

// Validate checks the configuration for specific rules
func (a *AuthenticationConfig) Validate() error {
	err := validation.ValidateStruct(a,
		validation.Field(&a.AdminKey, validation.Required, validation.Length(32, 111)),
		validation.Field(&a.Scheme, validation.Required, validation.In(AuthenticationSchemeXpub)),
	)
	if err != nil {
		err = errors.New("error while validating authentication config: " + err.Error())
		return err
	}
	return nil
}

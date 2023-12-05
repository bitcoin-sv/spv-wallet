package config

import (
	"github.com/spf13/viper"
)

const (
	AuthAdminKeyDefault        = "xpub661MyMwAqRbcFaYeQLxmExXvTCjw9jjBRpifkoGggkAitXNNjva4TStLJuYjjEmU4AzXRPGwoECjXo3Rgqg8zQqW6UPVfkKtsrogGBw8xz7"
	AuthRequireSigningDefault  = false
	AuthSchemeDefault          = "xpub"
	AuthSigningDisabledDefault = true
)

func setDefaults() {
	setAuthDefaults()
	// TODO: set next defaults
}

func setAuthDefaults() {
	viper.SetDefault(AuthAdminKey, AuthAdminKeyDefault)
	viper.SetDefault(AuthRequireSigning, AuthRequireSigningDefault)
	viper.SetDefault(AuthScheme, AuthSchemeDefault)
	viper.SetDefault(AuthSigningDisabled, AuthSigningDisabledDefault)
}

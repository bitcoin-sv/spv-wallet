package testabilities

import "github.com/bitcoin-sv/spv-wallet/config"

type ConfigOpts func(*config.AppConfig)

func WithV2() ConfigOpts {
	return func(c *config.AppConfig) {
		c.ExperimentalFeatures.V2 = true
	}
}

func WithDomainValidationDisabled() ConfigOpts {
	return func(c *config.AppConfig) {
		c.Paymail.DomainValidationEnabled = false
	}
}

func WithNotificationsEnabled() ConfigOpts {
	return func(c *config.AppConfig) {
		c.Notifications.Enabled = true
	}
}

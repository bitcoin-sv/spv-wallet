package testabilities

import (
	"github.com/bitcoin-sv/spv-wallet/config"
	"github.com/bitcoin-sv/spv-wallet/engine/testabilities/testmode"
	"os"
)

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

// WithoutCleanup prevents database cleanup after tests
func WithoutCleanup() ConfigOpts {
	return func(c *config.AppConfig) {
		os.Setenv(testmode.EnvSkipCleanup, "true")
	}
}

// WithPostgresContainer configures the test to use a PostgreSQL
func WithPostgresContainer() ConfigOpts {
	return func(c *config.AppConfig) {
		// This will be detected by the engine fixture and trigger usePostgresContainer
		os.Setenv(testmode.EnvDBMode, "postgres")
	}
}

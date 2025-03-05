package testabilities

import (
	"fmt"
	"os"

	"github.com/bitcoin-sv/spv-wallet/config"
	"github.com/bitcoin-sv/spv-wallet/engine/testabilities/testmode"
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
		if err := os.Setenv(testmode.EnvSkipCleanup, "true"); err != nil {
			fmt.Printf("Warning: failed to set environment variable %s: %v", testmode.EnvSkipCleanup, err)
		}
	}
}

// WithPostgresContainer configures the test to use a PostgreSQL
func WithPostgresContainer() ConfigOpts {
	return func(c *config.AppConfig) {
		if err := os.Setenv(testmode.EnvDBMode, "postgres"); err != nil {
			fmt.Printf("Warning: failed to set environment variable %s: %v", testmode.EnvDBMode, err)
		}
	}
}

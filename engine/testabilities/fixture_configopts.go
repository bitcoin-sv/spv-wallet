package testabilities

import "github.com/bitcoin-sv/spv-wallet/config"

type ConfigOpts func(*config.AppConfig)

func WithNewTransactionFlowEnabled() ConfigOpts {
	return func(c *config.AppConfig) {
		c.ExperimentalFeatures.NewTransactionFlowEnabled = true
	}
}

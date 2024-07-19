package config

import (
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/mrz1836/go-cachestore"
)

// Validate checks the configuration for specific rules
func (c *CacheConfig) Validate() error {
	// Valid engine
	if c.Engine == cachestore.Empty || c.Engine == "" {
		return spverrors.Newf("missing a valid cachestore engine")
	}

	if c.Engine == cachestore.Redis {
		if c.Redis == nil {
			return spverrors.Newf("missing redis config")
		} else if len(c.Redis.URL) == 0 {
			return spverrors.Newf("missing redis url")
		}
	}

	return nil
}

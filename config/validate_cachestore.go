package config

import (
	"errors"

	"github.com/mrz1836/go-cachestore"
)

// Validate checks the configuration for specific rules
func (c *CacheConfig) Validate() error {
	// Valid engine
	if c.Engine == cachestore.Empty || c.Engine == "" {
		return errors.New("missing a valid cachestore engine")
	}

	if c.Engine == cachestore.Redis {
		if c.Redis == nil {
			return errors.New("missing redis config")
		} else if len(c.Redis.URL) == 0 {
			return errors.New("missing redis url")
		}
	}

	return nil
}

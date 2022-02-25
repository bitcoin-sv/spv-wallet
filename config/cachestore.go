package config

import (
	"errors"

	"github.com/BuxOrg/bux/cachestore"
)

// Validate checks the configuration for specific rules
func (c *CachestoreConfig) Validate() error {

	// Valid engine
	if c.Engine == cachestore.Empty || c.Engine == "" {
		return errors.New("missing a valid cachestore engine")
	}

	return nil
}

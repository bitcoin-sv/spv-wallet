package config

import (
	"errors"

	"github.com/BuxOrg/bux/datastore"
)

// Validate checks the configuration for specific rules
func (d *datastoreConfig) Validate() error {

	// Valid engine
	if d.Engine == datastore.Empty || d.Engine == "" {
		return errors.New("missing a valid datastore engine")
	}

	return nil
}

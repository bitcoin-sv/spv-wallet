package config

import (
	"errors"

	"github.com/mrz1836/go-datastore"
)

// Validate checks the configuration for specific rules
func (d *DatastoreConfig) Validate() error {

	// Valid engine
	if d.Engine == datastore.Empty || d.Engine == "" {
		return errors.New("missing a valid datastore engine")
	}

	return nil
}

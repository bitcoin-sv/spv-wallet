package config

import (
	"errors"

	"github.com/mrz1836/go-validate"
)

// Validate checks the configuration for specific rules
func (n *NewRelicConfig) Validate() error {

	// If it's enabled
	if n.Enabled {
		if len(n.LicenseKey) != 40 {
			return errors.New("new_relic license key is missing or invalid")
		}
		if len(n.DomainName) <= 0 {
			return errors.New("domain name is missing and required")
		} else if !validate.IsValidDNSName(n.DomainName) {
			return errors.New("domain name is not valid")
		}
	}

	return nil
}

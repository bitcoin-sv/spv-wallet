package config

import (
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/mrz1836/go-validate"
)

// Validate checks the configuration for specific rules
func (n *NewRelicConfig) Validate() error {

	// If it's enabled
	if n.Enabled {
		if len(n.LicenseKey) != 40 {
			return spverrors.Newf("new_relic license key is missing or invalid")
		}
		if len(n.DomainName) <= 0 {
			return spverrors.Newf("domain name is missing and required")
		} else if !validate.IsValidDNSName(n.DomainName) {
			return spverrors.Newf("domain name is not valid")
		}
	}

	return nil
}

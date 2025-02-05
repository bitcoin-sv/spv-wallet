package config

import (
	"fmt"
	"slices"

	configerrors "github.com/bitcoin-sv/spv-wallet/config/errors"
)

// CheckDomain will check if the domain is allowed
func (p *PaymailConfig) CheckDomain(domain string) error {
	if p.DomainValidationEnabled {
		if !slices.Contains(p.Domains, domain) {
			return configerrors.ErrUnsupportedDomain
		}
	}
	return nil
}

// Enabled will return true if the BEEF functionality is enabled
func (b *BeefConfig) Enabled() bool {
	return b != nil && b.UseBeef
}

// GetUserAgent will return the outgoing user agent
func (c *AppConfig) GetUserAgent() string {
	return fmt.Sprintf("%s version %s", applicationName, c.Version)
}

// IsBeefEnabled returns true if the Beef capability will be used for paymail transactions
func (c *AppConfig) IsBeefEnabled() bool {
	return c.Paymail != nil && c.Paymail.Beef.Enabled()
}

package config

import (
	"errors"

	"github.com/mrz1836/go-sanitize"
	"github.com/mrz1836/go-validate"
)

// Validate checks the configuration for specific rules
func (p *PaymailConfig) Validate() error {
	if p == nil {
		return errors.New("paymail config is required")
	}
	if p.Beef.enabled() && p.Beef.BlockHeaderServiceHeaderValidationURL == "" {
		return errors.New("beef_url is required for beef")
	}
	if len(p.Domains) == 0 {
		return errors.New("at least one domain is required for paymail")
	}

	var err error
	for _, domain := range p.Domains {
		domain, err = sanitize.Domain(domain, false, true)
		if err != nil {
			return err
		}
		if !validate.IsValidHost(domain) {
			return errors.New("domain [" + domain + "] is not a valid hostname")
		}
	}

	// Todo: validate the default_from_paymail and default_note values

	return nil
}

package config_test

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet/config"
	"github.com/stretchr/testify/require"
)

func TestValidatePaymailConfig(t *testing.T) {
	t.Parallel()

	validConfigTests := map[string]struct {
		scenario func(cfg *config.AppConfig)
	}{
		"valid beef enabled": {
			scenario: func(cfg *config.AppConfig) {
				cfg.Paymail.Beef.UseBeef = true
				cfg.Paymail.Beef.BlockHeadersServiceHeaderValidationURL = "http://localhost:8080/api/v1/chain/merkleroot/verify"
			},
		},
		"valid multiple paymail domains": {
			scenario: func(cfg *config.AppConfig) {
				cfg.Paymail.Domains = []string{"test.com", "domain.com"}
			},
		},
	}
	for name, test := range validConfigTests {
		t.Run(name, func(t *testing.T) {
			// given:
			cfg := config.GetDefaultAppConfig()

			test.scenario(cfg)

			// when:
			err := cfg.Validate()

			// then:
			require.NoError(t, err)
		})
	}

	invalidConfigTests := map[string]struct {
		scenario func(cfg *config.AppConfig)
	}{
		"invalid for paymail options not set": {
			scenario: func(cfg *config.AppConfig) {
				cfg.Paymail = nil
			},
		},
		"invalid with nil domains": {
			scenario: func(cfg *config.AppConfig) {
				cfg.Paymail.Domains = nil
			},
		},
		"invalid with zero domains": {
			scenario: func(cfg *config.AppConfig) {
				cfg.Paymail.Domains = []string{}
			},
		},
		"invalid with empty domain": {
			scenario: func(cfg *config.AppConfig) {
				cfg.Paymail.Domains = []string{""}
			},
		},
		"invalid with empty domain in the middle": {
			scenario: func(cfg *config.AppConfig) {
				cfg.Paymail.Domains = []string{"test.com", "", "domain.com"}
			},
		},
		"invalid with empty domain at the end": {
			scenario: func(cfg *config.AppConfig) {
				cfg.Paymail.Domains = []string{"test.com", "domain.com", ""}
			},
		},
		"invalid for invalid domain": {
			scenario: func(cfg *config.AppConfig) {
				cfg.Paymail.Domains = []string{"..."}
			},
		},
		"invalid for spaces in domain": {
			scenario: func(cfg *config.AppConfig) {
				cfg.Paymail.Domains = []string{"spaces in domain"}
			},
		},
	}
	for name, test := range invalidConfigTests {
		t.Run(name, func(t *testing.T) {
			// given:
			cfg := config.GetDefaultAppConfig()

			test.scenario(cfg)

			// when:
			err := cfg.Validate()

			// then:
			require.Error(t, err)
		})
	}
}

package config_test

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet/config"
	"github.com/stretchr/testify/require"
)

func TestValidateAuthenticationConfig(t *testing.T) {
	t.Parallel()

	invalidConfigTests := map[string]struct {
		scenario func(cfg *config.AppConfig)
	}{
		"empty scheme": {
			scenario: func(cfg *config.AppConfig) {
				cfg.Authentication.Scheme = ""
			},
		},
		"invalid scheme": {
			scenario: func(cfg *config.AppConfig) {
				cfg.Authentication.Scheme = "invalid"
			},
		},
		"invalid admin key (missing)": {
			scenario: func(cfg *config.AppConfig) {
				cfg.Authentication.AdminKey = ""
			},
		},
		"invalid admin key (to short)": {
			scenario: func(cfg *config.AppConfig) {
				cfg.Authentication.AdminKey = "1234567"
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

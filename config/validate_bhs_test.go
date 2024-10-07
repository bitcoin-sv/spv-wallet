package config_test

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet/config"
	"github.com/stretchr/testify/require"
)

func TestValidateBHSConfig(t *testing.T) {
	t.Parallel()

	validConfigTests := map[string]struct {
		scenario func(cfg *config.AppConfig)
	}{
		"valid with no auth token": {
			scenario: func(cfg *config.AppConfig) {
				cfg.BHS.AuthToken = ""
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
		"return error when url is empty": {
			scenario: func(cfg *config.AppConfig) {
				cfg.BHS.URL = ""
			},
		},
		"return error when config is nil": {
			scenario: func(cfg *config.AppConfig) {
				cfg.BHS = nil
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

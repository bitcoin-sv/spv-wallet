package config_test

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet/config"
	"github.com/stretchr/testify/require"
)

func TestValidateServerConfig(t *testing.T) {
	t.Parallel()

	invalidConfigTests := map[string]struct {
		scenario func(cfg *config.AppConfig)
	}{
		"invalid for missing port": {
			scenario: func(cfg *config.AppConfig) {
				cfg.Server.Port = 0
			},
		},
		"invalid for too high port number": {
			scenario: func(cfg *config.AppConfig) {
				cfg.Server.Port = 1234567890
			},
		},
		"invalid for missing idle timeout": {
			scenario: func(cfg *config.AppConfig) {
				cfg.Server.IdleTimeout = 0
			},
		},
		"invalid for missing read timeout": {
			scenario: func(cfg *config.AppConfig) {
				cfg.Server.ReadTimeout = 0
			},
		},
		"invalid for missing write timeout": {
			scenario: func(cfg *config.AppConfig) {
				cfg.Server.WriteTimeout = 0
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

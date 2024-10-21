package config_test

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet/config"
	"github.com/stretchr/testify/require"
)

func TestValidateFeeUnit(t *testing.T) {
	validConfigTests := map[string]struct {
		scenario func(cfg *config.AppConfig)
	}{
		"Not defined is valid": {
			scenario: func(cfg *config.AppConfig) {
				cfg.CustomFeeUnit = nil
			},
		},
		"Standard": {
			scenario: func(cfg *config.AppConfig) {
				cfg.CustomFeeUnit = &config.FeeUnitConfig{
					Satoshis: 1,
					Bytes:    1000,
				}
			},
		},
		"Zero Satoshi is valid": {
			scenario: func(cfg *config.AppConfig) {
				cfg.CustomFeeUnit = &config.FeeUnitConfig{
					Satoshis: 0,
					Bytes:    1000,
				}
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
		"Empty is not ok": {
			scenario: func(cfg *config.AppConfig) {
				cfg.CustomFeeUnit = &config.FeeUnitConfig{}
			},
		},
		"Negative satoshis": {
			scenario: func(cfg *config.AppConfig) {
				cfg.CustomFeeUnit = &config.FeeUnitConfig{
					Satoshis: -1,
					Bytes:    1000,
				}
			},
		},
		"Zero bytes": {
			scenario: func(cfg *config.AppConfig) {
				cfg.CustomFeeUnit = &config.FeeUnitConfig{
					Satoshis: 1,
					Bytes:    0,
				}
			},
		},
		"Negative bytes": {
			scenario: func(cfg *config.AppConfig) {
				cfg.CustomFeeUnit = &config.FeeUnitConfig{
					Satoshis: 1,
					Bytes:    -1,
				}
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

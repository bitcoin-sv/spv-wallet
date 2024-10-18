package config_test

import (
	"github.com/bitcoin-sv/spv-wallet/config"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestValidateFeeUnit(t *testing.T) {
	tests := map[string]struct {
		customFee  *config.FeeUnitConfig
		shouldFail bool
	}{
		"Not defined is valid": {
			customFee: nil,
		},
		"Standard": {
			customFee: &config.FeeUnitConfig{
				Satoshis: 1,
				Bytes:    1000,
			},
		},
		"Zero Satoshi is valid": {
			customFee: &config.FeeUnitConfig{
				Satoshis: 0,
				Bytes:    1000,
			},
		},
		"Empty is not ok": {
			customFee:  &config.FeeUnitConfig{},
			shouldFail: true,
		},
		"Negative satoshis": {
			customFee: &config.FeeUnitConfig{
				Satoshis: -1,
				Bytes:    1000,
			},
			shouldFail: true,
		},
		"Zero bytes": {
			customFee: &config.FeeUnitConfig{
				Satoshis: 1,
				Bytes:    0,
			},
			shouldFail: true,
		},
		"Negative bytes": {
			customFee: &config.FeeUnitConfig{
				Satoshis: 1,
				Bytes:    -1,
			},
			shouldFail: true,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := test.customFee.Validate()
			if test.shouldFail {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

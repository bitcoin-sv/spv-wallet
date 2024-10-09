package config_test

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet/config"
	"github.com/stretchr/testify/require"
)

func TestValidateAppConfigForDefaultConfig(t *testing.T) {
	t.Parallel()

	// given:
	cfg := config.GetDefaultAppConfig()

	// when:
	err := cfg.Validate()

	// then:
	require.NoError(t, err)
}

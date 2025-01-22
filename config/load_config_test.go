package config_test

import (
	"os"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/config"
	"github.com/bitcoin-sv/spv-wallet/engine/tester"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	t.Run("empty configFilePath", func(t *testing.T) {
		// given
		logger := tester.Logger(t)

		// when
		cfg, err := config.Load("test", logger)

		// then
		assert.NoError(t, err)
		assert.Equal(t, viper.GetString(config.ConfigFilePathKey), config.DefaultConfigFilePath)
		assert.Equal(t, "test", cfg.Version)
	})

	t.Run("custom configFilePath overridden by ENV", func(t *testing.T) {
		// given
		anotherPath := "anotherPath.yml"
		logger := tester.Logger(t)

		// when
		// IMPORTANT! If you need to change the name of this variable, it means you're
		// making backwards incompatible changes. Please inform all SPV Wallet adopters and
		// update your configs on all servers and scripts.
		os.Setenv("SPVWALLET_CONFIG_FILE", anotherPath)
		_, err := config.Load("test", logger)

		// then
		assert.Equal(t, viper.GetString(config.ConfigFilePathKey), anotherPath)
		assert.Error(t, err)

		// cleanup
		os.Unsetenv("SPVWALLET_CONFIG_FILE")
	})
}

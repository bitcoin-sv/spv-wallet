package config

import (
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

// TestLoadConfig will test the method Load()
func TestLoadConfig(t *testing.T) {
	t.Run("empty configFilePath", func(t *testing.T) {
		// when
		_, err := Load()

		// then
		assert.NoError(t, err)
		assert.Equal(t, viper.GetString(ConfigFilePathKey), DefaultConfigFilePath)
	})

	t.Run("custom configFilePath overridden by ENV", func(t *testing.T) {
		// given
		anotherPath := "anotherPath.yml"

		// when
		// IMPORTANT! If you need to change the name of this variable, it means you're
		// making backwards incompatible changes. Please inform all Bux adoptors and
		// update your configs on all servers and scripts.
		os.Setenv("BUX_CONFIG_FILE", anotherPath)
		_, err := Load()

		// then
		assert.Equal(t, viper.GetString(ConfigFilePathKey), anotherPath)
		assert.Error(t, err)

		// cleanup
		os.Unsetenv("BUX_CONFIG_FILE")
	})
}

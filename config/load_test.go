package config

import (
	"os"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

// TestLoadConfig will test the method Load()
func TestLoadConfig(t *testing.T) {
	t.Run("empty configFilePath", func(t *testing.T) {
		// when
		_, err := Load("")

		// then
		assert.NoError(t, err)
		assert.Equal(t, viper.GetString(ConfigFilePathKey), "")
	})

	t.Run("custom configFilePath", func(t *testing.T) {
		// given
		path := "custom/config/file/path.json"

		// when
		_, err := Load(path)

		// then
		assert.Equal(t, viper.GetString(ConfigFilePathKey), path)
		assert.Error(t, err)
	})

	t.Run("custom configFilePath overriden by ENV", func(t *testing.T) {
		// given
		path := "custom/config/file/path.json"
		anotherPath := "anotherPath.json"

		// when
		os.Setenv("BUX_"+strings.ToUpper(ConfigFilePathKey), anotherPath)
		_, err := Load(path)

		// then
		assert.Equal(t, viper.GetString(ConfigFilePathKey), anotherPath)
		assert.Error(t, err)

		// cleanup
		os.Unsetenv("BUX_" + strings.ToUpper(ConfigFilePathKey))
	})
}

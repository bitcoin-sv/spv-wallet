package config

import (
	"testing"

	"github.com/BuxOrg/bux/cachestore"
	"github.com/stretchr/testify/require"
)

// TestCachestoreConfig_Validate will test the method Validate()
func TestCachestoreConfig_Validate(t *testing.T) {
	t.Parallel()

	t.Run("valid datastore config", func(t *testing.T) {
		c := cachestoreConfig{
			Engine: cachestore.MCache,
		}
		require.NotNil(t, c)

		err := c.Validate()
		require.NoError(t, err)
	})

	t.Run("empty datastore", func(t *testing.T) {
		c := cachestoreConfig{
			Engine: cachestore.Empty,
		}
		require.NotNil(t, c)

		err := c.Validate()
		require.Error(t, err)
	})

	t.Run("invalid datastore engine", func(t *testing.T) {
		c := cachestoreConfig{
			Engine: "",
		}
		require.NotNil(t, c)

		err := c.Validate()
		require.Error(t, err)
	})
}

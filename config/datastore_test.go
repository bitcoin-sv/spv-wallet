package config

import (
	"testing"

	"github.com/BuxOrg/bux/datastore"
	"github.com/stretchr/testify/require"
)

// TestDatastoreConfig_Validate will test the method Validate()
func TestDatastoreConfig_Validate(t *testing.T) {
	t.Parallel()

	t.Run("valid datastore config", func(t *testing.T) {
		d := datastoreConfig{
			Engine: datastore.SQLite,
		}
		require.NotNil(t, d)

		err := d.Validate()
		require.NoError(t, err)
	})

	t.Run("empty datastore", func(t *testing.T) {
		d := datastoreConfig{
			Engine: datastore.Empty,
		}
		require.NotNil(t, d)

		err := d.Validate()
		require.Error(t, err)
	})

	t.Run("invalid datastore engine", func(t *testing.T) {
		d := datastoreConfig{
			Engine: "",
		}
		require.NotNil(t, d)

		err := d.Validate()
		require.Error(t, err)
	})
}

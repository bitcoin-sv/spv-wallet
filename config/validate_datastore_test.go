package config

import (
	"testing"

	"github.com/mrz1836/go-datastore"
	"github.com/stretchr/testify/require"
)

// TestDatastoreConfig_Validate will test the method Validate()
func TestDatastoreConfig_Validate(t *testing.T) {
	t.Parallel()

	t.Run("valid datastore config", func(t *testing.T) {
		d := DbConfig{
			Datastore: &DatastoreConfig{Engine: datastore.SQLite},
			Mongo:     &datastore.MongoDBConfig{},
			SQL:       &datastore.SQLConfig{},
			SQLite:    &datastore.SQLiteConfig{},
		}
		require.NotNil(t, d)

		err := d.Validate()
		require.NoError(t, err)
	})

	t.Run("empty datastore", func(t *testing.T) {
		d := DbConfig{
			Datastore: &DatastoreConfig{Engine: datastore.Empty},
			Mongo:     &datastore.MongoDBConfig{},
			SQL:       &datastore.SQLConfig{},
			SQLite:    &datastore.SQLiteConfig{},
		}
		require.NotNil(t, d)

		err := d.Validate()
		require.Error(t, err)
	})

	t.Run("invalid datastore engine", func(t *testing.T) {
		d := DbConfig{
			Datastore: &DatastoreConfig{Engine: ""},
			Mongo:     &datastore.MongoDBConfig{},
			SQL:       &datastore.SQLConfig{},
			SQLite:    &datastore.SQLiteConfig{},
		}
		require.NotNil(t, d)

		err := d.Validate()
		require.Error(t, err)
	})
}

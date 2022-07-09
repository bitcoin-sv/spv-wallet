package config

import (
	"context"
	"os"
	"testing"

	"github.com/BuxOrg/bux/datastore"
	"github.com/mrz1836/go-cachestore"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newTestConfig will make a new test config
func newTestConfig(t *testing.T) (ac *AppConfig) {
	err := os.Setenv(EnvironmentKey, EnvironmentTest)
	require.NoError(t, err)

	ac, err = Load("../")
	require.NoError(t, err)
	require.NotNil(t, ac)
	return
}

func baseTestConfig(t *testing.T) (*AppConfig, *AppServices, *newrelic.Transaction) {
	app := newTestConfig(t)
	require.NotNil(t, app)
	assert.Equal(t, EnvironmentTest, app.Environment)

	services := newTestServices(context.Background(), t, app)
	require.NotNil(t, services)

	txn := services.NewRelic.StartTransaction("test-tx")
	require.NotNil(t, txn)
	return app, services, txn
}

// TestAppConfig_Validate will test the method Validate()
func TestAppConfig_Validate(t *testing.T) {
	t.Parallel()

	t.Run("new test config", func(t *testing.T) {
		app := newTestConfig(t)
		require.NotNil(t, app)
		assert.Equal(t, EnvironmentTest, app.Environment)
	})

	t.Run("validate test config json", func(t *testing.T) {
		app, services, txn := baseTestConfig(t)
		require.NotNil(t, services)
		err := app.Validate(txn)
		assert.NoError(t, err)
	})

	t.Run("authentication - invalid admin_key", func(t *testing.T) {
		app, _, txn := baseTestConfig(t)
		app.Authentication.AdminKey = "12345678"
		err := app.Validate(txn)
		assert.Error(t, err)
	})

	t.Run("authentication - invalid scheme", func(t *testing.T) {
		app, _, txn := baseTestConfig(t)
		app.Authentication.Scheme = "BAD"
		err := app.Validate(txn)
		assert.Error(t, err)
	})

	t.Run("cachestore - invalid engine", func(t *testing.T) {
		app, _, txn := baseTestConfig(t)
		app.Cachestore.Engine = cachestore.Empty
		err := app.Validate(txn)
		assert.Error(t, err)
	})

	t.Run("datastore - invalid engine", func(t *testing.T) {
		app, _, txn := baseTestConfig(t)
		app.Datastore.Engine = datastore.Empty
		err := app.Validate(txn)
		assert.Error(t, err)
	})

	t.Run("new relic - bad license key", func(t *testing.T) {
		app, _, txn := baseTestConfig(t)
		app.NewRelic.Enabled = true
		app.NewRelic.LicenseKey = "1234567"
		err := app.Validate(txn)
		assert.Error(t, err)
	})

	t.Run("new relic - bad domain name", func(t *testing.T) {
		app, _, txn := baseTestConfig(t)
		app.NewRelic.Enabled = true
		app.NewRelic.DomainName = ""
		err := app.Validate(txn)
		assert.Error(t, err)
	})

	t.Run("new relic - invalid domain name", func(t *testing.T) {
		app, _, txn := baseTestConfig(t)
		app.NewRelic.Enabled = true
		app.NewRelic.DomainName = "some domain"
		err := app.Validate(txn)
		assert.Error(t, err)
	})

	t.Run("paymail - no domains", func(t *testing.T) {
		app, _, txn := baseTestConfig(t)
		app.Paymail.Enabled = true
		app.Paymail.Domains = nil
		err := app.Validate(txn)
		assert.Error(t, err)
	})

	t.Run("server - no port", func(t *testing.T) {
		app, _, txn := baseTestConfig(t)
		app.Server.Port = ""
		err := app.Validate(txn)
		assert.Error(t, err)
	})

	t.Run("cachestore - invalid redis url", func(t *testing.T) {
		app, _, txn := baseTestConfig(t)
		app.Cachestore.Engine = cachestore.Redis
		app.Redis.URL = ""
		err := app.Validate(txn)
		assert.Error(t, err)
	})

	t.Run("cachestore - invalid redis config", func(t *testing.T) {
		app, _, txn := baseTestConfig(t)
		app.Cachestore.Engine = cachestore.Redis
		app.Redis = nil
		err := app.Validate(txn)
		assert.Error(t, err)
	})

	t.Run("cachestore - valid freecache", func(t *testing.T) {
		app, _, txn := baseTestConfig(t)
		app.Cachestore.Engine = cachestore.FreeCache
		err := app.Validate(txn)
		assert.NoError(t, err)
	})

	t.Run("datastore - invalid sqlite config", func(t *testing.T) {
		app, _, txn := baseTestConfig(t)
		app.Datastore.Engine = datastore.SQLite
		app.SQLite = nil
		err := app.Validate(txn)
		assert.Error(t, err)
	})

	t.Run("datastore - invalid sql config", func(t *testing.T) {
		app, _, txn := baseTestConfig(t)
		app.Datastore.Engine = datastore.MySQL
		app.SQL = nil
		err := app.Validate(txn)
		assert.Error(t, err)
	})

	t.Run("datastore - invalid sql user", func(t *testing.T) {
		app, _, txn := baseTestConfig(t)
		app.Datastore.Engine = datastore.MySQL
		app.SQL.User = ""
		err := app.Validate(txn)
		assert.Error(t, err)
	})

	t.Run("datastore - invalid sql name", func(t *testing.T) {
		app, _, txn := baseTestConfig(t)
		app.Datastore.Engine = datastore.MySQL
		app.SQL.Name = ""
		err := app.Validate(txn)
		assert.Error(t, err)
	})

	t.Run("datastore - invalid sql host", func(t *testing.T) {
		app, _, txn := baseTestConfig(t)
		app.Datastore.Engine = datastore.MySQL
		app.SQL.Host = ""
		err := app.Validate(txn)
		assert.Error(t, err)
	})

	t.Run("datastore - invalid mongo config", func(t *testing.T) {
		app, _, txn := baseTestConfig(t)
		app.Datastore.Engine = datastore.MongoDB
		app.Mongo = nil
		err := app.Validate(txn)
		assert.Error(t, err)
	})

	t.Run("datastore - invalid mongo uri", func(t *testing.T) {
		app, _, txn := baseTestConfig(t)
		app.Datastore.Engine = datastore.MongoDB
		app.Mongo.URI = ""
		err := app.Validate(txn)
		assert.Error(t, err)
	})

	t.Run("datastore - invalid mongo database name", func(t *testing.T) {
		app, _, txn := baseTestConfig(t)
		app.Datastore.Engine = datastore.MongoDB
		app.Mongo.DatabaseName = ""
		err := app.Validate(txn)
		assert.Error(t, err)
	})

	t.Run("config - invalid env", func(t *testing.T) {
		app, _, txn := baseTestConfig(t)
		app.Environment = ""
		err := app.Validate(txn)
		assert.Error(t, err)
	})

	t.Run("config - invalid working directory", func(t *testing.T) {
		app, _, txn := baseTestConfig(t)
		app.WorkingDirectory = ""
		err := app.Validate(txn)
		assert.Error(t, err)
	})
}

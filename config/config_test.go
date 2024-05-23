package config

import (
	"context"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/engine/datastore"
	"github.com/mrz1836/go-cachestore"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newTestConfig will make a new test config
func newTestConfig(t *testing.T) (ac *AppConfig) {
	nop := zerolog.Nop()
	ac, err := Load(&nop)

	require.NoError(t, err)
	require.NotNil(t, ac)
	return
}

func baseTestConfig(t *testing.T) (*AppConfig, *AppServices) {
	app := newTestConfig(t)
	require.NotNil(t, app)

	services := newTestServices(context.Background(), t, app)
	require.NotNil(t, services)

	return app, services
}

// TestAppConfig_Validate will test the method Validate()
func TestAppConfig_Validate(t *testing.T) {
	t.Parallel()

	t.Run("new test config", func(t *testing.T) {
		app := newTestConfig(t)
		require.NotNil(t, app)
	})

	t.Run("validate test config json", func(t *testing.T) {
		app, services := baseTestConfig(t)
		require.NotNil(t, services)
		err := app.Validate()
		assert.NoError(t, err)
	})

	t.Run("authentication - invalid admin_key", func(t *testing.T) {
		app, _ := baseTestConfig(t)
		app.Authentication.AdminKey = "12345678"
		err := app.Validate()
		assert.Error(t, err)
	})

	t.Run("authentication - invalid scheme", func(t *testing.T) {
		app, _ := baseTestConfig(t)
		app.Authentication.Scheme = "BAD"
		err := app.Validate()
		assert.Error(t, err)
	})

	t.Run("cachestore - invalid engine", func(t *testing.T) {
		app, _ := baseTestConfig(t)
		app.Cache.Engine = cachestore.Empty
		err := app.Validate()
		assert.Error(t, err)
	})

	t.Run("datastore - invalid engine", func(t *testing.T) {
		app, _ := baseTestConfig(t)
		app.Db.Datastore.Engine = datastore.Empty
		err := app.Validate()
		assert.Error(t, err)
	})

	t.Run("new relic - bad license key", func(t *testing.T) {
		app, _ := baseTestConfig(t)
		app.NewRelic.Enabled = true
		app.NewRelic.LicenseKey = "1234567"
		err := app.Validate()
		assert.Error(t, err)
	})

	t.Run("new relic - bad domain name", func(t *testing.T) {
		app, _ := baseTestConfig(t)
		app.NewRelic.Enabled = true
		app.NewRelic.DomainName = ""
		err := app.Validate()
		assert.Error(t, err)
	})

	t.Run("new relic - invalid domain name", func(t *testing.T) {
		app, _ := baseTestConfig(t)
		app.NewRelic.Enabled = true
		app.NewRelic.DomainName = "some domain"
		err := app.Validate()
		assert.Error(t, err)
	})

	t.Run("paymail - no domains", func(t *testing.T) {
		app, _ := baseTestConfig(t)
		app.Paymail.Domains = nil
		err := app.Validate()
		assert.Error(t, err)
	})

	t.Run("server - no port", func(t *testing.T) {
		app, _ := baseTestConfig(t)
		app.Server.Port = 0
		err := app.Validate()
		assert.Error(t, err)
	})

	t.Run("cachestore - invalid redis url", func(t *testing.T) {
		app, _ := baseTestConfig(t)
		app.Cache.Engine = cachestore.Redis
		app.Cache.Redis.URL = ""
		err := app.Validate()
		assert.Error(t, err)
	})

	t.Run("cachestore - invalid redis config", func(t *testing.T) {
		app, _ := baseTestConfig(t)
		app.Cache.Engine = cachestore.Redis
		app.Cache.Redis = nil
		err := app.Validate()
		assert.Error(t, err)
	})

	t.Run("cachestore - valid freecache", func(t *testing.T) {
		app, _ := baseTestConfig(t)
		app.Cache.Engine = cachestore.FreeCache
		err := app.Validate()
		assert.NoError(t, err)
	})

	t.Run("datastore - invalid sqlite config", func(t *testing.T) {
		app, _ := baseTestConfig(t)
		app.Db.Datastore.Engine = datastore.SQLite
		app.Db.SQLite = nil
		err := app.Validate()
		assert.Error(t, err)
	})

	t.Run("datastore - invalid mongo config", func(t *testing.T) {
		app, _ := baseTestConfig(t)
		app.Db.Datastore.Engine = datastore.MongoDB
		app.Db.Mongo = nil
		err := app.Validate()
		assert.Error(t, err)
	})

	t.Run("datastore - invalid mongo uri", func(t *testing.T) {
		app, _ := baseTestConfig(t)
		app.Db.Datastore.Engine = datastore.MongoDB
		app.Db.Mongo.URI = ""
		err := app.Validate()
		assert.Error(t, err)
	})

	t.Run("datastore - invalid mongo database name", func(t *testing.T) {
		app, _ := baseTestConfig(t)
		app.Db.Datastore.Engine = datastore.MongoDB
		app.Db.Mongo.DatabaseName = ""
		err := app.Validate()
		assert.Error(t, err)
	})
}

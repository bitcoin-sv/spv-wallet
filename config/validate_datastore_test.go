package config_test

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet/config"
	"github.com/bitcoin-sv/spv-wallet/engine/datastore"
	"github.com/stretchr/testify/require"
)

func TestValidateDataStoreConfig(t *testing.T) {
	t.Parallel()

	validConfigTests := map[string]struct {
		scenario func(cfg *config.AppConfig)
	}{
		"valid default postgres config": {
			scenario: func(cfg *config.AppConfig) {
				cfg.Db.Datastore.Engine = datastore.PostgreSQL
			},
		},
	}
	for name, test := range validConfigTests {
		t.Run(name, func(t *testing.T) {
			// given:
			cfg := config.GetDefaultAppConfig()

			test.scenario(cfg)

			// when:
			err := cfg.Validate()

			// then:
			require.NoError(t, err)
		})
	}

	invalidConfigTests := map[string]struct {
		scenario func(cfg *config.AppConfig)
	}{
		"invalid when sqlite engine and sqlite config not set": {
			scenario: func(cfg *config.AppConfig) {
				cfg.Db.Datastore.Engine = datastore.SQLite
				cfg.Db.SQLite = nil
			},
		},
		"invalid when empty datastore engine": {
			scenario: func(cfg *config.AppConfig) {
				cfg.Db.Datastore.Engine = datastore.Empty
			},
		},
		"invalid when unknown datastore engine": {
			scenario: func(cfg *config.AppConfig) {
				cfg.Db.Datastore.Engine = ""
			},
		},
		"invalid when sqlite engine and sql config not set": {
			scenario: func(cfg *config.AppConfig) {
				cfg.Db.Datastore.Engine = datastore.PostgreSQL
				cfg.Db.SQL = nil
			},
		},
		"invalid when sqlite engine and host is empty": {
			scenario: func(cfg *config.AppConfig) {
				cfg.Db.Datastore.Engine = datastore.PostgreSQL
				cfg.Db.SQL.Host = ""
			},
		},
		"invalid when sqlite engine and user is empty": {
			scenario: func(cfg *config.AppConfig) {
				cfg.Db.Datastore.Engine = datastore.PostgreSQL
				cfg.Db.SQL.User = ""
			},
		},
		"invalid when sqlite engine and database name is empty": {
			scenario: func(cfg *config.AppConfig) {
				cfg.Db.Datastore.Engine = datastore.PostgreSQL
				cfg.Db.SQL.Name = ""
			},
		},
	}
	for name, test := range invalidConfigTests {
		t.Run(name, func(t *testing.T) {
			// given:
			cfg := config.GetDefaultAppConfig()

			test.scenario(cfg)

			// when:
			err := cfg.Validate()

			// then:
			require.Error(t, err)
		})
	}
}

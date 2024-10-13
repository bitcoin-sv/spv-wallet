package config_test

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet/config"
	"github.com/mrz1836/go-cachestore"
	"github.com/stretchr/testify/require"
)

func TestValidateCacheStoreConfig(t *testing.T) {
	t.Parallel()

	invalidConfigTests := map[string]struct {
		scenario func(cfg *config.AppConfig)
	}{
		"invalid empty cachestore": {
			scenario: func(cfg *config.AppConfig) {
				cfg.Cache.Engine = cachestore.Empty
			},
		},
		"invalid empty string as cachestore engine": {
			scenario: func(cfg *config.AppConfig) {
				cfg.Cache.Engine = ""
			},
		},
		"invalid when cache engine is redis and redis config not provided": {
			scenario: func(cfg *config.AppConfig) {
				cfg.Cache.Engine = cachestore.Redis
				cfg.Cache.Redis = nil
			},
		},
		"invalid when cache engine is redis and redis url is empty": {
			scenario: func(cfg *config.AppConfig) {
				cfg.Cache.Engine = cachestore.Redis
				cfg.Cache.Redis.URL = ""
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

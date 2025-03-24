package config_test

import (
	"fmt"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/config"
	"github.com/stretchr/testify/require"
)

func TestValidateArcConfig(t *testing.T) {
	t.Parallel()

	t.Run("no arc url", func(t *testing.T) {
		// given:
		cfg := config.GetDefaultAppConfig()

		cfg.ARC.URL = ""

		// when:
		err := cfg.Validate()

		// then:
		require.Error(t, err)
	})

	t.Run("if callback is disabled, then empty callback url is valid", func(t *testing.T) {
		// given:
		cfg := config.GetDefaultAppConfig()

		cfg.ARC.Callback.Enabled = false
		cfg.ARC.Callback.Host = ""

		// when:
		err := cfg.Validate()

		// then:
		require.NoError(t, err)
	})

	validCallbackURLTests := []string{
		"http://example.com",
		"https://example.com",
		"http://subdomain.example.com",
		"https://subdomain.example.com",
		"https://subdomain.example.com:3003",
	}
	for _, test := range validCallbackURLTests {
		t.Run(fmt.Sprintf("url %s should be valid callback url", test), func(t *testing.T) {
			// given:
			cfg := config.GetDefaultAppConfig()

			cfg.ARC.Callback.Enabled = true
			cfg.ARC.Callback.Host = test

			// when:
			err := cfg.Validate()

			// then:
			require.NoError(t, err)
		})
	}

	invalidCallbackURLTests := map[string]struct {
		url string
	}{
		"empty callback url is invalid callback url": {
			url: "",
		},
		"external url without schema is invalid callback url": {
			url: "example.com",
		},
		"external url with ftp schema is invalid callback url": {
			url: "ftp://example.com",
		},
		"localhost is invalid callback url": {
			url: "https://localhost",
		},
		"localhost IP is invalid callback url": {
			url: "https://127.0.0.1",
		},
		"local network address is invalid callback url": {
			url: "https://10.0.0.1",
		},
		"url with wrong https schema part (no colon) is invalid callback url": {
			url: "https//example.com",
		},
		"url with wrong http schema part (no colon) is invalid callback url": {
			url: "http//example.com",
		},
	}
	for name, test := range invalidCallbackURLTests {
		t.Run(name, func(t *testing.T) {
			// given:
			cfg := config.GetDefaultAppConfig()

			cfg.ARC.Callback.Enabled = true
			cfg.ARC.Callback.Host = test.url

			// when:
			err := cfg.Validate()

			// then:
			require.Error(t, err)
		})
	}
}

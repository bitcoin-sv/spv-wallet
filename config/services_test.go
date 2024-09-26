package config

import (
	"context"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newTestServices will make a new test services
func newTestServices(ctx context.Context, t *testing.T,
	appConfig *AppConfig,
) *AppServices {
	s, err := appConfig.LoadTestServices(ctx)
	require.NoError(t, err)
	require.NotNil(t, s)
	return s
}

// TestAppServices_CloseAll will test the method CloseAll()
func TestAppServices_CloseAll(t *testing.T) {
	t.Parallel()

	t.Run("no services", func(_ *testing.T) {
		s := new(AppServices)
		s.CloseAll(context.Background())
	})

	t.Run("close all services", func(t *testing.T) {
		ac := newTestConfig(t)
		require.NotNil(t, ac)
		s := newTestServices(context.Background(), t, ac)
		require.NotNil(t, s)
		s.CloseAll(context.Background())

		assert.Nil(t, s.SpvWalletEngine)
		assert.Nil(t, s.NewRelic)
	})
}

// TestAppConfig_GetUserAgent will test the method GetUserAgent()
func TestAppConfig_GetUserAgent(t *testing.T) {
	t.Parallel()

	t.Run("get valid user agent", func(t *testing.T) {
		ac := newTestConfig(t)
		require.NotNil(t, ac)
		agent := ac.GetUserAgent()
		assert.Equal(t, "SPV Wallet "+Version, agent)
	})
}

// TestCallback_HostPattern will test the callback host pattern defined by the regex
func TestCallback_HostPattern(t *testing.T) {
	validURLs := []string{
		"http://example.com",
		"https://example.com",
		"http://subdomain.example.com",
		"https://subdomain.example.com",
		"https://subdomain.example.com:3003",
	}

	invalidURLs := []string{
		"example.com",
		"ftp://example.com",
		"localhost",
		"http//example.com",
		"https//example.com",
		"https://localhost",
		"https://127.0.0.1",
	}

	for _, url := range validURLs {
		if !isValidURL(url) {
			t.Errorf("expected %v to be valid, but it was not", url)
		}
	}

	for _, url := range invalidURLs {
		if isValidURL(url) {
			t.Errorf("expected %v to be invalid, but it was not", url)
		}
	}
}

// TestCallback_ConfigureCallback will test the method configureCallback()
func TestCallback_ConfigureCallback(t *testing.T) {
	tests := []struct {
		appConfig    AppConfig
		name         string
		expectedErr  string
		expectedOpts int
	}{
		{
			appConfig: AppConfig{
				ARC: &ARCConfig{
					Callback: &CallbackConfig{
						Host:    "http://example.com",
						Token:   "",
						Enabled: true,
					},
				},
			},
			name:         "Valid URL with empty token and http",
			expectedErr:  "",
			expectedOpts: 1,
		},
		{
			appConfig: AppConfig{
				ARC: &ARCConfig{
					Callback: &CallbackConfig{
						Host:    "https://example.com",
						Token:   "existingToken",
						Enabled: true,
					},
				},
			},
			name:         "Valid URL with existing token and https",
			expectedErr:  "",
			expectedOpts: 1,
		},
		{
			appConfig: AppConfig{
				ARC: &ARCConfig{
					Callback: &CallbackConfig{
						Host:    "ftp://example.com",
						Token:   "",
						Enabled: true,
					},
				},
			},
			name:         "Invalid URL without http/https",
			expectedErr:  "invalid callback host: ftp://example.com - must be a valid external url - not a localhost",
			expectedOpts: 0,
		},
		{
			appConfig: AppConfig{
				ARC: &ARCConfig{
					Callback: &CallbackConfig{
						Host:    "http://localhost:3003",
						Token:   "",
						Enabled: true,
					},
				},
			},
			name:         "Invalid URL with localhost",
			expectedErr:  "invalid callback host: http://localhost:3003 - must be a valid external url - not a localhost",
			expectedOpts: 0,
		},
		{
			appConfig: AppConfig{
				ARC: &ARCConfig{
					Callback: &CallbackConfig{
						Host:    "http://example.com",
						Token:   "",
						Enabled: false,
					},
				},
			},
			name:         "Callback disabled",
			expectedErr:  "",
			expectedOpts: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var options []engine.ClientOps
			ops, err := configureCallback(options, &tt.appConfig)
			if tt.expectedErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedOpts, len(ops))
		})
	}
}

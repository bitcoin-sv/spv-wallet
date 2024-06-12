package config

import (
	"context"
	"github.com/bitcoin-sv/spv-wallet/engine"
	"testing"

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

func TestCallback_HostPattern(t *testing.T) {
	validURLs := []string{
		"http://example.com",
		"https://example.com",
		"http://subdomain.example.com",
		"https://subdomain.example.com",
	}

	invalidURLs := []string{
		"example.com",
		"ftp://example.com",
		"localhost",
		"http//example.com",
		"https//example.com",
	}

	for _, url := range validURLs {
		if !callbackURLPattern.MatchString(url) {
			t.Errorf("expected %v to be valid, but it was not", url)
		}
	}

	for _, url := range invalidURLs {
		if callbackURLPattern.MatchString(url) {
			t.Errorf("expected %v to be invalid, but it was not", url)
		}
	}
}

func TestCallback_ConfigureCallback(t *testing.T) {
	tests := []struct {
		name         string
		appConfig    AppConfig
		expectedErr  string
		expectedOpts int
	}{
		{
			name: "Valid URL with empty token and http",
			appConfig: AppConfig{
				Nodes: &NodesConfig{
					Callback: &CallbackConfig{
						Enabled: true,
						Host:    "http://example.com",
						Token:   "",
					},
				},
			},
			expectedErr:  "",
			expectedOpts: 1,
		},
		{
			name: "Valid URL with existing token and https",
			appConfig: AppConfig{
				Nodes: &NodesConfig{
					Callback: &CallbackConfig{
						Enabled: true,
						Host:    "https://example.com",
						Token:   "existingToken",
					},
				},
			},
			expectedErr:  "",
			expectedOpts: 1,
		},
		{
			name: "Invalid URL without http/https",
			appConfig: AppConfig{
				Nodes: &NodesConfig{
					Callback: &CallbackConfig{
						Enabled: true,
						Host:    "ftp://example.com",
						Token:   "",
					},
				},
			},
			expectedErr:  "invalid callback host: ftp://example.com - must be a https:// or http:// valid external url",
			expectedOpts: 0,
		},
		{
			name: "Invalid URL with localhost",
			appConfig: AppConfig{
				Nodes: &NodesConfig{
					Callback: &CallbackConfig{
						Enabled: true,
						Host:    "http://localhost:3003",
						Token:   "",
					},
				},
			},
			expectedErr:  "invalid callback host: http://localhost:3003 - must be a valid external url - not a localhost",
			expectedOpts: 0,
		},
		{
			name: "Callback disabled",
			appConfig: AppConfig{
				Nodes: &NodesConfig{
					Callback: &CallbackConfig{
						Enabled: false,
						Host:    "http://example.com",
						Token:   "",
					},
				},
			},
			expectedErr:  "",
			expectedOpts: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			options := []engine.ClientOps{}
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

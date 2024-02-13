package config

import (
	"context"
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

	t.Run("no services", func(t *testing.T) {
		s := new(AppServices)
		s.CloseAll(context.Background())
	})

	t.Run("close all services", func(t *testing.T) {
		ac := newTestConfig(t)
		require.NotNil(t, ac)
		s := newTestServices(context.Background(), t, ac)
		require.NotNil(t, s)
		s.CloseAll(context.Background())

		assert.Nil(t, s.SPV)
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
		assert.Equal(t, "SPV-Wallet "+Version, agent)
	})
}

package config

import (
	"context"
	"net/http"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newTestServices will make a new test services
func newTestServices(ctx context.Context, t *testing.T,
	appConfig *AppConfig) *AppServices {
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

		assert.Nil(t, s.Bux)
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
		assert.Equal(t, "BUX-Server "+EnvironmentTest+" "+Version, agent)
	})
}

// TestAppServices_NewRestyClient will test the method NewRestyClient()
func TestAppServices_NewRestyClient(t *testing.T) {
	t.Parallel()

	t.Run("new client", func(t *testing.T) {
		ac := newTestConfig(t)
		require.NotNil(t, ac)

		s := newTestServices(context.Background(), t, ac)
		require.NotNil(t, s)
		defer s.CloseAll(context.Background())

		r := s.NewRestyClient()
		require.NotNil(t, r)
		assert.IsType(t, &resty.Client{}, r)

		assert.Equal(t, 2, r.RetryCount)
	})
}

// TestAppServices_NewHTTPClient will test the method NewHTTPClient()
func TestAppServices_NewHTTPClient(t *testing.T) {
	t.Parallel()

	t.Run("new client", func(t *testing.T) {
		ac := newTestConfig(t)
		require.NotNil(t, ac)

		s := newTestServices(context.Background(), t, ac)
		require.NotNil(t, s)
		defer s.CloseAll(context.Background())

		h := s.NewHTTPClient()
		require.NotNil(t, h)
		assert.IsType(t, &http.Client{}, h)
	})
}

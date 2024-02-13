package server

import (
	"testing"

	"github.com/BuxOrg/spv-wallet/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewServer will test the method NewServer()
func TestNewServer(t *testing.T) {
	t.Parallel()

	t.Run("empty values", func(t *testing.T) {
		s := NewServer(nil, nil)
		require.NotNil(t, s)
		assert.Nil(t, s.AppConfig)
		assert.Nil(t, s.Services)
		assert.Nil(t, s.Router)
		assert.Nil(t, s.WebServer)
	})

	t.Run("set values", func(t *testing.T) {
		s := NewServer(&config.AppConfig{}, &config.AppServices{})
		require.NotNil(t, s)
		assert.NotNil(t, s.AppConfig)
		assert.NotNil(t, s.Services)
		assert.Nil(t, s.Router)
		assert.Nil(t, s.WebServer)
	})
}

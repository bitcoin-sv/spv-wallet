package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestBHSConfig_Validate will test the method Validate()
func TestBHSConfig_Validate(t *testing.T) {
	t.Parallel()

	t.Run("no auth token", func(t *testing.T) {
		b := BHSConfig{
			AuthToken:  "",
			URL:        "http://localhost:8080",
			APIVersion: "v1",
		}

		err := b.Validate()
		require.NoError(t, err)
	})

	t.Run("no url", func(t *testing.T) {
		b := BHSConfig{
			AuthToken:  "token",
			APIVersion: "v1",
			URL:        "",
		}

		err := b.Validate()
		require.Error(t, err)
	})

	t.Run("no api version", func(t *testing.T) {
		b := BHSConfig{
			AuthToken:  "token",
			URL:        "http://localhost:8080",
			APIVersion: "",
		}

		err := b.Validate()
		require.Error(t, err)
	})

	t.Run("config is nil", func(t *testing.T) {
		var b *BHSConfig

		err := b.Validate()
		require.Error(t, err)
	})

	t.Run("full config", func(t *testing.T) {
		b := BHSConfig{
			AuthToken:  "token",
			URL:        "http://localhost:8080",
			APIVersion: "v1",
		}

		err := b.Validate()
		require.NoError(t, err)
	})
}

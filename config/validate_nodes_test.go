package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewRelicConfig_Validate will test the method Validate()
func TestNodesConfig_Validate(t *testing.T) {
	t.Parallel()

	t.Run("valid default nodes config", func(t *testing.T) {
		n := getARCDefaults()
		assert.NoError(t, n.Validate())
	})

	t.Run("no arc url", func(t *testing.T) {
		n := getARCDefaults()

		n.URL = ""
		assert.Error(t, n.Validate())
	})
}

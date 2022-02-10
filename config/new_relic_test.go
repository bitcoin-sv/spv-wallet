package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewRelicConfig_Validate will test the method Validate()
func TestNewRelicConfig_Validate(t *testing.T) {
	t.Parallel()

	t.Run("valid new relic", func(t *testing.T) {
		n := newRelicConfig{
			Enabled:    true,
			LicenseKey: "1234567890123456789012345678901234567890",
			DomainName: "domain.com",
		}
		assert.NoError(t, n.Validate())
	})

	t.Run("not enabled", func(t *testing.T) {
		n := newRelicConfig{
			Enabled:    false,
			DomainName: "domain.com",
			LicenseKey: "1234567890123456789012345678901234567890",
		}
		assert.NoError(t, n.Validate())
	})

	t.Run("missing domain", func(t *testing.T) {
		n := newRelicConfig{
			Enabled:    true,
			DomainName: "",
			LicenseKey: "1234567890123456789012345678901234567890",
		}
		assert.Error(t, n.Validate())
	})

	t.Run("invalid key", func(t *testing.T) {
		n := newRelicConfig{
			Enabled:    true,
			DomainName: "domain.com",
			LicenseKey: "1234567",
		}
		assert.Error(t, n.Validate())
	})
}

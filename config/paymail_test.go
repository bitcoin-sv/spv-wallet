package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestPaymailConfig_Validate will test the method Validate()
func TestPaymailConfig_Validate(t *testing.T) {
	t.Parallel()

	t.Run("no domains", func(t *testing.T) {
		p := paymailConfig{
			Domains: nil,
			Enabled: true,
		}
		err := p.Validate()
		require.Error(t, err)
	})

	t.Run("zero domains", func(t *testing.T) {
		p := paymailConfig{
			Domains: []string{},
			Enabled: true,
		}
		err := p.Validate()
		require.Error(t, err)
	})

	t.Run("empty domains", func(t *testing.T) {
		p := paymailConfig{
			Domains: []string{""},
			Enabled: true,
		}
		err := p.Validate()
		require.Error(t, err)
	})

	t.Run("invalid hostname", func(t *testing.T) {
		p := paymailConfig{
			Domains: []string{"..."},
			Enabled: true,
		}
		err := p.Validate()
		require.Error(t, err)
	})

	t.Run("spaces in hostname", func(t *testing.T) {
		p := paymailConfig{
			Domains: []string{"spaces in domain"},
			Enabled: true,
		}
		err := p.Validate()
		require.Error(t, err)
	})

	t.Run("valid domains", func(t *testing.T) {
		p := paymailConfig{
			Domains: []string{"test.com", "domain.com"},
			Enabled: true,
		}
		err := p.Validate()
		require.NoError(t, err)
	})

}

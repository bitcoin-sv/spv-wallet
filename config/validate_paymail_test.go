package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestPaymailConfig_Validate will test the method Validate()
func TestPaymailConfig_Validate(t *testing.T) {
	t.Parallel()

	t.Run("no domains", func(t *testing.T) {
		p := PaymailConfig{
			Domains: nil,
		}
		err := p.Validate()
		require.Error(t, err)
	})

	t.Run("zero domains", func(t *testing.T) {
		p := PaymailConfig{
			Domains: []string{},
		}
		err := p.Validate()
		require.Error(t, err)
	})

	t.Run("empty domains", func(t *testing.T) {
		p := PaymailConfig{
			Domains: []string{""},
		}
		err := p.Validate()
		require.Error(t, err)
	})

	t.Run("invalid hostname", func(t *testing.T) {
		p := PaymailConfig{
			Domains: []string{"..."},
		}
		err := p.Validate()
		require.Error(t, err)
	})

	t.Run("spaces in hostname", func(t *testing.T) {
		p := PaymailConfig{
			Domains: []string{"spaces in domain"},
		}
		err := p.Validate()
		require.Error(t, err)
	})

	t.Run("valid domains", func(t *testing.T) {
		p := PaymailConfig{
			Domains: []string{"test.com", "domain.com"},
		}
		err := p.Validate()
		require.NoError(t, err)
	})

	t.Run("invalid beef", func(t *testing.T) {
		p := PaymailConfig{
			Domains: []string{"test.com", "domain.com"},
			Beef: &BeefConfig{
				UseBeef:                                true,
				BlockHeadersServiceHeaderValidationURL: "",
			},
		}
		err := p.Validate()
		require.Error(t, err)
	})

	t.Run("valid beef", func(t *testing.T) {
		p := PaymailConfig{
			Domains: []string{"test.com", "domain.com"},
			Beef: &BeefConfig{
				UseBeef:                                true,
				BlockHeadersServiceHeaderValidationURL: "http://localhost:8080/api/v1/chain/merkleroot/verify",
			},
		}
		err := p.Validate()
		require.NoError(t, err)
	})
}

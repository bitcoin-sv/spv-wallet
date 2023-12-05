package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testAdminKey = "12345678901234567890123456789012"
)

// TestAuthenticationConfig_IsAdmin will test the method IsAdmin()
func TestAuthenticationConfig_IsAdmin(t *testing.T) {
	t.Run("admin valid", func(t *testing.T) {
		a := AuthenticationConfig{
			Scheme:   AuthenticationSchemeXpub,
			AdminKey: testAdminKey,
		}
		assert.Equal(t, true, a.IsAdmin(testAdminKey))
	})

	t.Run("admin invalid", func(t *testing.T) {
		a := AuthenticationConfig{
			Scheme:   AuthenticationSchemeXpub,
			AdminKey: testAdminKey,
		}
		assert.Equal(t, false, a.IsAdmin("invalid"))
	})

}

// TestNewRelicConfig_Validate will test the method Validate()
func TestAuthenticationConfig_Validate(t *testing.T) {
	t.Parallel()

	t.Run("valid scheme and admin key", func(t *testing.T) {
		a := AuthenticationConfig{
			Scheme:   AuthenticationSchemeXpub,
			AdminKey: testAdminKey,
		}
		assert.NoError(t, a.Validate())
	})

	t.Run("empty scheme", func(t *testing.T) {
		a := AuthenticationConfig{
			Scheme:   "",
			AdminKey: testAdminKey,
		}
		assert.Error(t, a.Validate())
	})

	t.Run("invalid scheme", func(t *testing.T) {
		a := AuthenticationConfig{
			Scheme:   "invalid",
			AdminKey: testAdminKey,
		}
		assert.Error(t, a.Validate())
	})

	t.Run("invalid admin key (missing)", func(t *testing.T) {
		a := AuthenticationConfig{
			Scheme:   AuthenticationSchemeXpub,
			AdminKey: "",
		}
		assert.Error(t, a.Validate())
	})

	t.Run("invalid admin key (to short)", func(t *testing.T) {
		a := AuthenticationConfig{
			Scheme:   AuthenticationSchemeXpub,
			AdminKey: "1234567",
		}
		assert.Error(t, a.Validate())
	})
}

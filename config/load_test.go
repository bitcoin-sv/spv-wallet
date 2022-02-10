package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestIsValidEnvironment will test the method isValidEnvironment()
func TestIsValidEnvironment(t *testing.T) {
	t.Run("empty env", func(t *testing.T) {
		valid := isValidEnvironment("")
		assert.Equal(t, false, valid)
	})

	t.Run("unknown env", func(t *testing.T) {
		valid := isValidEnvironment("unknown")
		assert.Equal(t, false, valid)
	})

	t.Run("different case of letters", func(t *testing.T) {
		valid := isValidEnvironment("DEVELOPment")
		assert.Equal(t, true, valid)
	})

	t.Run("valid envs", func(t *testing.T) {
		valid := isValidEnvironment(EnvironmentTest)
		assert.Equal(t, true, valid)

		valid = isValidEnvironment(EnvironmentDevelopment)
		assert.Equal(t, true, valid)

		valid = isValidEnvironment(EnvironmentStaging)
		assert.Equal(t, true, valid)

		valid = isValidEnvironment(EnvironmentProduction)
		assert.Equal(t, true, valid)
	})
}

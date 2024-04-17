package taskmanager

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestFactory_String will test the method String()
func TestFactory_String(t *testing.T) {

	t.Run("test all factories", func(t *testing.T) {
		assert.Equal(t, "empty", FactoryEmpty.String())
		assert.Equal(t, "memory", FactoryMemory.String())
		assert.Equal(t, "redis", FactoryRedis.String())
	})
}

// TestFactory_IsEmpty will test the method IsEmpty()
func TestFactory_IsEmpty(t *testing.T) {
	t.Run("test empty factory", func(t *testing.T) {
		f := FactoryEmpty
		assert.Equal(t, true, f.IsEmpty())
	})

	t.Run("test regular factory", func(t *testing.T) {
		f := FactoryMemory
		assert.Equal(t, false, f.IsEmpty())
	})
}

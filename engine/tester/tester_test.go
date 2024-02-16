package tester

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testDomain          = "domain.com"
	testRedisConnection = "redis://localhost:6379"
)

// TestRandomTablePrefix will test the method RandomTablePrefix()
func TestRandomTablePrefix(t *testing.T) {
	t.Parallel()

	t.Run("valid prefix", func(t *testing.T) {
		prefix := RandomTablePrefix()
		assert.Equal(t, 17, len(prefix))
	})
}

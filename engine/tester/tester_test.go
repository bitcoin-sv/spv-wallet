package tester

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomTablePrefix(t *testing.T) {
	t.Parallel()

	t.Run("valid prefix", func(t *testing.T) {
		prefix := RandomTablePrefix()
		assert.Equal(t, 17, len(prefix))
	})
}

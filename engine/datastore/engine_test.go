package datastore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEngine_String(t *testing.T) {
	t.Run("valid name", func(t *testing.T) {
		assert.Equal(t, "empty", Empty.String())
		assert.Equal(t, "postgresql", PostgreSQL.String())
		assert.Equal(t, "sqlite", SQLite.String())
	})
}

func TestEngine_IsEmpty(t *testing.T) {
	t.Run("actually empty", func(t *testing.T) {
		assert.True(t, Empty.IsEmpty())
	})
}

func TestIsSQLEngine(t *testing.T) {
	t.Run("test sql databases", func(t *testing.T) {
		assert.True(t, IsSQLEngine(PostgreSQL))
		assert.True(t, IsSQLEngine(SQLite))
	})

	t.Run("test other databases", func(t *testing.T) {
		assert.False(t, IsSQLEngine(Empty))
	})
}

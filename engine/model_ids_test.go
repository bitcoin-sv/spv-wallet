package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIDs_GormDataType(t *testing.T) {
	t.Parallel()

	i := new(IDs)
	assert.Equal(t, gormTypeText, i.GormDataType())
}

func TestIDs_Scan(t *testing.T) {
	t.Parallel()

	t.Run("nil value", func(t *testing.T) {
		i := IDs{}
		err := i.Scan(nil)
		require.NoError(t, err)
		assert.Equal(t, 0, len(i))
	})

	t.Run("empty string", func(t *testing.T) {
		i := IDs{}
		err := i.Scan("\"\"")
		assert.Error(t, err)
		assert.Equal(t, 0, len(i))
	})

	t.Run("valid slice of ids", func(t *testing.T) {
		i := IDs{}
		err := i.Scan("[\"test1\",\"test2\"]")
		require.NoError(t, err)
		assert.Equal(t, 2, len(i))
		assert.Equal(t, "test1", i[0])
		assert.Equal(t, "test2", i[1])
	})

	t.Run("empty id slice", func(t *testing.T) {
		i := IDs{}
		err := i.Scan("[\"\"]")
		require.NoError(t, err)
		assert.Equal(t, 1, len(i))
		assert.Equal(t, "", i[0])
	})

	t.Run("invalid JSON", func(t *testing.T) {
		i := IDs{}
		err := i.Scan("[test1]")
		assert.Error(t, err)
	})
}

func TestIDs_Value(t *testing.T) {
	t.Parallel()

	t.Run("empty object", func(t *testing.T) {
		i := IDs{}
		value, err := i.Value()
		require.NoError(t, err)
		assert.Equal(t, "[]", value)
	})

	t.Run("ids present", func(t *testing.T) {
		i := IDs{"test1"}
		value, err := i.Value()
		require.NoError(t, err)
		assert.Equal(t, "[\"test1\"]", value)
	})
}

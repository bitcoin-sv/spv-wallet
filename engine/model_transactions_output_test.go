package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestXpubOutputValue_Scan(t *testing.T) {
	t.Parallel()

	t.Run("nil value", func(t *testing.T) {
		x := XpubOutputValue{}
		err := x.Scan(nil)
		require.NoError(t, err)
		assert.Equal(t, 0, len(x))
	})

	t.Run("empty string", func(t *testing.T) {
		x := XpubOutputValue{}
		err := x.Scan([]byte("\"\""))
		assert.NoError(t, err)
		assert.Equal(t, 0, len(x))
	})

	t.Run("empty string - incorrectly coded", func(t *testing.T) {
		x := XpubOutputValue{}
		err := x.Scan([]byte(""))
		assert.NoError(t, err)
		assert.Equal(t, 0, len(x))
	})

	t.Run("object", func(t *testing.T) {
		x := XpubOutputValue{}
		err := x.Scan([]byte("{\"xPubId\":543}"))
		require.NoError(t, err)
		assert.Equal(t, 1, len(x))
		assert.Equal(t, int64(543), x["xPubId"])
	})
}

func TestXpubOutputValue_Value(t *testing.T) {
	t.Parallel()

	t.Run("empty object", func(t *testing.T) {
		x := XpubOutputValue{}
		value, err := x.Value()
		require.NoError(t, err)
		assert.Equal(t, "{}", value)
	})

	t.Run("map present", func(t *testing.T) {
		x := XpubOutputValue{
			"xPubId": 123,
		}
		value, err := x.Value()
		require.NoError(t, err)
		assert.Equal(t, "{\"xPubId\":123}", value)
	})
}

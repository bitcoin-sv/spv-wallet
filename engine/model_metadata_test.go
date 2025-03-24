package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetadata_Scan(t *testing.T) {
	t.Parallel()

	t.Run("nil value", func(t *testing.T) {
		m := Metadata{}
		err := m.Scan(nil)
		require.NoError(t, err)
		assert.Equal(t, 0, len(m))
	})

	t.Run("empty string", func(t *testing.T) {
		m := Metadata{}
		err := m.Scan([]byte("\"\""))
		assert.NoError(t, err)
		assert.Equal(t, 0, len(m))
	})

	t.Run("empty string - incorrectly coded", func(t *testing.T) {
		m := Metadata{}
		err := m.Scan([]byte(""))
		assert.NoError(t, err)
		assert.Equal(t, 0, len(m))
	})

	t.Run("object", func(t *testing.T) {
		m := Metadata{}
		err := m.Scan([]byte("{\"test\":\"test2\"}"))
		require.NoError(t, err)
		assert.Equal(t, 1, len(m))
		assert.Equal(t, "test2", m["test"])
	})
}

func TestMetadata_Value(t *testing.T) {
	t.Parallel()

	t.Run("empty object", func(t *testing.T) {
		m := Metadata{}
		value, err := m.Value()
		require.NoError(t, err)
		assert.Equal(t, "{}", value)
	})

	t.Run("map present", func(t *testing.T) {
		m := Metadata{}
		m["test"] = "test2"
		value, err := m.Value()
		require.NoError(t, err)
		assert.Equal(t, "{\"test\":\"test2\"}", value)
	})
}

func TestXpubMetadata_Scan(t *testing.T) {
	t.Parallel()

	t.Run("nil value", func(t *testing.T) {
		x := XpubMetadata{}
		err := x.Scan(nil)
		require.NoError(t, err)
		assert.Equal(t, 0, len(x))
	})

	t.Run("empty string", func(t *testing.T) {
		x := XpubMetadata{}
		err := x.Scan([]byte("\"\""))
		assert.NoError(t, err)
		assert.Equal(t, 0, len(x))
	})

	t.Run("empty string - incorrectly coded", func(t *testing.T) {
		x := XpubMetadata{}
		err := x.Scan([]byte(""))
		assert.NoError(t, err)
		assert.Equal(t, 0, len(x))
	})

	t.Run("object", func(t *testing.T) {
		x := XpubMetadata{}
		err := x.Scan([]byte("{\"xPubId\":{\"test\":\"test2\"}}"))
		require.NoError(t, err)
		assert.Equal(t, 1, len(x))
		assert.Equal(t, 1, len(x["xPubId"]))
		assert.Equal(t, "test2", x["xPubId"]["test"])
	})
}

func TestXpubMetadata_Value(t *testing.T) {
	t.Parallel()

	t.Run("empty object", func(t *testing.T) {
		x := XpubMetadata{}
		value, err := x.Value()
		require.NoError(t, err)
		assert.Equal(t, "{}", value)
	})

	t.Run("map present", func(t *testing.T) {
		x := XpubMetadata{
			"xPubId": Metadata{
				"test": "test2",
			},
		}
		value, err := x.Value()
		require.NoError(t, err)
		assert.Equal(t, "{\"xPubId\":{\"test\":\"test2\"}}", value)
	})
}

package datastore

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestQueryParams_UnmarshalQueryParams will test the db Scanner of the QueryParams model
func TestQueryParams_UnmarshalQueryParams(t *testing.T) {
	t.Parallel()

	t.Run("nil value", func(t *testing.T) {
		m, err := UnmarshalQueryParams(nil)
		require.NoError(t, err)
		assert.Equal(t, QueryParams{}, m)
	})

	t.Run("empty string", func(t *testing.T) {
		m, err := UnmarshalQueryParams("\"\"")
		require.Error(t, err)
		assert.Equal(t, QueryParams{}, m)
	})

	t.Run("empty string - incorrectly coded", func(t *testing.T) {
		m, err := UnmarshalQueryParams("")
		require.Error(t, err)
		assert.Equal(t, QueryParams{}, m)
	})

	t.Run("object", func(t *testing.T) {
		var data map[string]interface{}
		err := json.Unmarshal([]byte(`{"page": 100}`), &data)
		require.NoError(t, err)
		var m QueryParams
		m, err = UnmarshalQueryParams(data)
		require.NoError(t, err)
		assert.Equal(t, QueryParams{Page: 100}, m)
	})
}

// TestMetadata_MarshalMetadata will test the db Valuer of the Metadata model
func TestMetadata_MarshalMetadata(t *testing.T) {
	t.Parallel()

	t.Run("empty object", func(t *testing.T) {
		q := QueryParams{}
		writer := MarshalQueryParams(q)
		require.NotNil(t, writer)
		b := bytes.NewBufferString("")
		writer.MarshalGQL(b)
		assert.Equal(t, "null", b.String())
	})

	t.Run("map present", func(t *testing.T) {
		q := QueryParams{Page: 11, PageSize: 35}
		writer := MarshalQueryParams(q)
		require.NotNil(t, writer)
		b := bytes.NewBufferString("")
		writer.MarshalGQL(b)
		assert.Equal(t, `{"page":11,"page_size":35}`+"\n", b.String())
	})
}

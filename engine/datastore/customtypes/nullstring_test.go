package customtypes

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"strconv"
	"testing"

	"github.com/99designs/gqlgen/graphql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testString = "test-string"

func TestNullString(t *testing.T) {
	t.Run("empty string", func(t *testing.T) {
		nt := new(NullString)
		assert.False(t, nt.Valid)
		assert.True(t, driver.IsValue(nt.String))
	})

	t.Run("string", func(t *testing.T) {
		nt := new(NullString)
		nt.String = testString
		nt.Valid = true
		assert.True(t, nt.Valid)
		assert.True(t, driver.IsValue(nt.String))
		assert.Equal(t, testString, nt.String)
	})
}

func TestNullString_IsZero(t *testing.T) {
	t.Run("nil string", func(t *testing.T) {
		nt := new(NullString)
		assert.True(t, nt.IsZero())
	})

	t.Run("string", func(t *testing.T) {
		nt := NullString{sql.NullString{Valid: true, String: testString}}
		assert.False(t, nt.IsZero())
	})
}

func TestMarshalNullString(t *testing.T) {
	t.Run("empty string", func(t *testing.T) {
		nt := new(NullString)
		marshaller := MarshalNullString(*nt)
		assert.Equal(t, graphql.Null, marshaller)
	})

	t.Run("string", func(t *testing.T) {
		nt := NullString{sql.NullString{Valid: true, String: testString}}
		marshaller := MarshalNullString(nt)
		var b bytes.Buffer
		marshaller.MarshalGQL(&b)
		assert.Equal(t, strconv.Quote(testString), b.String())
	})
}

func TestUnmarshalNullString(t *testing.T) {
	t.Run("nil string", func(t *testing.T) {
		nt, err := UnmarshalNullString(nil)
		require.NoError(t, err)
		assert.IsType(t, NullString{}, nt)
		assert.False(t, nt.Valid)
	})

	t.Run("empty string", func(t *testing.T) {
		val := NullString{}
		nt, err := UnmarshalNullString(val)
		require.Error(t, err)
		assert.IsType(t, NullString{}, nt)
		assert.False(t, nt.Valid)
	})

	t.Run("string", func(t *testing.T) {
		nt, err := UnmarshalNullString(testString)
		require.NoError(t, err)
		assert.IsType(t, NullString{}, nt)
		assert.True(t, nt.Valid)
		assert.Equal(t, testString, nt.String)
	})
}

func TestNullString_MarshalJSON(t *testing.T) {
	t.Run("nil string", func(t *testing.T) {
		nt := new(NullString)
		outBytes, err := nt.MarshalJSON()
		assert.Equal(t, []byte("null"), outBytes)
		require.NoError(t, err)
	})

	t.Run("empty string", func(t *testing.T) {
		nt := new(NullString)
		nt.Valid = true
		nt.String = ""
		outBytes, err := nt.MarshalJSON()
		assert.Equal(t, []byte("\"\""), outBytes)
		require.NoError(t, err)
	})

	t.Run("valid string", func(t *testing.T) {
		nt := NullString{sql.NullString{Valid: true, String: testString}}
		outBytes, err := nt.MarshalJSON()
		require.NoError(t, err)
		assert.NotNil(t, outBytes)
		assert.Equal(t, "\""+testString+"\"", string(outBytes))
	})
}

func TestNullString_UnmarshalJSON(t *testing.T) {
	t.Run("nil string", func(t *testing.T) {
		var nt NullString
		err := nt.UnmarshalJSON([]byte(nil))
		require.NoError(t, err)
		assert.False(t, nt.Valid)
	})

	t.Run("empty string", func(t *testing.T) {
		var nt NullString
		err := nt.UnmarshalJSON([]byte("\"\""))
		require.NoError(t, err)
		assert.True(t, nt.Valid)
		assert.Equal(t, "", nt.String)
	})

	t.Run("string", func(t *testing.T) {
		var nt NullString
		b := []byte("\"" + testString + "\"")
		err := nt.UnmarshalJSON(b)
		require.NoError(t, err)
		assert.True(t, nt.Valid)
		assert.Equal(t, testString, nt.String)
	})
}

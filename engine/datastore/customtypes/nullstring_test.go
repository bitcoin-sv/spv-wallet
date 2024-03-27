package customtypes

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/hex"
	"strconv"
	"testing"

	"github.com/99designs/gqlgen/graphql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

const testString = "test-string"

// TestNullString will test the basics of the null time struct
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

// TestMarshalNullString will test the method MarshalNullString()
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

// TestUnmarshalNullString will test the method UnmarshalNullString()
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

// TestNullString_IsZero will test the method IsZero()
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

// TestNullString_MarshalBSONValue will test the method MarshalBSONValue()
func TestNullString_MarshalBSONValue(t *testing.T) {
	t.Run("nil string", func(t *testing.T) {
		nt := new(NullString)
		outType, outBytes, err := nt.MarshalBSONValue()
		require.Equal(t, bsontype.Null, outType)
		assert.Nil(t, outBytes)
		require.NoError(t, err)
	})

	t.Run("empty string", func(t *testing.T) {
		nt := new(NullString)
		nt.Valid = true
		nt.String = ""
		outType, outBytes, err := nt.MarshalBSONValue()
		require.Equal(t, bsontype.String, outType)
		assert.Equal(t, "0100000000", hex.EncodeToString(outBytes))
		require.NoError(t, err)
	})

	t.Run("valid string", func(t *testing.T) {
		nt := NullString{sql.NullString{Valid: true, String: testString}}
		outType, outBytes, err := nt.MarshalBSONValue()
		require.NoError(t, err)
		assert.Equal(t, bsontype.String, outType)
		assert.NotNil(t, outBytes)
		outHex := hex.EncodeToString(outBytes[:])
		_, inHex, _ := bson.MarshalValue(testString)
		assert.Equal(t, hex.EncodeToString(inHex[:]), outHex)
	})
}

// TestNullString_UnmarshalBSONValue will test the method UnmarshalBSONValue()
func TestNullString_UnmarshalBSONValue(t *testing.T) {
	t.Run("nil string", func(t *testing.T) {
		var nt NullString
		err := nt.UnmarshalBSONValue(bsontype.Null, nil)
		require.NoError(t, err)
		assert.False(t, nt.Valid)
	})

	t.Run("empty string", func(t *testing.T) {
		var nt NullString
		b, _ := hex.DecodeString("0100000000")
		err := nt.UnmarshalBSONValue(bsontype.String, b)
		require.NoError(t, err)
		assert.True(t, nt.Valid)
		assert.Equal(t, "", nt.String)
	})

	t.Run("string", func(t *testing.T) {
		var nt NullString
		b, _ := hex.DecodeString("0c000000746573742d737472696e6700")
		err := nt.UnmarshalBSONValue(bsontype.String, b)
		require.NoError(t, err)
		assert.True(t, nt.Valid)
		assert.Equal(t, testString, nt.String)
	})
}

// TestNullString_MarshalJSON will test the method MarshalJSON()
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

// TestNullString_UnmarshalJSON will test the method UnmarshalJSON()
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

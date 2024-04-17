package engine

import (
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

// TestMetaDataScan will test the db Scanner of the Metadata model
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

// TestMetaDataValue will test the db Valuer of the Metadata model
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

// TestXpubMetadata_Scan will test the db Scanner of the XpubMetadata model
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

// TestXpubMetadata_Value will test the db Valuer of the XpubMetadata model
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

// TestMetadata_MarshalBSONValue will test the method MarshalBSONValue()
func TestMetadata_MarshalBSONValue(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		m := Metadata{}
		outType, outBytes, err := m.MarshalBSONValue()
		require.Equal(t, bsontype.Null, outType)
		assert.Nil(t, outBytes)
		assert.NoError(t, err)
	})

	t.Run("valid", func(t *testing.T) {
		m := Metadata{
			"test-key": "test-value",
		}
		outType, outBytes, err := m.MarshalBSONValue()
		require.NoError(t, err)
		assert.Equal(t, bsontype.Array, outType)
		assert.NotNil(t, outBytes)
		outHex := hex.EncodeToString(outBytes[:])

		out := new(map[string]interface{})
		outBytes, hexErr := hex.DecodeString(outHex)
		require.NoError(t, hexErr)
		err = bson.Unmarshal(outBytes, out)
		require.NoError(t, err)
		jsonOut, jsonErr := json.Marshal(out)
		require.NoError(t, jsonErr)
		assert.Equal(t, string(jsonOut), "{\"0\":{\"k\":\"test-key\",\"v\":\"test-value\"}}")

		// check that it is not normal marshaling
		_, inHex, _ := bson.MarshalValue(m)
		assert.NotEqual(t, hex.EncodeToString(inHex[:]), outHex)
	})
}

// TestMetadata_UnmarshalBSONValue will test the method UnmarshalBSONValue()
func TestMetadata_UnmarshalBSONValue(t *testing.T) {
	t.Parallel()

	t.Run("nil string", func(t *testing.T) {
		var m Metadata
		err := m.UnmarshalBSONValue(bsontype.Null, nil)
		require.NoError(t, err)
		assert.Nil(t, m)
	})

	t.Run("string", func(t *testing.T) {
		var m Metadata
		// this hex is a bson array [{k: "test-key", v: "test-value"}]
		b, _ := hex.DecodeString("2f000000033000270000000276000b000000746573742d76616c756500026b0009000000746573742d6b6579000000")
		err := m.UnmarshalBSONValue(bsontype.Array, b)
		require.NoError(t, err)
		assert.Equal(t, Metadata{"test-key": "test-value"}, m)
	})
}

// TestXpubMetadata_MarshalBSONValue will test the method MarshalBSONValue()
func TestXpubMetadata_MarshalBSONValue(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		m := XpubMetadata{}
		outType, outBytes, err := m.MarshalBSONValue()
		require.Equal(t, bsontype.Null, outType)
		assert.Nil(t, outBytes)
		assert.NoError(t, err)
	})

	t.Run("valid", func(t *testing.T) {
		m := XpubMetadata{
			"xPubId": Metadata{
				"test-key": "test-value",
			},
		}
		outType, outBytes, err := m.MarshalBSONValue()
		require.NoError(t, err)
		assert.Equal(t, bsontype.Array, outType)
		assert.NotNil(t, outBytes)
		outHex := hex.EncodeToString(outBytes[:])

		out := new(map[string]interface{})
		outBytes, hexErr := hex.DecodeString(outHex)
		require.NoError(t, hexErr)
		err = bson.Unmarshal(outBytes, out)
		require.NoError(t, err)
		jsonOut, jsonErr := json.Marshal(out)
		require.NoError(t, jsonErr)
		assert.Equal(t, string(jsonOut), "{\"0\":{\"k\":\"test-key\",\"v\":\"test-value\",\"x\":\"xPubId\"}}")

		// check that it is not normal marshaling
		_, inHex, _ := bson.MarshalValue(m)
		assert.NotEqual(t, hex.EncodeToString(inHex[:]), outHex)
	})
}

// TestXpubMetadata_UnmarshalBSONValue will test the method UnmarshalBSONValue()
func TestXpubMetadata_UnmarshalBSONValue(t *testing.T) {
	t.Parallel()

	t.Run("nil string", func(t *testing.T) {
		var m XpubMetadata
		err := m.UnmarshalBSONValue(bsontype.Null, nil)
		require.NoError(t, err)
		assert.Nil(t, m)
	})

	t.Run("string", func(t *testing.T) {
		var m XpubMetadata
		// this hex is a bson array {xPubId: [{k: "test-key", v: "test-value", x: "xpubId"}]}
		b, _ := hex.DecodeString("3d000000033000350000000278000700000078507562496400026b0009000000746573742d6b6579000276000b000000746573742d76616c7565000000")
		err := m.UnmarshalBSONValue(bsontype.Array, b)
		require.NoError(t, err)
		assert.Equal(t, XpubMetadata{"xPubId": Metadata{"test-key": "test-value"}}, m)
	})
}

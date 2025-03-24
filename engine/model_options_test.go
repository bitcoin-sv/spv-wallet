package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testEncryption = "35dbe09a941a90a5f59e57020face68860d7b284b7b2973a58de8b4242ec5a925a40ac2933b7e45e78a0b3a13123520e46f9566815589ba2d345577dadee0d5e"
)

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("Get opts", func(t *testing.T) {
		opt := New()
		assert.IsType(t, *new(ModelOps), opt)
	})

	t.Run("apply opts", func(t *testing.T) {
		opt := New()
		m := new(Model)
		m.SetOptions(opt)
		assert.Equal(t, true, m.IsNew())
	})
}

func TestWithMetadata(t *testing.T) {
	t.Parallel()

	t.Run("Get opts", func(t *testing.T) {
		opt := WithMetadata("key", "value")
		assert.IsType(t, *new(ModelOps), opt)
	})

	t.Run("apply opts", func(t *testing.T) {
		opt := WithMetadata("key", "value")
		m := new(Model)
		m.SetOptions(opt)
		assert.Equal(t, "value", m.Metadata["key"])
	})
}

func TestWithMetadataFromJSON(t *testing.T) {
	t.Parallel()

	t.Run("Get opts", func(t *testing.T) {
		opt := WithMetadataFromJSON([]byte(`{"key": "value"}`))
		assert.IsType(t, *new(ModelOps), opt)
	})

	t.Run("apply opts", func(t *testing.T) {
		opt := WithMetadataFromJSON([]byte(`{"key": "value"}`))
		m := new(Model)
		m.SetOptions(opt)
		assert.Equal(t, "value", m.Metadata["key"])
	})
}

func TestWithXPub(t *testing.T) {
	t.Parallel()

	t.Run("Get opts", func(t *testing.T) {
		opt := WithXPub(testXPub)
		assert.IsType(t, *new(ModelOps), opt)
	})

	t.Run("apply opts", func(t *testing.T) {
		opt := WithXPub(testXPub)
		m := new(Model)
		m.SetOptions(opt)
		assert.Equal(t, testXPub, m.rawXpubKey)
	})
}

func TestWithEncryptionKey(t *testing.T) {
	t.Parallel()

	t.Run("Get opts", func(t *testing.T) {
		opt := WithEncryptionKey(testEncryption)
		assert.IsType(t, *new(ModelOps), opt)
	})

	t.Run("apply opts", func(t *testing.T) {
		opt := WithEncryptionKey(testEncryption)
		m := new(Model)
		m.SetOptions(opt)
		assert.Equal(t, testEncryption, m.encryptionKey)
	})
}

func TestWithMetadatas(t *testing.T) {
	t.Parallel()

	t.Run("Get opts", func(t *testing.T) {
		opt := WithMetadatas(map[string]interface{}{
			"key": "value",
		})
		assert.IsType(t, *new(ModelOps), opt)
	})

	t.Run("apply opts", func(t *testing.T) {
		opt := WithMetadatas(map[string]interface{}{
			"key": "value",
		})
		m := new(Model)
		m.SetOptions(opt)
		assert.Equal(t, "value", m.Metadata["key"])
	})
}

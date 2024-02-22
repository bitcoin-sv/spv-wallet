package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNew will test the method New()
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

// TestWithMetadata will test the method WithMetadata()
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

// TestWithMetadataFromJSON will test the method WithMetadataFromJSON()
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

// TestWithXPub will test the method WithXPub()
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

// TestWithEncryptionKey will test the method WithEncryptionKey()
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

// TestWithMetadatas will test the method WithMetadatas()
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

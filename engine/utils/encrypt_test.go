package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testEncryptValue = "##!(TEST)!##"
)

func TestEncrypt(t *testing.T) {
	t.Run("empty key", func(t *testing.T) {
		encrypted, err := Encrypt("", "")
		require.Error(t, err)
		assert.Equal(t, "", encrypted)
	})

	t.Run("invalid key", func(t *testing.T) {
		encrypted, err := Encrypt("123", "")
		require.Error(t, err)
		assert.Equal(t, "", encrypted)
	})

	t.Run("valid small key, no value", func(t *testing.T) {
		encryptionKey, err := RandomHex(64)
		require.NoError(t, err)

		var encrypted string
		encrypted, err = Encrypt(encryptionKey, "")
		require.NoError(t, err)
		assert.NotEqual(t, 0, len(encrypted))

		var decrypted string
		decrypted, err = Decrypt(encryptionKey, encrypted)
		require.NoError(t, err)
		assert.Equal(t, "", decrypted)
	})

	t.Run("hardcoded small key with value", func(t *testing.T) {
		encryptionKey := "a7f024b811012a88"

		encrypted, err := Encrypt(encryptionKey, testEncryptValue)
		require.NoError(t, err)

		var decrypted string
		decrypted, err = Decrypt(encryptionKey, encrypted)
		require.NoError(t, err)
		assert.Equal(t, testEncryptValue, decrypted)
	})

	t.Run("hardcoded 32 key with value", func(t *testing.T) {
		encryptionKey := "be5d67424e5e3d7bb0ca69da68e423774062aebf76cb265490ac2d57d2fa2933"

		encrypted, err := Encrypt(encryptionKey, testEncryptValue)
		require.NoError(t, err)

		var decrypted string
		decrypted, err = Decrypt(encryptionKey, encrypted)
		require.NoError(t, err)
		assert.Equal(t, testEncryptValue, decrypted)
	})

	t.Run("hardcoded 64 key with value", func(t *testing.T) {
		encryptionKey := "35dbe09a941a90a5f59e57020face68860d7b284b7b2973a58de8b4242ec5a925a40ac2933b7e45e78a0b3a13123520e46f9566815589ba2d345577dadee0d5e"

		encrypted, err := Encrypt(encryptionKey, testEncryptValue)
		require.NoError(t, err)

		var decrypted string
		decrypted, err = Decrypt(encryptionKey, encrypted)
		require.NoError(t, err)
		assert.Equal(t, testEncryptValue, decrypted)
	})
}

package utils

import (
	"github.com/bitcoinschema/go-bitcoin/v2"
)

// Encrypt will encrypt the value using the encryption key
func Encrypt(encryptionKey, encryptValue string) (string, error) {

	// Get the keys seeded with the encryption key
	privateKey, _, err := bitcoin.PrivateAndPublicKeys(encryptionKey)
	if err != nil {
		return "", err
	}

	// Encrypt the private key
	var encryptedValue string
	if encryptedValue, err = bitcoin.EncryptWithPrivateKey(
		privateKey, encryptValue,
	); err != nil {
		return "", err
	}

	return encryptedValue, nil
}

// Decrypt will take the data and decrypt using a char(64) key
func Decrypt(encryptionKey, data string) (string, error) {
	return bitcoin.DecryptWithPrivateKeyString(encryptionKey, data)
}

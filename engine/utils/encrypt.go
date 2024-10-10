package utils

import (
	ecies "github.com/bitcoin-sv/go-sdk/compat/ecies"
	ec "github.com/bitcoin-sv/go-sdk/primitives/ec"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
)

// Encrypt will encrypt the value using the encryption key
func Encrypt(encryptionKey, encryptValue string) (string, error) {
	// Get the keys seeded with the encryption key
	privateKey, err := ec.PrivateKeyFromHex(encryptionKey)
	if err != nil {
		return "", spverrors.Wrapf(err, "error getting private keys from encryption key")
	}

	// Encrypt the private key
	var encryptedValue string
	if encryptedValue, err = ecies.EncryptSingle(
		encryptValue, privateKey,
	); err != nil {
		return "", spverrors.Wrapf(err, "error encrypting data with private key")
	}

	return encryptedValue, nil
}

// Decrypt will take the data and decrypt using a char(64) key
func Decrypt(encryptionKey, data string) (string, error) {
	// keyString, err := bitcoin.DecryptWithPrivateKeyString(encryptionKey, data)
	privateKey, err := ec.PrivateKeyFromHex(encryptionKey)
	keyString, err := ecies.DecryptSingle(data, privateKey)
	return keyString, spverrors.Wrapf(err, "error decrypting data with private key")
}

package utils

import (
	bip32 "github.com/bitcoin-sv/go-sdk/compat/bip32"
	compat "github.com/bitcoin-sv/go-sdk/compat/bip32"
	ec "github.com/bitcoin-sv/go-sdk/primitives/ec"
	"github.com/bitcoin-sv/go-sdk/script"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
)

// DeriveChildKeyFromHex derive the child extended key from the hex string
func DeriveChildKeyFromHex(hdKey *bip32.ExtendedKey, hexHash string) (*bip32.ExtendedKey, error) {
	var childKey *bip32.ExtendedKey
	childKey = hdKey

	childNums, err := GetChildNumsFromHex(hexHash)
	if err != nil {
		return nil, err
	}

	for _, num := range childNums {
		if childKey, err = childKey.Child(num); err != nil {
			return nil, spverrors.Wrapf(err, "failed to derive child key")
		}
	}

	return childKey, nil
}

// DerivePublicKey will derive the internal and external address from a key
func DerivePublicKey(hdKey *bip32.ExtendedKey, chain uint32, num uint32) (*ec.PublicKey, error) {
	if hdKey == nil {
		return nil, ErrHDKeyNil
	}

	pubKeys, err := compat.GetPublicKeysForPath(hdKey, num)
	if err != nil {
		return nil, spverrors.Wrapf(err, "failed to derive public key")
	}

	return pubKeys[chain], nil
}

// ValidateXPub will check the xPub key for length & validation
func ValidateXPub(rawKey string) (*bip32.ExtendedKey, error) {

	// Validate the xpub (length)
	if len(rawKey) != XpubKeyLength {
		return nil, spverrors.ErrXpubInvalidLength
	}

	// Parse the xPub into an HD key
	hdKey, err := compat.GetHDKeyFromExtendedPublicKey(rawKey)
	if err != nil {
		return nil, spverrors.Wrapf(err, "failed to parse xpub")
	} else if hdKey.String() != rawKey { // Sanity check (might not be needed)
		return nil, spverrors.ErrXpubNoMatch
	}
	return hdKey, nil
}

// DeriveAddress will derive the given address from a key
func DeriveAddress(hdKey *bip32.ExtendedKey, chain uint32, num uint32) (address string, err error) {

	// Don't panic
	if hdKey == nil {
		return "", ErrHDKeyNil
	}

	var child *bip32.ExtendedKey
	if child, err = compat.GetHDKeyByPath(hdKey, chain, num); err != nil {
		return "", spverrors.Wrapf(err, "failed to derive child key")
	}

	var pubKey *ec.PublicKey
	if pubKey, err = child.ECPubKey(); err != nil {
		// Should never error since the previous method ensures a valid hdKey
		return "", spverrors.Wrapf(err, "failed to derive public key")
	}

	var addressScript *script.Address
	if addressScript, err = script.NewAddressFromPublicKey(pubKey, true); err != nil {
		// Should never error if the pubKeys are valid keys
		return "", spverrors.Wrapf(err, "failed to derive address")
	}

	return addressScript.AddressString, nil
}

// DeriveAddresses will derive the internal and external address from a key
func DeriveAddresses(hdKey *bip32.ExtendedKey, num uint32) (external, internal string, err error) {

	// Don't panic
	if hdKey == nil {
		return "", "", ErrHDKeyNil
	}

	// Derive the address
	var addresses []string
	if addresses, err = compat.GetAddressesForPath(
		hdKey, num,
	); err != nil {
		return
	} else if len(addresses) != 2 { // Sanity check might not be needed
		return "", "", ErrDeriveFailed
	}
	external = addresses[0]
	internal = addresses[1]
	return
}

// DerivePrivateKeyFromHex will derive the private key from the extended key using the hex as the derivation paths
func DerivePrivateKeyFromHex(hdKey *bip32.ExtendedKey, hexString string) (*ec.PrivateKey, error) {
	if hdKey == nil {
		return nil, ErrHDKeyNil
	}

	childKey, err := DeriveChildKeyFromHex(hdKey, hexString)
	if err != nil {
		return nil, err
	}

	var privKey *ec.PrivateKey
	if privKey, err = childKey.ECPrivKey(); err != nil {
		return nil, spverrors.Wrapf(err, "failed to derive private key")
	}

	return privKey, nil
}

// DerivePublicKeyFromHex will derive the public key from the extended key using the hex as the derivation paths
func DerivePublicKeyFromHex(hdKey *bip32.ExtendedKey, hexString string) (*ec.PublicKey, error) {
	if hdKey == nil {
		return nil, ErrHDKeyNil
	}

	childKey, err := DeriveChildKeyFromHex(hdKey, hexString)
	if err != nil {
		return nil, err
	}

	var pubKey *ec.PublicKey
	if pubKey, err = childKey.ECPubKey(); err != nil {
		return nil, spverrors.Wrapf(err, "failed to derive public key")
	}

	return pubKey, nil
}

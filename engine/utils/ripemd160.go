package utils

import (
	"crypto/sha256"

	"golang.org/x/crypto/ripemd160"
)

// NOTE: Temporary implementation - needs to be tracked in GO-SDK SPV-1035

// Sha256 hashes with SHA256
func Sha256(b []byte) []byte {
	data := sha256.Sum256(b)
	return data[:]
}

// Ripemd160 hashes with RIPEMD160
func Ripemd160(b []byte) ([]byte, error) {
	ripe := ripemd160.New()
	_, err := ripe.Write(b[:])
	return ripe.Sum(nil), err
}

// Hash160 hashes with SHA256 and then hashes again with RIPEMD160.
func Hash160(b []byte) ([]byte, error) {
	hash := Sha256(b)
	return Ripemd160(hash[:])
}

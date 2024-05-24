package type42

import (
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"fmt"
	"math/big"

	"github.com/libsv/go-bk/bec"
)

// calculateHMAC calculates the HMAC of the provided public shared secret using a reference string.
// The reference string, which acts as the HMAC key, can be an invoice number or any other string.
// The function returns the HMAC result as a byte slice.
//
// Parameters:
// - pubSharedSecret: A byte slice representing the public shared secret to be hashed.
// - reference: A string used as the HMAC key.
//
// Returns:
// - A byte slice containing the HMAC result as a byte slice or an error if the HMAC calculation fails.
func calculateHMAC(pubSharedSecret []byte, reference string) ([]byte, error) {
	if reference == "" {
		return nil, errors.New("invalid invoice number")
	}
	h := hmac.New(sha256.New, pubSharedSecret)
	if _, err := h.Write([]byte(reference)); err != nil {
		return nil, fmt.Errorf("error writing HMAC message - %w", err)
	}
	return h.Sum(nil), nil
}

// calculateDedicatedPublicKey calculates the dedicated public key (dPK) using the HMAC result and the receiver's public key.
// The HMAC result is used as a scalar to perform elliptic curve scalar multiplication and point addition.
// Returns the resulting public key or an error if the calculation fails.
func calculateDedicatedPublicKey(hmacResult []byte, receiverPubKey *bec.PublicKey) (*bec.PublicKey, error) {
	if len(hmacResult) == 0 {
		return nil, fmt.Errorf("HMAC result is empty")
	}
	if receiverPubKey == nil {
		return nil, fmt.Errorf("receiver public key is nil")
	}

	// Convert HMAC result to a big integer
	hn := new(big.Int).SetBytes(hmacResult)
	curve := bec.S256()          // Use secp256k1 curve
	hn.Mod(hn, curve.Params().N) // Ensure the scalar is within the curve order

	// Perform scalar multiplication: hn * G
	rx, ry := curve.ScalarBaseMult(hn.Bytes())

	// Perform point addition: (receiverPubKey + (hn * G))
	dedicatedPubKeyX, dedicatedPubKeyY := curve.Add(receiverPubKey.X, receiverPubKey.Y, rx, ry)

	// Verify that the resulting point is on the curve
	if !curve.IsOnCurve(dedicatedPubKeyX, dedicatedPubKeyY) {
		return nil, fmt.Errorf("resulting public key is not on curve")
	}

	// Create the dedicated public key
	dPK := &bec.PublicKey{
		Curve: curve,
		X:     dedicatedPubKeyX,
		Y:     dedicatedPubKeyY,
	}
	return dPK, nil
}

// DeriveLinkedKey derives a child public key from the source public key and link it with public key
// with use of invoiceNumber as reference of this derivation.
func DeriveLinkedKey(source bec.PublicKey, linkPubKey bec.PublicKey, invoiceNumber string) (*bec.PublicKey, error) {
	// Check for nil receiver public key
	if linkPubKey.X == nil || linkPubKey.Y == nil {
		return nil, fmt.Errorf("receiver public key is nil")
	}

	// Compute the shared secret
	sharedSecret := source.SerialiseCompressed()

	// Compute the HMAC result
	hmacResult, err := calculateHMAC(sharedSecret, invoiceNumber)
	if err != nil {
		return nil, err
	}

	// Calculate the dedicated public key
	dPK, err := calculateDedicatedPublicKey(hmacResult, &linkPubKey)
	if err != nil {
		return nil, err
	}

	return dPK, nil
}

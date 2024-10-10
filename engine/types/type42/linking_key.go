package type42

import (
	"crypto/hmac"
	"crypto/sha256"
	"math/big"

	ec "github.com/bitcoin-sv/go-sdk/primitives/ec"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
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
func calculateHMAC(pubSharedSecret []byte, message string) ([]byte, error) {
	if message == "" {
		return nil, spverrors.Newf("invalid invoice number")
	}
	h := hmac.New(sha256.New, pubSharedSecret)
	if _, err := h.Write([]byte(message)); err != nil {
		return nil, spverrors.Wrapf(err, "error writing HMAC message -")
	}
	return h.Sum(nil), nil
}

// calculateLinkedPublicKey calculates the dedicated public key (dPK) using the HMAC result and the receiver's public key.
// The HMAC result is used as a scalar to perform elliptic curve scalar multiplication and point addition.
// Returns the resulting public key or an error if the calculation fails.
func calculateLinkedPublicKey(hmacResult []byte, receiverPubKey *ec.PublicKey) (*ec.PublicKey, error) {
	if len(hmacResult) == 0 {
		return nil, spverrors.Newf("HMAC result is empty")
	}
	if receiverPubKey == nil {
		return nil, spverrors.Newf("receiver public key is nil")
	}

	// Convert HMAC result to a big integer
	hn := new(big.Int).SetBytes(hmacResult)
	curve := ec.S256()           // Use secp256k1 curve
	hn.Mod(hn, curve.Params().N) // Ensure the scalar is within the curve order

	// Perform scalar multiplication: hn * G
	rx, ry := curve.ScalarBaseMult(hn.Bytes())

	// Perform point addition: (receiverPubKey + (hn * G))
	dedicatedPubKeyX, dedicatedPubKeyY := curve.Add(receiverPubKey.X, receiverPubKey.Y, rx, ry)

	// Verify that the resulting point is on the curve
	if !curve.IsOnCurve(dedicatedPubKeyX, dedicatedPubKeyY) {
		return nil, spverrors.Newf("resulting public key is not on curve")
	}

	// Create the dedicated public key
	dPK := &ec.PublicKey{
		Curve: curve,
		X:     dedicatedPubKeyX,
		Y:     dedicatedPubKeyY,
	}
	return dPK, nil
}

// DeriveLinkedKey derives a child public key from the source public key and link it with public key
// with use of invoiceNumber as reference of this derivation.
func DeriveLinkedKey(source *ec.PublicKey, linkPubKey *ec.PublicKey, invoiceNumber string) (*ec.PublicKey, error) {
	if source == nil || linkPubKey == nil {
		return nil, spverrors.Newf("source or receiver public key is nil")
	}

	// Check for nil receiver public key
	if source.X == nil || source.Y == nil {
		return nil, spverrors.Newf("source public key is nil")
	}
	if linkPubKey.X == nil || linkPubKey.Y == nil {
		return nil, spverrors.Newf("receiver public key is nil")
	}

	// Compute the shared secret
	publicKeyBytes := source.SerializeCompressed()

	// Compute the HMAC result
	hmacResult, err := calculateHMAC(publicKeyBytes, invoiceNumber)
	if err != nil {
		return nil, err
	}

	// Calculate the dedicated public key
	linkedPK, err := calculateLinkedPublicKey(hmacResult, linkPubKey)
	if err != nil {
		return nil, err
	}

	return linkedPK, nil
}

// Package pike2 provides functionality to work with Pay-to-PubKey-Hash (P2PKH) scripts in a blockchain context.
//
// P2PKH is a common type of Bitcoin transaction that locks funds to a specific public key hash, requiring the
// corresponding private key to produce a valid signature for spending the funds. This package offers utilities
// to create, decode, and validate P2PKH scripts.
//
// The package includes:
// - Functions to generate P2PKH addresses from public keys.
// - Methods to construct P2PKH scriptPubKey and scriptSig scripts.
// - Utilities to decode and inspect P2PKH scripts and addresses.
// - Validation functions to ensure the integrity and correctness of P2PKH scripts.
//
// This package is intended for developers working with Bitcoin or other cryptocurrencies that support P2PKH transactions.
// It abstracts the low-level details and provides a high-level interface for creating and handling P2PKH scripts.
package pike2

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"

	"github.com/libsv/go-bk/bec"
	"github.com/libsv/go-bk/crypto"
	"github.com/libsv/go-bt/v2/bscript"
	"github.com/libsv/go-bt/v2/bscript/interpreter"
)

// PikeOutputsTemplate represents the script and satoshis for a Pike output
type PikeOutputsTemplate struct {
	Script   string `json:"script"`
	Satoshis uint64 `json:"satoshis"`
}

// CreatePikeOutputsTemplate creates a single output with the PIKE template
func CreatePikeOutputsTemplate(satoshis uint64) ([]PikeOutputsTemplate, error) {

	opcodes := []byte{
		bscript.OpDUP,
		bscript.OpHASH160,
		bscript.OpPUBKEYHASH,
		bscript.OpEQUALVERIFY,
		bscript.OpCHECKSIG,
	}

	// Convert opcodes to hexadecimal string
	scriptHex := hex.EncodeToString(opcodes)

	// Check if the conversion was successful
	if scriptHex == "" {
		return nil, errors.New("failed to create script hex")
	}

	// Create and return the PikeOutputsTemplate
	return []PikeOutputsTemplate{
		{
			Script:   scriptHex,
			Satoshis: satoshis,
		},
	}, nil
}

// CalculatePubSharedSecret creates public shared secret base on sender public key
func CalculatePubSharedSecret(senderPubKey *bec.PublicKey) *bec.PublicKey {
	privateKeyOne := newPrivateKey(1)
	pubSharedSecret := new(bec.PublicKey)
	pubSharedSecret.X, pubSharedSecret.Y = bec.S256().ScalarMult(senderPubKey.X, senderPubKey.Y, privateKeyOne.D.Bytes())
	return pubSharedSecret
}

// newPrivateKey create a new private key with a given integer value
func newPrivateKey(value int64) *bec.PrivateKey {
	const privKeyBytesLen = 32
	privKeyBytes := make([]byte, privKeyBytesLen)
	bigValue := big.NewInt(value)
	copy(privKeyBytes[privKeyBytesLen-len(bigValue.Bytes()):], bigValue.Bytes())
	privateKey, _ := bec.PrivKeyFromBytes(bec.S256(), privKeyBytes)
	return privateKey
}

// calculateHMAC calculates the HMAC of the public shared secret using a reference string and an index.
// The reference string and index are concatenated to form the HMAC key.
// Returns the HMAC result as a byte slice or an error if the HMAC calculation fails.
func calculateHMAC(pubSharedSecret []byte, reference string, idx int) ([]byte, error) {
	h := hmac.New(sha256.New, []byte(fmt.Sprintf("%s-%d", reference, idx)))
	if _, err := h.Write(pubSharedSecret); err != nil {
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

// convertToCompressedPublicKey converts the dedicated public key to its compressed format
// and then hashes it using SHA-256 followed by RIPEMD-160 (Hash160).
// Returns the compressed and hashed public key as a byte slice.
func convertToCompressedPublicKey(dPK *bec.PublicKey) []byte {
	// Serialize the public key to compressed format
	dPKBytes := dPK.SerialiseCompressed()

	// Apply Hash160 (SHA-256 followed by RIPEMD-160) to the compressed public key
	dPKBytesCompressed := crypto.Hash160(dPKBytes)

	// Return the compressed and hashed public key
	return dPKBytesCompressed
}

// GenerateLockingScriptsFromTemplates converts Pike outputs templates to scripts
func GenerateLockingScriptsFromTemplates(outputsTemplate []PikeOutputsTemplate, senderPubKey, receiverPubKey *bec.PublicKey, reference string) ([]string, error) {
	lockingScripts := make([]string, len(outputsTemplate))
	parser := interpreter.DefaultOpcodeParser{}
	pubSharedSecret := CalculatePubSharedSecret(senderPubKey)

	for idx, output := range outputsTemplate {
		templateScript, err := bscript.NewFromHexString(output.Script)
		if err != nil {
			return nil, fmt.Errorf("error creating script from hex string - %w", err)
		}
		parsedScript, err := parser.Parse(templateScript)
		if err != nil {
			return nil, fmt.Errorf("error parsing template script - %w", err)
		}

		hmacResult, err := calculateHMAC(pubSharedSecret.SerialiseCompressed(), reference, idx)
		if err != nil {
			return nil, err
		}

		dPK, err := calculateDedicatedPublicKey(hmacResult, receiverPubKey)
		if err != nil {
			return nil, err
		}

		dPKBytesCompressed := convertToCompressedPublicKey(dPK)

		// Create a new script with the dedicated public key
		newScript := new(bscript.Script)
		if err := newScript.AppendPushData(dPKBytesCompressed); err != nil {
			return nil, fmt.Errorf("error appending public key hash to script: %w", err)
		}

		// Parse the public key hash script
		pkhParsed, err := parser.Parse(newScript)
		if err != nil {
			return nil, fmt.Errorf("error parsing public key hash script: %w", err)
		}

		// Replace OP_PUBKEYHASH with the actual public key hash
		evaluated := make([]interpreter.ParsedOpcode, 0, len(parsedScript))
		for _, op := range parsedScript {
			switch op.Value() {
			case bscript.OpPUBKEYHASH:
				evaluated = append(evaluated, pkhParsed...)
			case bscript.OpPUBKEY:
				evaluated = append(evaluated, pkhParsed...) // Currently just appending the same key
			default:
				evaluated = append(evaluated, op)
			}
		}

		// Unparse the evaluated opcodes back into a script
		finalScript, err := parser.Unparse(evaluated)
		if err != nil {
			return nil, fmt.Errorf("error unparsing evaluated opcodes: %w", err)
		}

		lockingScripts[idx] = finalScript.String()
	}

	return lockingScripts, nil
}

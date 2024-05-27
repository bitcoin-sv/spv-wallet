package type42

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"

	"github.com/libsv/go-bk/bec"
)

// Helper function to generate a random public key
func generateRandomPublicKey() *bec.PublicKey {
	privKey, _ := bec.NewPrivateKey(bec.S256())
	return privKey.PubKey()
}

// Helper function to compute the SHA-256 hash of a string and return it as a slice of bytes
func sha256Hash(data string) []byte {
	hash := sha256.Sum256([]byte(data))
	return hash[:]
}

// Helper function to compute the expected HMAC
func computeExpectedHMAC(pubSharedSecret []byte, reference string, idx int) []byte {
	h := hmac.New(sha256.New, pubSharedSecret)
	h.Write([]byte(fmt.Sprintf("%s-%d", reference, idx)))
	return h.Sum(nil)
}

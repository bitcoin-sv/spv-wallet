package type42

import (
	"encoding/hex"
	"testing"

	"github.com/libsv/go-bk/bec"
	assert "github.com/stretchr/testify/require"
)

func TestDeriveLinkedKey(t *testing.T) {
	sourcePubKeyHex := "027c1404c3ecb034053e6dd90bc68f7933284559c7d0763367584195a8796d9b0e"
	sourcePubKeyBytes, err := hex.DecodeString(sourcePubKeyHex)
	assert.NoError(t, err)
	sourcePubKey, err := bec.ParsePubKey(sourcePubKeyBytes, bec.S256())
	assert.NoError(t, err)

	linkPubKeyHex := "03a34e456deecb6e6e9237e63e5b7d045d1d2a456eb6be43de1ec4e9ac9a07b50d"
	linkPubKeyBytes, err := hex.DecodeString(linkPubKeyHex)
	assert.NoError(t, err)
	linkPubKey, err := bec.ParsePubKey(linkPubKeyBytes, bec.S256())
	assert.NoError(t, err)

	validHMAC, err := calculateHMAC(sourcePubKey.SerialiseCompressed(), "valid-invoice")
	assert.NoError(t, err)
	validDerivedKey, err := calculateLinkedPublicKey(validHMAC, linkPubKey)
	assert.NoError(t, err)

	validTests := []struct {
		name           string
		source         bec.PublicKey
		linkPubKey     bec.PublicKey
		invoiceNumber  string
		expectedResult *bec.PublicKey
	}{
		{
			name:           "valid case",
			source:         *sourcePubKey,
			linkPubKey:     *linkPubKey,
			invoiceNumber:  "valid-invoice",
			expectedResult: validDerivedKey,
		},
	}

	t.Run("Valid Cases", func(t *testing.T) {
		for _, tt := range validTests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				result, err := DeriveLinkedKey(&tt.source, &tt.linkPubKey, tt.invoiceNumber)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			})
		}
	})

	errorTests := []struct {
		name          string
		source        bec.PublicKey
		linkPubKey    bec.PublicKey
		invoiceNumber string
	}{
		{
			name:          "invalid HMAC result",
			source:        *sourcePubKey,
			linkPubKey:    *linkPubKey,
			invoiceNumber: "", // Empty invoice number causing HMAC calculation to fail
		},
		{
			name:          "nil receiver public key",
			source:        *sourcePubKey,
			linkPubKey:    bec.PublicKey{}, // Empty public key causing dedicated public key calculation to fail
			invoiceNumber: "valid-invoice",
		},
	}

	t.Run("Error Cases", func(t *testing.T) {
		for _, tt := range errorTests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				result, err := DeriveLinkedKey(&tt.source, &tt.linkPubKey, tt.invoiceNumber)
				assert.Error(t, err)
				assert.Nil(t, result)
			})
		}
	})
}

package type42

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/libsv/go-bk/bec"
	"github.com/stretchr/testify/assert"
)

func TestCalculateHMAC(t *testing.T) {
	tests := []struct {
		name            string
		pubSharedSecret []byte
		reference       string
		idx             int
		expected        []byte
		wantErr         bool
	}{
		{
			name:            "valid input",
			pubSharedSecret: []byte("shared_secret"),
			reference:       "reference",
			idx:             1,
			expected:        computeExpectedHMAC([]byte("shared_secret"), "reference", 1),
			wantErr:         false,
		},
		{
			name:            "different index",
			pubSharedSecret: []byte("shared_secret"),
			reference:       "reference",
			idx:             2,
			expected:        computeExpectedHMAC([]byte("shared_secret"), "reference", 2),
			wantErr:         false,
		},
		{
			name:            "empty shared secret",
			pubSharedSecret: []byte(""),
			reference:       "reference",
			idx:             1,
			expected:        computeExpectedHMAC([]byte(""), "reference", 1),
			wantErr:         false,
		},
		{
			name:            "empty reference",
			pubSharedSecret: []byte("shared_secret"),
			reference:       "",
			idx:             1,
			expected:        computeExpectedHMAC([]byte("shared_secret"), "", 1),
			wantErr:         false,
		},
		{
			name:            "both empty",
			pubSharedSecret: []byte(""),
			reference:       "",
			idx:             1,
			expected:        computeExpectedHMAC([]byte(""), "", 1),
			wantErr:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := calculateHMAC(tt.pubSharedSecret, fmt.Sprintf("%s-%d", tt.reference, tt.idx))
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, got)
			}
		})
	}
}

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

	validHMAC, _ := calculateHMAC(sourcePubKey.SerialiseCompressed(), "valid-invoice")
	validDerivedKey, _ := calculateDedicatedPublicKey(validHMAC, linkPubKey)

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

	t.Run("Valid Cases", func(t *testing.T) {
		for _, tt := range validTests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				result, err := DeriveLinkedKey(tt.source, tt.linkPubKey, tt.invoiceNumber)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			})
		}
	})

	t.Run("Error Cases", func(t *testing.T) {
		for _, tt := range errorTests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				result, err := DeriveLinkedKey(tt.source, tt.linkPubKey, tt.invoiceNumber)
				assert.Error(t, err)
				assert.Nil(t, result)
			})
		}
	})
}

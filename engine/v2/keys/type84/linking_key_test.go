package type84

import (
	"encoding/hex"
	"fmt"
	"testing"

	ec "github.com/bitcoin-sv/go-sdk/primitives/ec"
	assert "github.com/stretchr/testify/require"
)

func TestDeriveLinkedKey(t *testing.T) {
	sourcePubKeyHex := "027c1404c3ecb034053e6dd90bc68f7933284559c7d0763367584195a8796d9b0e"
	sourcePubKeyBytes, err := hex.DecodeString(sourcePubKeyHex)
	assert.NoError(t, err)
	sourcePubKey, err := ec.ParsePubKey(sourcePubKeyBytes)
	assert.NoError(t, err)

	linkPubKeyHex := "03a34e456deecb6e6e9237e63e5b7d045d1d2a456eb6be43de1ec4e9ac9a07b50d"
	linkPubKeyBytes, err := hex.DecodeString(linkPubKeyHex)
	assert.NoError(t, err)
	linkPubKey, err := ec.ParsePubKey(linkPubKeyBytes)
	assert.NoError(t, err)

	validHMAC, err := calculateHMAC(sourcePubKey.Compressed(), "valid-invoice")
	assert.NoError(t, err)
	validDerivedKey, err := calculateLinkedPublicKey(validHMAC, linkPubKey)
	assert.NoError(t, err)

	validTests := []struct {
		name           string
		source         ec.PublicKey
		linkPubKey     ec.PublicKey
		invoiceNumber  string
		expectedResult *ec.PublicKey
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
			t.Run(tt.name, func(t *testing.T) {
				result, err := DeriveLinkedKey(&tt.source, &tt.linkPubKey, tt.invoiceNumber)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			})
		}
	})

	errorTests := []struct {
		name          string
		source        ec.PublicKey
		linkPubKey    ec.PublicKey
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
			linkPubKey:    ec.PublicKey{}, // Empty public key causing dedicated public key calculation to fail
			invoiceNumber: "valid-invoice",
		},
	}

	t.Run("Error Cases", func(t *testing.T) {
		for _, tt := range errorTests {
			t.Run(tt.name, func(t *testing.T) {
				result, err := DeriveLinkedKey(&tt.source, &tt.linkPubKey, tt.invoiceNumber)
				assert.Error(t, err)
				assert.Nil(t, result)
			})
		}
	})
}

// Helper function to parse a public key from a hex string
func mustParsePubKey(hexKey string) *ec.PublicKey {
	keyBytes, err := hex.DecodeString(hexKey)
	if err != nil {
		panic(fmt.Sprintf("invalid hex key: %s", err))
	}
	key, err := ec.ParsePubKey(keyBytes)
	if err != nil {
		panic(fmt.Sprintf("invalid public key: %s", err))
	}
	return key
}

func TestDeriveLinkedKeyCases(t *testing.T) {
	validTests := []struct {
		name           string
		source         *ec.PublicKey
		linkPubKey     *ec.PublicKey
		invoiceNumber  string
		expectedResult *ec.PublicKey
	}{
		{
			name:           "case 1",
			source:         mustParsePubKey("033f9160df035156f1c48e75eae99914fa1a1546bec19781e8eddb900200bff9d1"),
			linkPubKey:     mustParsePubKey("02133b035cda4ba15f93b5fdde11c1f73eb9f1a79b60c6caa1c78e1c4c64ed72ce"),
			invoiceNumber:  "f3WCaUmnN9U=",
			expectedResult: mustParsePubKey("02d4049747553b4b956a419a7e1ddef418d57f317de8cc7d024d05c75f29309f26"),
		},
		{
			name:           "case 2",
			source:         mustParsePubKey("027775fa43959548497eb510541ac34b01d5ee9ea768de74244a4a25f7b60fae8d"),
			linkPubKey:     mustParsePubKey("02dfcbe35d95b55b5f3168ea8f12717e266ceddf88d04d2ff741272dfb0e542c2a"),
			invoiceNumber:  "2Ska++APzEc=",
			expectedResult: mustParsePubKey("02679645bc44a771154f66eb52ca93507c3e777997a774e87611a7b17d30d748ee"),
		},
		{
			name:           "case 3",
			source:         mustParsePubKey("0338d2e0d12ba645578b0955026ee7554889ae4c530bd7a3b6f688233d763e169f"),
			linkPubKey:     mustParsePubKey("023c1db4cce57a44c95e05f0ed499085a9ce8ac83bf80dd4c7b4658a9d6c4a122e"),
			invoiceNumber:  "cN/yQ7+k7pg=",
			expectedResult: mustParsePubKey("03c07884a8e8f02bb1cfb276edddfa8dade7654d8cb744e0ff9e4b0c5dffa63ae0"),
		},
		{
			name:           "case 4",
			source:         mustParsePubKey("02830212a32a47e68b98d477000bde08cb916f4d44ef49d47ccd4918d9aaabe9c8"),
			linkPubKey:     mustParsePubKey("0363ec92e374974e9c904875ff8daff50bfa22d8ec0a39891a31e5e760aca5e9cf"),
			invoiceNumber:  "m2/QAsmwaA4=",
			expectedResult: mustParsePubKey("024d369f8440e91fdd58577bd6bc471d868fdfc58e605eb2ec5a2362df589d43cd"),
		},
		{
			name:           "case 5",
			source:         mustParsePubKey("03f20a7e71c4b276753969e8b7e8b67e2dbafc3958d66ecba98dedc60a6615336d"),
			linkPubKey:     mustParsePubKey("02cef40bd6826d1b1960bc94094c0c0a19547b291c33b7a07380cbdc4580f3678b"),
			invoiceNumber:  "jgpUIjWFlVQ=",
			expectedResult: mustParsePubKey("03aa3ea0de12f44642a51593b2f8a0e0ff9813cdb4ab39ab8039654dce2ae5bdb3"),
		},
		{
			name:           "case 6",
			source:         mustParsePubKey("03f20a7e71c4b276753969e8b7e8b67e2dbafc3958d66ecba98dedc60a6615336d"),
			linkPubKey:     mustParsePubKey("0363ec92e374974e9c904875ff8daff50bfa22d8ec0a39891a31e5e760aca5e9cf"),
			invoiceNumber:  "jgpUIjWFlVQ=",
			expectedResult: mustParsePubKey("032ac4dccb237d777361a7fade10fc2d21b642df9b0e4634d3e27fdce131526e4f"),
		},
		{
			name:           "case 7",
			source:         mustParsePubKey("02830212a32a47e68b98d477000bde08cb916f4d44ef49d47ccd4918d9aaabe9c8"),
			linkPubKey:     mustParsePubKey("02cef40bd6826d1b1960bc94094c0c0a19547b291c33b7a07380cbdc4580f3678b"),
			invoiceNumber:  "jgpUIjWFlVQ=",
			expectedResult: mustParsePubKey("02caa04e2020a0642e5283c978257d2135c00882c049adcd71a71dffe9d7208a39"),
		},
		{
			name:           "case 8",
			source:         mustParsePubKey("02830212a32a47e68b98d477000bde08cb916f4d44ef49d47ccd4918d9aaabe9c8"),
			linkPubKey:     mustParsePubKey("02cef40bd6826d1b1960bc94094c0c0a19547b291c33b7a07380cbdc4580f3678b"),
			invoiceNumber:  "m2/QAsmwaA4=",
			expectedResult: mustParsePubKey("030467472149ac58d9d04e4182b03af99593a7af312623c4be7d96f2fde08f6421"),
		},
	}

	t.Run("Test Keys", func(t *testing.T) {
		for _, tt := range validTests {
			t.Run(tt.name, func(t *testing.T) {
				result, err := DeriveLinkedKey(tt.source, tt.linkPubKey, tt.invoiceNumber)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			})
		}
	})
}

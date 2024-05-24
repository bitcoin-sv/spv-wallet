package template

import (
	"encoding/hex"
	"testing"

	"github.com/libsv/go-bk/bec"
	"github.com/libsv/go-bk/crypto"
	"github.com/libsv/go-bt/bscript"
	assert "github.com/stretchr/testify/require"
)

func TestP2PKH(t *testing.T) {
	validTests := []struct {
		name     string
		satoshis uint64
		expected []P2PKHTemplate
	}{
		{
			name:     "valid input",
			satoshis: 1000,
			expected: []P2PKHTemplate{
				{
					Script:   "76a9fd88ac",
					Satoshis: 1000,
				},
			},
		},
		{
			name:     "zero satoshis",
			satoshis: 0,
			expected: []P2PKHTemplate{
				{
					Script:   "76a9fd88ac",
					Satoshis: 0,
				},
			},
		},
	}

	errorTests := []struct {
		name     string
		satoshis uint64
	}{
		{
			name:     "negative satoshis",
			satoshis: ^uint64(0), // Simulating a case that would cause an error, maximum uint64 value, bitwise NOT of 0 is -1
		},
	}

	t.Run("Valid Cases", func(t *testing.T) {
		for _, tt := range validTests {
			tt := tt // capture range variable
			t.Run(tt.name, func(t *testing.T) {
				got, err := P2PKH(tt.satoshis)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, got)
			})
		}
	})

	t.Run("Error Cases", func(t *testing.T) {
		for _, tt := range errorTests {
			tt := tt // capture range variable
			t.Run(tt.name, func(t *testing.T) {
				_, err := P2PKH(tt.satoshis)
				assert.Error(t, err)
			})
		}
	})
}

func TestEvaluate(t *testing.T) {
	pubKeyHex := "027c1404c3ecb034053e6dd90bc68f7933284559c7d0763367584195a8796d9b0e"
	pubKeyBytes, err := hex.DecodeString(pubKeyHex)
	assert.NoError(t, err)
	mockPublicKey, err := bec.ParsePubKey(pubKeyBytes, bec.S256())
	assert.NoError(t, err)
	mockPubKeyHash := crypto.Hash160(mockPublicKey.SerialiseCompressed())

	validTests := []struct {
		name      string
		script    []byte
		publicKey *bec.PublicKey
		expected  []byte
	}{
		{
			name:      "valid script with OP_PUBKEYHASH",
			script:    []byte{bscript.OpDUP, bscript.OpHASH160, bscript.OpPUBKEYHASH, bscript.OpEQUALVERIFY, bscript.OpCHECKSIG},
			publicKey: mockPublicKey,
			expected:  append([]byte{bscript.OpDUP, bscript.OpHASH160, bscript.OpDATA20}, append(mockPubKeyHash, bscript.OpEQUALVERIFY, bscript.OpCHECKSIG)...),
		},
		{
			name:      "valid script with OP_PUBKEY",
			script:    []byte{bscript.OpDUP, bscript.OpHASH160, bscript.OpPUBKEY, bscript.OpEQUALVERIFY, bscript.OpCHECKSIG},
			publicKey: mockPublicKey,
			expected:  append([]byte{bscript.OpDUP, bscript.OpHASH160, bscript.OpDATA20}, append(mockPubKeyHash, bscript.OpEQUALVERIFY, bscript.OpCHECKSIG)...),
		},
		{
			name:      "valid script without OP_PUBKEYHASH or OP_PUBKEY",
			script:    []byte{bscript.OpDUP, bscript.OpHASH160, bscript.OpEQUALVERIFY, bscript.OpCHECKSIG},
			publicKey: mockPublicKey,
			expected:  []byte{bscript.OpDUP, bscript.OpHASH160, bscript.OpEQUALVERIFY, bscript.OpCHECKSIG},
		},
	}

	t.Run("Valid Cases", func(t *testing.T) {
		for _, tt := range validTests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				result := Evaluate(tt.script, tt.publicKey)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expected, result)
			})
		}
	})
}

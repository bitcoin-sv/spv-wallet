package template

import (
	"encoding/hex"
	"testing"

	ec "github.com/bitcoin-sv/go-sdk/primitives/ec"
	crypto "github.com/bitcoin-sv/go-sdk/primitives/hash"
	script "github.com/bitcoin-sv/go-sdk/script"
	assert "github.com/stretchr/testify/require"
)

func TestP2PKH(t *testing.T) {
	validTests := []struct {
		name     string
		satoshis uint64
		expected *OutputTemplate
	}{
		{
			name:     "valid input",
			satoshis: 1000,
			expected: &OutputTemplate{
				Script:   "76a9fd88ac",
				Satoshis: 1000,
			},
		},
		{
			name:     "zero satoshis",
			satoshis: 1,
			expected: &OutputTemplate{
				Script:   "76a9fd88ac",
				Satoshis: 1,
			},
		},
	}

	t.Run("Valid Cases", func(t *testing.T) {
		for _, tt := range validTests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := P2PKH(tt.satoshis)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, got)
			})
		}
	})

	t.Run("Valid Cases", func(t *testing.T) {
		for _, tt := range validTests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := P2PKH(tt.satoshis)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, got)
			})
		}
	})

	errorTests := []struct {
		name     string
		satoshis uint64
	}{
		{
			name:     "negative satoshis",
			satoshis: ^uint64(0), // Simulating a case that would cause an error, maximum uint64 value, bitwise NOT of 0 is -1
		},
		{
			name:     "zero satoshis",
			satoshis: 0,
		},
	}

	t.Run("Error Cases", func(t *testing.T) {
		for _, tt := range errorTests {
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
	mockPublicKey, err := ec.ParsePubKey(pubKeyBytes)
	assert.NoError(t, err)
	mockPubKeyHash := crypto.Hash160(mockPublicKey.Compressed())

	t.Run("Valid Cases", func(t *testing.T) {
		validTests := []struct {
			name      string
			script    []byte
			publicKey *ec.PublicKey
			expected  []byte
		}{
			{
				name:      "valid script with OP_PUBKEYHASH",
				script:    []byte{script.OpDUP, script.OpHASH160, script.OpPUBKEYHASH, script.OpEQUALVERIFY, script.OpCHECKSIG},
				publicKey: mockPublicKey,
				expected:  append([]byte{script.OpDUP, script.OpHASH160, script.OpDATA20}, append(mockPubKeyHash, script.OpEQUALVERIFY, script.OpCHECKSIG)...),
			},
			{
				name:      "valid script without OP_PUBKEYHASH or OP_PUBKEY",
				script:    []byte{script.OpDUP, script.OpHASH160, script.OpEQUALVERIFY, script.OpCHECKSIG},
				publicKey: mockPublicKey,
				expected:  []byte{script.OpDUP, script.OpHASH160, script.OpEQUALVERIFY, script.OpCHECKSIG},
			},
			{
				name:      "script with OP_PUSHDATA1 and hex data matching PUBKEY and PUBKEYHASH opcodes",
				script:    []byte{script.OpPUSHDATA1, 1, script.OpPUBKEYHASH, script.OpADD, script.OpPUSHDATA1, 1, script.OpPUBKEY, script.OpEQUALVERIFY},
				publicKey: mockPublicKey,
				expected:  []byte{script.OpPUSHDATA1, 1, script.OpPUBKEYHASH, script.OpADD, script.OpPUSHDATA1, 1, script.OpPUBKEY, script.OpEQUALVERIFY},
			},
			{
				name:      "empty script",
				script:    []byte{},
				publicKey: mockPublicKey,
				expected:  []byte{},
			},
			{
				name:      "script with only valid push data",
				script:    []byte{script.OpPUSHDATA1, 2, 0xaa, 0xbb},
				publicKey: mockPublicKey,
				expected:  []byte{script.OpPUSHDATA1, 2, 0xaa, 0xbb},
			},
		}

		for _, tt := range validTests {
			t.Run(tt.name, func(t *testing.T) {
				result, err := Evaluate(tt.script, tt.publicKey)
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expected, result)
			})
		}
	})

	t.Run("Invalid Cases", func(t *testing.T) {
		invalidTests := []struct {
			name      string
			script    []byte
			publicKey *ec.PublicKey
		}{
			{
				name:      "invalid script",
				script:    []byte{0xFF}, // Invalid opcode
				publicKey: mockPublicKey,
			},
			{
				name:      "valid script with OP_PUBKEY",
				script:    []byte{script.OpDUP, script.OpHASH160, script.OpPUBKEY, script.OpEQUALVERIFY, script.OpCHECKSIG},
				publicKey: mockPublicKey,
			},
		}

		for _, tt := range invalidTests {
			t.Run(tt.name, func(t *testing.T) {
				_, err := Evaluate(tt.script, tt.publicKey)
				assert.Error(t, err)
			})
		}
	})
}

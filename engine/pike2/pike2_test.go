package pike2

import (
	"encoding/hex"
	"testing"

	"github.com/libsv/go-bk/bec"
	assert "github.com/stretchr/testify/require"
)

func TestCreatePikeOutput(t *testing.T) {
	tests := []struct {
		name     string
		satoshis uint64
		wantErr  bool
		expected []PikeOutputsTemplate
	}{
		{
			name:     "valid input",
			satoshis: 1000,
			wantErr:  false,
			expected: []PikeOutputsTemplate{
				{
					Script:   "76a9fd88ac",
					Satoshis: 1000,
				},
			},
		},
		{
			name:     "zero satoshis",
			satoshis: 0,
			wantErr:  false,
			expected: []PikeOutputsTemplate{
				{
					Script:   "76a9fd88ac",
					Satoshis: 0,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreatePikeOutputsTemplate(tt.satoshis)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, got)
			}
		})
	}
}

func TestNewPrivateKey(t *testing.T) {
	tests := []struct {
		name     string
		value    int64
		expected string
	}{
		{
			name:     "value 1",
			value:    1,
			expected: "0000000000000000000000000000000000000000000000000000000000000001",
		},
		{
			name:     "value 255",
			value:    255,
			expected: "00000000000000000000000000000000000000000000000000000000000000ff",
		},
		{
			name:     "value 256",
			value:    256,
			expected: "0000000000000000000000000000000000000000000000000000000000000100",
		},
		{
			name:     "value max int64",
			value:    int64(9223372036854775807),
			expected: "0000000000000000000000000000000000000000000000007fffffffffffffff",
		},
		{
			name:     "value 0",
			value:    0,
			expected: "0000000000000000000000000000000000000000000000000000000000000000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			privKey := newPrivateKey(tt.value)
			privKeyBytes := privKey.Serialise()
			privKeyHex := hex.EncodeToString(privKeyBytes)
			assert.Equal(t, tt.expected, privKeyHex)
		})
	}
}

func TestCalculatePubSharedSecret(t *testing.T) {
	// Helper function to generate a random public key
	generateRandomPublicKey := func() *bec.PublicKey {
		privKey, _ := bec.NewPrivateKey(bec.S256())
		return privKey.PubKey()
	}

	tests := []struct {
		name         string
		senderPubKey *bec.PublicKey
	}{
		{
			name:         "valid public key",
			senderPubKey: generateRandomPublicKey(),
		},
		{
			name:         "another valid public key",
			senderPubKey: generateRandomPublicKey(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculatePubSharedSecret(tt.senderPubKey)
			assert.NotNil(t, got, "The returned public key should not be nil")
			assert.NotNil(t, got.X, "The X coordinate of the public key should not be nil")
			assert.NotNil(t, got.Y, "The Y coordinate of the public key should not be nil")
		})
	}
}

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
			got, err := calculateHMAC(tt.pubSharedSecret, tt.reference, tt.idx)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, got)
			}
		})
	}
}

func TestCalculateDedicatedPublicKey(t *testing.T) {
	tests := []struct {
		name        string
		hmacResult  []byte
		receiverPub *bec.PublicKey
		wantErr     bool
		errMsg      string
	}{
		{
			name:        "valid input",
			hmacResult:  sha256Hash("test_hmac_1"),
			receiverPub: generateRandomPublicKey(),
			wantErr:     false,
		},
		{
			name:        "another valid input",
			hmacResult:  sha256Hash("test_hmac_2"),
			receiverPub: generateRandomPublicKey(),
			wantErr:     false,
		},
		{
			name:        "empty HMAC result",
			hmacResult:  []byte{},
			receiverPub: generateRandomPublicKey(),
			wantErr:     true,
			errMsg:      "HMAC result is empty",
		},
		{
			name:        "nil receiver public key",
			hmacResult:  sha256Hash("test_hmac_3"),
			receiverPub: nil,
			wantErr:     true,
			errMsg:      "receiver public key is nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := calculateDedicatedPublicKey(tt.hmacResult, tt.receiverPub)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				if got != nil {
					assert.True(t, bec.S256().IsOnCurve(got.X, got.Y))
				}
			}
		})
	}
}

func TestGenerateLockingScriptsFromTemplates(t *testing.T) {
	tests := []struct {
		name            string
		outputsTemplate []PikeOutputsTemplate
		senderPubKey    *bec.PublicKey
		receiverPubKey  *bec.PublicKey
		reference       string
		wantErr         bool
	}{
		{
			name: "valid input",
			outputsTemplate: []PikeOutputsTemplate{
				{Script: "76a91489abcdefabbaabbaabbaabbaabbaabbaabbaabba88ac", Satoshis: 1000}, // Valid P2PKH script
			},
			senderPubKey:   generateRandomPublicKey(),
			receiverPubKey: generateRandomPublicKey(),
			reference:      "test_ref",
			wantErr:        false,
		},
		{
			name: "invalid hex script",
			outputsTemplate: []PikeOutputsTemplate{
				{Script: "invalid_hex", Satoshis: 1000},
			},
			senderPubKey:   generateRandomPublicKey(),
			receiverPubKey: generateRandomPublicKey(),
			reference:      "test_ref",
			wantErr:        true,
		},
		{
			name:            "empty outputs template",
			outputsTemplate: []PikeOutputsTemplate{},
			senderPubKey:    generateRandomPublicKey(),
			receiverPubKey:  generateRandomPublicKey(),
			reference:       "test_ref",
			wantErr:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateLockingScriptsFromTemplates(tt.outputsTemplate, tt.senderPubKey, tt.receiverPubKey, tt.reference)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
			}
		})
	}
}

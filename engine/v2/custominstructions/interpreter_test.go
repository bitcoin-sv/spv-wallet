package custominstructions

import (
	"testing"

	primitives "github.com/bitcoin-sv/go-sdk/primitives/ec"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"github.com/stretchr/testify/require"
)

func TestLockingScriptInterpreter(t *testing.T) {
	tests := map[string]struct {
		customInstructions bsv.CustomInstructions
		expectAddress      string
	}{
		"standard": {
			customInstructions: bsv.CustomInstructions{
				{
					Type:        Type42,
					Instruction: "1-paymail_pki-test@example.com_0",
				},
				{
					Type:        Type42,
					Instruction: "1-destination-6a5dbb7df22a265de809c35dd8d703c1",
				},
				{
					Type:        Sign,
					Instruction: "P2PKH",
				},
			},
			expectAddress: "18RAjvMrT1HMzejGn9qyP4zb1c5BgHidRa",
		},
		"with no explicit sign method": {
			customInstructions: bsv.CustomInstructions{
				{
					Type:        Type42,
					Instruction: "1-paymail_pki-test@example.com_0",
				},
				{
					Type:        Type42,
					Instruction: "1-destination-6a5dbb7df22a265de809c35dd8d703c1",
				},
			},
			expectAddress: "18RAjvMrT1HMzejGn9qyP4zb1c5BgHidRa",
		},
		"with empty instructions": {
			customInstructions: bsv.CustomInstructions{},
			expectAddress:      "1GtetpoX4eraGj7FgMhwdRhrd13Xn96USN",
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			pubKey := makePubKey(t, "033014c226b8fe8260e21e75479a47a654e7b631b3bd13484d85c484f7791aa75b")

			// and:
			processor := NewAddressInterpreter()

			// when:
			res, err := processor.Process(pubKey, test.customInstructions)

			// then:
			require.NoError(t, err)

			// and:
			require.NoError(t, err)
			require.Equal(t, test.expectAddress, res.Address.AddressString)
		})
	}
}

func makePubKey(t *testing.T, pubDERHex string) *primitives.PublicKey {
	t.Helper()
	pk, err := primitives.PublicKeyFromString(pubDERHex)
	require.NoError(t, err)
	return pk
}

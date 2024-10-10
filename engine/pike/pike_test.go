package pike

import (
	"encoding/hex"
	"testing"

	ec "github.com/bitcoin-sv/go-sdk/primitives/ec"
	"github.com/bitcoin-sv/spv-wallet/engine/script/template"
	assert "github.com/stretchr/testify/require"
)

func TestGenerateLockingScriptsFromTemplates(t *testing.T) {
	// Define sample data
	senderPubKeyHex := "027c1404c3ecb034053e6dd90bc68f7933284559c7d0763367584195a8796d9b0e"
	senderPubKeyBytes, err := hex.DecodeString(senderPubKeyHex)
	assert.NoError(t, err)
	senderPubKey, err := ec.ParsePubKey(senderPubKeyBytes)
	assert.NoError(t, err)

	receiverPubKeyHex := "03a34e456deecb6e6e9237e63e5b7d045d1d2a456eb6be43de1ec4e9ac9a07b50d"
	receiverPubKeyBytes, err := hex.DecodeString(receiverPubKeyHex)
	assert.NoError(t, err)
	receiverPubKey, err := ec.ParsePubKey(receiverPubKeyBytes)
	assert.NoError(t, err)

	outputsTemplate := []*template.OutputTemplate{
		{Script: "76a914000000000000000000000000000000000000000088ac"},
		{Script: "76a914111111111111111111111111111111111111111188ac"},
	}

	t.Run("Valid Case", func(t *testing.T) {
		lockingScripts, err := GenerateLockingScriptsFromTemplates(outputsTemplate, senderPubKey, receiverPubKey, "test-reference")
		assert.NoError(t, err)
		assert.Len(t, lockingScripts, len(outputsTemplate))
		assert.Equal(t, outputsTemplate[0].Script, lockingScripts[0])
		assert.Equal(t, outputsTemplate[1].Script, lockingScripts[1])
	})

	t.Run("Invalid Template Script", func(t *testing.T) {
		invalidTemplate := []*template.OutputTemplate{
			{Script: "invalid-hex-string"}, // Invalid hex string
		}
		lockingScripts, err := GenerateLockingScriptsFromTemplates(invalidTemplate, senderPubKey, receiverPubKey, "test-reference")
		assert.Error(t, err)
		assert.Nil(t, lockingScripts)
	})
}

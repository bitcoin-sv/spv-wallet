package pike_test

import (
	"encoding/hex"
	"fmt"

	"github.com/bitcoin-sv/spv-wallet/engine/pike"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/libsv/go-bk/bec"
)

func Example_generateLockingScripts() {
	// Example sender's public key (replace with actual sender's public key)
	senderPublicKeyHex := "034252e5359a1de3b8ec08e6c29b80594e88fb47e6ae9ce65ee5a94f0d371d2cde"
	senderPublicKeyBytes, err := hex.DecodeString(senderPublicKeyHex)
	if err != nil {
		panic(err)
	}
	senderPubKey, err := bec.ParsePubKey(senderPublicKeyBytes, bec.S256())
	if err != nil {
		panic(err)
	}

	receiverPublicKeyHex := "027c1404c3ecb034053e6dd90bc68f7933284559c7d0763367584195a8796d9b0e"
	receiverPublicKeyBytes, err := hex.DecodeString(receiverPublicKeyHex)
	if err != nil {
		panic(err)
	}
	receiverPubKey, err := bec.ParsePubKey(receiverPublicKeyBytes, bec.S256())
	if err != nil {
		panic(err)
	}

	// Example usage of GenerateOutputsTemplate
	outputsTemplate, err := pike.GenerateOutputsTemplate(10000)
	if err != nil {
		panic(spverrors.Wrapf(err, "Error generating outputs template"))
	}

	// Example usage of GenerateLockingScriptsFromTemplates
	lockingScripts, err := pike.GenerateLockingScriptsFromTemplates(outputsTemplate, senderPubKey, receiverPubKey, "reference")
	if err != nil {
		panic(spverrors.Wrapf(err, "Error generating locking scripts"))
	}

	for _, script := range lockingScripts {
		fmt.Println("Locking Script:", script)
	}

	// Output:
	// Locking Script: 76a9147327490be831259f38b0f9ab019413e51d1b40c688ac
}

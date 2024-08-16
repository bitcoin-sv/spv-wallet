package main

import (
	"encoding/hex"
	"fmt"

	"github.com/bitcoin-sv/spv-wallet/engine/pike"
	"github.com/bitcoin-sv/spv-wallet/engine/script/template"
	"github.com/libsv/go-bk/bec"
)

func main() {
	// Example sender's public key (replace with actual sender's public key)
	// generating keys
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

	// example of usage pike.GenerateOutputsTemplate
	outputsTemplate, _ := pike.GenerateOutputsTemplate(10000)
	fmt.Println(formatOutputs(outputsTemplate))

	lockingScripts, err := pike.GenerateLockingScriptsFromTemplates(outputsTemplate, senderPubKey, receiverPubKey, "reference")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for _, script := range lockingScripts {
		fmt.Println("Locking Script:", script)
	}

}

// Helper function to format the outputs into a string
func formatOutputs(outputs []*template.OutputTemplate) string {
	var result string
	for i, output := range outputs {
		result += fmt.Sprintf("Output %d: %v\n", i+1, output)
	}
	return result
}

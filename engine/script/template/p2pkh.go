// Package template provides a collection of functions and types for working with script templates.
package template

import (
	"encoding/hex"
	"sync"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/libsv/go-bt/v2/bscript"
)

var (
	scriptHex string
	once      sync.Once
)

func initScriptHex() {
	opcodes := []byte{
		bscript.OpDUP,
		bscript.OpHASH160,
		bscript.OpPUBKEYHASH,
		bscript.OpEQUALVERIFY,
		bscript.OpCHECKSIG,
	}

	// Convert opcodes to hexadecimal string
	scriptHex = hex.EncodeToString(opcodes)
}

// OutputTemplate represents the script and satoshis for a Pike output
type OutputTemplate struct {
	Script   string `json:"script"`
	Satoshis uint64 `json:"satoshis"`
}

// P2PKH creates a single output with the PIKE template
func P2PKH(satoshis uint64) (*OutputTemplate, error) {

	if satoshis == 0 {
		return nil, spverrors.Newf("satoshis cannot be zero")
	}
	if satoshis == ^uint64(0) {
		return nil, spverrors.Newf("invalid satoshis")
	}

	// Initialize the scriptHex once
	once.Do(initScriptHex)

	// Create and return the PikeOutputsTemplate
	return &OutputTemplate{
		Script:   scriptHex,
		Satoshis: satoshis,
	}, nil
}

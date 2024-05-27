// Package template provides a collection of functions and types for working with script templates.
package template

import (
	"encoding/hex"
	"errors"

	"github.com/libsv/go-bt/v2/bscript"
)

// OutputTemplate represents the script and satoshis for a Pike output
type OutputTemplate struct {
	Script   string `json:"script"`
	Satoshis uint64 `json:"satoshis"`
}

// P2PKH creates a single output with the PIKE template
func P2PKH(satoshis uint64) (*OutputTemplate, error) {

	if satoshis == 0 {
		return nil, errors.New("satoshis cannot be zero")
	}
	if satoshis == ^uint64(0) {
		return nil, errors.New("invalid satoshis")
	}

	opcodes := []byte{
		bscript.OpDUP,
		bscript.OpHASH160,
		bscript.OpPUBKEYHASH,
		bscript.OpEQUALVERIFY,
		bscript.OpCHECKSIG,
	}

	// Convert opcodes to hexadecimal string
	scriptHex := hex.EncodeToString(opcodes)

	// Check if the conversion was successful
	if scriptHex == "" {
		return nil, errors.New("failed to create script hex")
	}

	// Create and return the PikeOutputsTemplate
	return &OutputTemplate{
		Script:   scriptHex,
		Satoshis: satoshis,
	}, nil
}

// Package template provides a collection of functions and types for working with script templates.
package template

import (
	"encoding/hex"
	"errors"

	"github.com/libsv/go-bk/bec"
	"github.com/libsv/go-bk/crypto"
	"github.com/libsv/go-bt/v2/bscript"
	"github.com/libsv/go-bt/v2/bscript/interpreter"
)

// P2PKHTemplate represents the script and satoshis for a Pike output
type P2PKHTemplate struct {
	Script   string `json:"script"`
	Satoshis uint64 `json:"satoshis"`
}

// P2PKH creates a single output with the PIKE template
func P2PKH(satoshis uint64) ([]P2PKHTemplate, error) {

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
	return []P2PKHTemplate{
		{
			Script:   scriptHex,
			Satoshis: satoshis,
		},
	}, nil
}

// Evaluate processes a given Bitcoin script by parsing it, replacing certain opcodes
// with the public key hash, and returning the resulting script as a byte array.
// Will replace any OP_PUBKEYHASH or OP_PUBKEY
//
// Parameters:
// - script: A byte array representing the input script.
// - dPK: A pointer to a bec.PublicKey which provides the dedicated public key to be used in the evaluation.
//
// Returns:
// - A byte array representing the evaluated script, or nil if an error occurs.
func Evaluate(script []byte, dPK *bec.PublicKey) []byte {
	s := bscript.Script(script)

	parser := interpreter.DefaultOpcodeParser{}
	parsedScript, err := parser.Parse(&s)
	if err != nil {
		return nil
	}

	// Validate parsed opcodes
	for _, op := range parsedScript {
		if op.Value() == 0xFF {
			return nil
		}
	}

	// Serialize the public key to compressed format
	dPKBytes := dPK.SerialiseCompressed()

	// Apply Hash160 (SHA-256 followed by RIPEMD-160) to the compressed public key
	dPKBytesCompressed := crypto.Hash160(dPKBytes)

	// Create a new script with the public key hash
	newScript := new(bscript.Script)
	if err := newScript.AppendPushData(dPKBytesCompressed); err != nil {
		return nil
	}

	// Parse the public key hash script
	pkhParsed, err := parser.Parse(newScript)
	if err != nil {
		return nil
	}

	// Replace OP_PUBKEYHASH with the actual public key hash
	evaluated := make([]interpreter.ParsedOpcode, 0, len(parsedScript))
	for _, op := range parsedScript {
		switch op.Value() {
		case bscript.OpPUBKEYHASH:
			evaluated = append(evaluated, pkhParsed...)
		case bscript.OpPUBKEY:
			evaluated = append(evaluated, pkhParsed...) // Currently just appending the same key
		default:
			evaluated = append(evaluated, op)
		}
	}

	// Unparse the evaluated opcodes back into a script
	finalScript, err := parser.Unparse(evaluated)
	if err != nil {
		return nil
	}

	// Cast *bscript.Script back to []byte
	return []byte(*finalScript)
}

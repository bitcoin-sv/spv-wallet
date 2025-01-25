package bsv

// CustomInstruction represents a custom instruction how to unlock a UTXO.
type CustomInstruction struct {
	Type        string `json:"type"`
	Instruction string `json:"instruction"`
}

// CustomInstructions is a slice of CustomInstruction.
type CustomInstructions []CustomInstruction

package bsv

import "fmt"

// Outpoint is a struct that represents a pair consisting of a transaction ID and an output index
// This represents a specific unspent transaction output (UTXO)
type Outpoint struct {
	TxID string
	Vout uint32
}

// String returns a string representation of outpoint
func (o *Outpoint) String() string {
	return fmt.Sprintf("%s-%d", o.TxID, o.Vout)
}

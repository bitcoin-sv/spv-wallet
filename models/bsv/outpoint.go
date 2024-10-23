package bsv

// Outpoint is a struct that represents a pair consisting of a transaction ID and an output index
// This represents a specific unspent transaction output (UTXO)
type Outpoint struct {
	TxID string
	Vout uint32
}

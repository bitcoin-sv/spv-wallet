package bsv

import (
	"fmt"
	"math"
)

// Outpoint is a struct that represents a pair consisting of a transaction ID and an output index
// This represents a specific unspent transaction output (UTXO)
type Outpoint struct {
	TxID string
	Vout uint32
}

// String returns a string representation of outpoint
func (o Outpoint) String() string {
	return fmt.Sprintf("%s-%d", o.TxID, o.Vout)
}

// OutpointFromString creates an Outpoint from a string
func OutpointFromString(s string) (Outpoint, error) {
	if s == "" {
		return Outpoint{}, fmt.Errorf("empty string")
	}

	var txID string
	var voutTmp int
	n, err := fmt.Sscanf(s, "%64s-%d", &txID, &voutTmp)
	if err != nil {
		return Outpoint{}, fmt.Errorf("invalid outpoint format: %w", err)
	} else if n != 2 {
		return Outpoint{}, fmt.Errorf("invalid outpoint format")
	}

	vout, err := toUint32(voutTmp)
	if err != nil {
		return Outpoint{}, err
	}

	return Outpoint{TxID: txID, Vout: vout}, nil
}

func toUint32(value int) (uint32, error) {
	if value < 0 || value > math.MaxUint32 {
		return 0, fmt.Errorf("invalid vout")
	}
	return uint32(value), nil
}

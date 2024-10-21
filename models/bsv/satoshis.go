package bsv

import (
	"errors"
)

type Satoshis uint64

// SatoshisFromInt creates a new Satoshis from an integer.
func SatoshisFromInt(s int) (Satoshis, error) {
	if s < 0 {
		return 0, errors.New("value cannot be negative")
	}
	return Satoshis(s), nil //nolint:gosec
}

package bsv

import "fmt"

// FeeUnit displays the amount of Satoshis neededz
// for a specific amount of Bytes in a transaction
// see https://github.com/bitcoin-sv-specs/brfc-misc/tree/master/feespec
// Imported from deprecated go-bt library
type FeeUnit struct {
	Satoshis int `json:"satoshis"` // Fee in satoshis of the amount of Bytes
	Bytes    int `json:"bytes"`    // Number of bytes that the Fee covers
}

// IsLowerThan compare two fee units
func (f *FeeUnit) IsLowerThan(other *FeeUnit) bool {
	return float64(f.Satoshis)/float64(f.Bytes) < float64(other.Satoshis)/float64(other.Bytes)
}

// String returns the fee unit as a string
func (f *FeeUnit) String() string {
	return fmt.Sprintf("FeeUnit(%d satoshis / %d bytes)", f.Satoshis, f.Bytes)
}

// IsZero returns true if the fee unit suggest no fees (free)
func (f *FeeUnit) IsZero() bool {
	return f.Satoshis == 0
}

// IsValid returns true if the Bytes in fee are greater than 0
func (f *FeeUnit) IsValid() bool {
	return f.Bytes > 0
}

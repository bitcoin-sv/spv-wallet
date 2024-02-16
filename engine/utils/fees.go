package utils

import (
	"encoding/hex"
	"fmt"

	"github.com/libsv/go-bt/v2"
)

// FeeUnit fee unit imported from go-bt/v2
type FeeUnit bt.FeeUnit

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

// ValidFees filters out invalid fees from a list of fee units
func ValidFees(feeUnits []FeeUnit) []FeeUnit {
	validFees := []FeeUnit{}
	for _, fee := range feeUnits {
		if fee.IsValid() {
			validFees = append(validFees, fee)
		}
	}
	return validFees
}

// LowestFee get the lowest fee from a list of fee units, if defaultValue exists and none is found, return defaultValue
func LowestFee(feeUnits []FeeUnit, defaultValue *FeeUnit) *FeeUnit {
	validFees := ValidFees(feeUnits)
	if len(validFees) == 0 {
		return defaultValue
	}
	minFee := validFees[0]
	for i := 1; i < len(validFees); i++ {
		if validFees[i].IsLowerThan(&minFee) {
			minFee = validFees[i]
		}
	}
	return &minFee
}

// GetInputSizeForType get an estimated size for the input based on the type
func GetInputSizeForType(inputType string) uint64 {
	switch inputType {
	case ScriptTypePubKeyHash:
		// 32 bytes txID
		// + 4 bytes vout index
		// + 1 byte script length
		// + 107 bytes script pub key
		// + 4 bytes nSequence
		return 148
	}

	return 500
}

// GetOutputSize get an estimated size for the output based on the type
func GetOutputSize(lockingScript string) uint64 {
	if lockingScript != "" {
		size, _ := hex.DecodeString(lockingScript)
		if size != nil {
			return uint64(len(size)) + 9 // 9 bytes = 8 bytes value, 1 byte length
		}
	}

	return 500
}

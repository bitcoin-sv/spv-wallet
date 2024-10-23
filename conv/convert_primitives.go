package conv

import (
	"math"

	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
)

// ConvertInt64ToUint32 will convert an int64 to a uint32, with range checks
func ConvertInt64ToUint32(value int64) (uint32, error) {
	if value < 0 || value > math.MaxUint32 {
		return 0, spverrors.ErrInvalidUint32
	}
	return uint32(value), nil
}

// ConvertUint32ToInt64 will convert a uint32 to an int64 (safe as uint32 fits into int64)
func ConvertUint32ToInt64(value uint32) int64 {
	return int64(value)
}

// ConvertUint64ToInt64 will convert a uint64 to an int64, with range checks
func ConvertUint64ToInt64(value uint64) (int64, error) {
	if value > math.MaxInt64 {
		return 0, spverrors.ErrInvalidInt64
	}
	return int64(value), nil
}

// ConvertInt64ToUint64 will convert an int64 to a uint64, with range checks
func ConvertInt64ToUint64(value int64) (uint64, error) {
	if value < 0 {
		return 0, spverrors.ErrInvalidUint64
	}
	return uint64(value), nil
}

// ConvertUint64ToInt will convert a uint64 to an int, with range checks
func ConvertUint64ToInt(value uint64) (int, error) {
	if value > math.MaxInt {
		return 0, spverrors.ErrInvalidInt
	}
	return int(value), nil
}

// ConvertIntToUint64 will convert an int to a uint64, with range checks
func ConvertIntToUint64(value int) (uint64, error) {
	if value < 0 {
		return 0, spverrors.ErrInvalidUint64
	}
	return uint64(value), nil
}

// ConvertIntToUint32 will convert an int to a uint32, with range checks
func ConvertIntToUint32(value int) (uint32, error) {
	if value < 0 || value > math.MaxUint32 {
		return 0, spverrors.ErrInvalidUint32
	}
	return uint32(value), nil
}

// SafeVarIntToInt will convert a VarInt to an int, with range checks
func SafeVarIntToInt(varInt *sdk.VarInt) (int, error) {
	if varInt == nil {
		return 0, spverrors.ErrInvalidInt
	}
	i := uint64(*varInt)
	// Ensure VarInt is not negative and within the range of int
	if i > uint64(math.MaxInt) {
		return 0, spverrors.ErrInvalidInt
	}
	// Convert the VarInt to an int
	return int(i), nil
}

// ConvertToIntUint64 will convert an int to a uint64, with range checks
func ConvertToIntUint64(value int) (uint64, error) {
	if value < 0 {
		return 0, spverrors.ErrInvalidUint64
	}
	return uint64(value), nil
}

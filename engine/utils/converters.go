package utils

import (
	"math"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
)

// ConvertInt64ToUint32 converts an int to a uint32
func ConvertInt64ToUint32(num int64) (uint32, error) {
	if num < 0 {
		return 0, spverrors.Newf("cannot convert negative number to uint32")
	}
	if num > math.MaxUint32 {
		return 0, spverrors.Newf("number is too large to convert to uint32")
	}
	return uint32(num), nil
}

// ConvertIntToUint32 converts an int to a uint32
func ConvertIntToUint32(num int) (uint32, error) {
	if num < 0 {
		return 0, spverrors.Newf("cannot convert negative number to uint32")
	}
	if num > math.MaxUint32 {
		return 0, spverrors.Newf("number is too large to convert to uint32")
	}
	return uint32(num), nil
}

// ConvertInt64ToUInt64 converts an int to a uint64
func ConvertInt64ToUInt64(num int64) (uint64, error) {
	if num < 0 {
		return 0, spverrors.Newf("cannot convert negative number to uint64")
	}
	return uint64(num), nil
}

// ConvertIntToUInt64 converts an int to a uint64
func ConvertIntToUInt64(num int) (uint64, error) {
	if num < 0 {
		return 0, spverrors.Newf("cannot convert negative number to uint64")
	}
	return uint64(num), nil
}

// ShouldConvertInt64ToUInt64 converts an int to a uint64 and panics if there is an error
func ShouldConvertInt64ToUInt64(num int64) uint64 {
	value, err := ConvertInt64ToUInt64(num)
	if err != nil {
		panic(err)
	}
	return value
}

// ShouldConvertIntToUInt64 converts an int to a uint64 and panics if there is an error
func ShouldConvertIntToUInt64(num int) uint64 {
	value, err := ConvertIntToUInt64(num)
	if err != nil {
		panic(err)
	}
	return value
}

// ConvertIntToUInt32 converts an int to an int32
func ConvertIntToUInt32(num int) (uint32, error) {
	if num < 0 {
		return 0, spverrors.Newf("cannot convert negative number to uint32")
	}
	if num > math.MaxUint32 {
		return 0, spverrors.Newf("number is too large to convert to uint32")
	}
	return uint32(num), nil
}

// ConvertUInt64ToInt64 converts an uint64 to an int64
func ConvertUInt64ToInt64(num uint64) (int64, error) {
	if num > math.MaxInt64 {
		return 0, spverrors.Newf("number is too large to convert to int64")
	}
	return int64(num), nil
}

// ConvertUInt64ToInt converts an uint64 to an int
func ConvertUInt64ToInt(num uint64) (int, error) {
	if num > math.MaxInt {
		return 0, spverrors.Newf("number is too large to convert to int")
	}
	return int(num), nil
}

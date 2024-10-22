package conversionkit

import (
	"math"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
)

// ConvertToUint64 is a generic function that safely converts signed integers to uint64
func ConvertToUint64[T any](value T) (uint64, error) {
	switch v := any(value).(type) {
	case int:
		if v < 0 {
			return 0, spverrors.ErrInvalidUint64Value
		}
		return uint64(v), nil
	case int64:
		if v < 0 {
			return 0, spverrors.ErrInvalidUint64Value
		}
		return uint64(v), nil
	case int32:
		if v < 0 {
			return 0, spverrors.ErrInvalidUint64Value
		}
		return uint64(v), nil
	case int16:
		if v < 0 {
			return 0, spverrors.ErrInvalidUint64Value
		}
		return uint64(v), nil
	case int8:
		if v < 0 {
			return 0, spverrors.ErrInvalidUint64Value
		}
		return uint64(v), nil
	default:
		return 0, spverrors.ErrUnsupportedDestinationType
	}
}

// ConvertToInt64 is a generic function that converts numeric values to int64
func ConvertToInt64[T any](value T) (int64, error) {
	switch v := any(value).(type) {
	case int:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int16:
		return int64(v), nil
	case int8:
		return int64(v), nil
	case uint:
		if v > uint(math.MaxInt64) {
			return math.MaxInt64, spverrors.ErrInvalidUintValue
		}
		return int64(v), nil
	case uint64:
		if v > uint64(math.MaxInt64) {
			return math.MaxInt64, spverrors.ErrInvalidUint64Value
		}
		return int64(v), nil
	case uint32:
		return int64(v), nil
	case uint16:
		return int64(v), nil
	case uint8:
		return int64(v), nil
	case int64:
		return v, nil
	default:
		return 0, spverrors.ErrUnsupportedDestinationType
	}
}

// ConvertToInt is a generic function that converts numeric values to int
func ConvertToInt[T any](value T) (int, error) {
	maxInt := int(math.MaxInt) // Max int value for the current platform (32-bit or 64-bit)
	minInt := int(math.MinInt) // Min int value for the current platform

	switch v := any(value).(type) {
	case int:
		return v, nil
	case int32:
		return int(v), nil
	case int16:
		return int(v), nil
	case int8:
		return int(v), nil
	case uint:
		if v > uint(maxInt) {
			return maxInt, spverrors.ErrInvalidUintValue
		}
		return int(v), nil
	case uint64:
		if v > uint64(maxInt) {
			return maxInt, spverrors.ErrInvalidUint64Value
		}
		return int(v), nil
	case uint32:
		if v > uint32(maxInt) {
			return maxInt, spverrors.ErrInvalidUintValue
		}
		return int(v), nil
	case uint16:
		return int(v), nil
	case uint8:
		return int(v), nil
	case int64:
		if v > int64(maxInt) || v < int64(minInt) {
			return 0, spverrors.ErrInvalidIntValue
		}
		return int(v), nil
	default:
		return 0, spverrors.ErrUnsupportedTypeForConversion
	}
}

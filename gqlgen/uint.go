package gqlgen

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"

	"github.com/99designs/gqlgen/graphql"
)

// MarshalUint uint support in gqlgen
func MarshalUint(i uint) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = io.WriteString(w, strconv.FormatUint(uint64(i), 10))
	})
}

// UnmarshalUint uint support in gqlgen
func UnmarshalUint(v interface{}) (uint, error) {
	switch v := v.(type) {
	case string:
		u64, err := strconv.ParseUint(v, 10, 64)
		return uint(u64), err
	case int:
		return uint(v), nil
	case int64:
		return uint(v), nil
	case json.Number:
		u64, err := strconv.ParseUint(string(v), 10, 64)
		if err != nil {
			return 0, err
		}
		if u64 <= math.MaxUint {
			return uint(u64), err
		}
		return 0, errors.New("value is > math.MaxUint")
	default:
		return 0, fmt.Errorf("%T is not an uint", v)
	}
}

// MarshalUint64 uint64 support in gqlgen
func MarshalUint64(i uint64) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = io.WriteString(w, strconv.FormatUint(i, 10))
	})
}

// UnmarshalUint64 uint64 support in gqlgen
func UnmarshalUint64(v interface{}) (uint64, error) {
	switch v := v.(type) {
	case string:
		return strconv.ParseUint(v, 10, 64)
	case int:
		return uint64(v), nil
	case int64:
		return uint64(v), nil
	case json.Number:
		return strconv.ParseUint(string(v), 10, 64)
	default:
		return 0, fmt.Errorf("%T is not an uint", v)
	}
}

// MarshalUint32 uint32 support in gqlgen
func MarshalUint32(i uint32) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = io.WriteString(w, strconv.FormatUint(uint64(i), 10))
	})
}

// UnmarshalUint32 uint32 support in gqlgen
func UnmarshalUint32(v interface{}) (uint32, error) {
	switch v := v.(type) {
	case string:
		iv, err := strconv.ParseInt(v, 10, 32)
		if err != nil {
			return 0, err
		}
		return uint32(iv), nil
	case int:
		return uint32(v), nil
	case int64:
		return uint32(v), nil
	case json.Number:
		iv, err := strconv.ParseUint(string(v), 10, 32)
		if err != nil {
			return 0, err
		}
		return uint32(iv), nil
	default:
		return 0, fmt.Errorf("%T is not an uint", v)
	}
}

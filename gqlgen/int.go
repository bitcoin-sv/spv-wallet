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

// MarshalInt16 int16 support in gqlgen
func MarshalInt16(i int16) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = io.WriteString(w, strconv.FormatInt(int64(i), 10))
	})
}

// UnmarshalInt16 int16 support in gqlgen
func UnmarshalInt16(v interface{}) (int16, error) {
	switch v := v.(type) {
	case string:
		u64, err := strconv.ParseUint(v, 10, 16)
		return int16(u64), err
	case int:
		return int16(v), nil
	case int32:
		return int16(v), nil
	case int64:
		return int16(v), nil
	case json.Number:
		u64, err := strconv.ParseUint(string(v), 10, 16)
		if err != nil {
			return 0, err
		}
		if u64 <= math.MaxUint16 {
			return int16(u64), err
		}
		return 0, errors.New("value is > math.MaxUint16")
	default:
		return 0, fmt.Errorf("%T is not an int16", v)
	}
}

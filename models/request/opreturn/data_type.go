package opreturn

import (
	"encoding/json"

	"github.com/bitcoin-sv/spv-wallet/models/request/internal"
)

// DataType represents the type of data in the OP_RETURN output.
type DataType int

// Enum values for DataType
const (
	DataTypeDefault DataType = iota
	DataTypeStrings
	DataTypeHexes
)

// UnmarshalJSON custom unmarshaler for DataTypeEnum
func (d *DataType) UnmarshalJSON(data []byte) error {
	var dataType string
	if err := json.Unmarshal(data, &dataType); err != nil {
		return internal.ErrorUnmarshal.Wrap(err)
	}

	switch dataType {
	case "strings":
		*d = DataTypeStrings
	case "hexes":
		*d = DataTypeHexes
	default:
		return internal.ErrorInvalidDataType
	}
	return nil
}

// MarshalJSON custom marshaler for DataType Enum
func (d DataType) MarshalJSON() ([]byte, error) {
	var dataType string
	switch d {
	case DataTypeDefault:
		dataType = ""
	case DataTypeStrings:
		dataType = "strings"
	case DataTypeHexes:
		dataType = "hexes"
	default:
		return nil, internal.ErrorInvalidDataType
	}
	data, err := json.Marshal(dataType)
	if err != nil {
		return nil, internal.ErrorMarshal.Wrap(err)
	}
	return data, nil
}

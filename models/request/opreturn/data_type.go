package opreturn

import (
	"encoding/json"
	"errors"
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
		return err //nolint:wrapcheck // UnmarshalJSON is run internally by json.Unmarshal on "DataType" object, so we don't want to wrap the error
	}

	switch dataType {
	case "strings":
		*d = DataTypeStrings
	case "hexes":
		*d = DataTypeHexes
	default:
		return errors.New("invalid data type")
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
		return nil, errors.New("invalid data type")
	}
	return json.Marshal(dataType) //nolint:wrapcheck // MarshalJSON is run internally by json.Marshal on "DataType" object, so we don't want to wrap the error
}

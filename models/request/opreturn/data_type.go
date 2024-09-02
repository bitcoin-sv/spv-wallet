package opreturn

import (
	"encoding/json"
	"errors"
)

type DataType int

const (
	DataTypeDefault DataType = iota
	DataTypeStrings
	DataTypeHexes
)

// UnmarshalJSON custom unmarshaler for DataTypeEnum
func (d *DataType) UnmarshalJSON(data []byte) error {
	var dataType string
	if err := json.Unmarshal(data, &dataType); err != nil {
		return err
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
	return json.Marshal(dataType)
}

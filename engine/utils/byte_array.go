package utils

import "fmt"

// ToByteArray converts string or []byte to byte array or returns an error
func ToByteArray(value interface{}) ([]byte, error) {
	switch typedValue := value.(type) {
	case []byte:
		return typedValue, nil
	case string:
		return []byte(typedValue), nil
	default:
		return nil, fmt.Errorf("unsupported type: %T", value)
	}
}

// StrOrBytesToString converts string or []byte to string or returns an error
func StrOrBytesToString(value interface{}) (string, error) {
	switch typedValue := value.(type) {
	case []byte:
		return string(typedValue), nil
	case string:
		return typedValue, nil
	default:
		return "", fmt.Errorf("unsupported type: %T", value)
	}
}

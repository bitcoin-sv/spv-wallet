package request

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"

	"github.com/bitcoin-sv/spv-wallet/models/request/opreturn"
)

func MakeOutputParser[T Output]() OutputParser {
	return OutputParser{
		Unmarshal: func(raw json.RawMessage) (Output, error) { return unmarshall[T](raw) },
	}
}

var supportedOutputTypes = map[string]OutputParser{
	"op_return": MakeOutputParser[opreturn.Output](),
}

// IsOutputSupported checks by type name if output is supported
func IsOutputSupported(typeName string) bool {
	_, ok := supportedOutputTypes[typeName]
	return ok
}

func unmarshall[T Output](raw json.RawMessage) (Output, error) {
	var desiredType T
	if err := json.Unmarshal(raw, &desiredType); err != nil {
		return nil, err //nolint:wrapcheck // TODO it later
	}
	return desiredType, nil
}

// OutputParser defines supported outputs for json unmarshall
type OutputParser struct {
	Unmarshal func(json.RawMessage) (Output, error)
}

// unmarshalOutput used by DraftTransaction unmarshalling to get Output object by type
// IMPORTANT: Every time a new output type is added, it must be handled here also.
func unmarshalOutput(rawOutput json.RawMessage, outputType string) (Output, error) {
	parser, ok := supportedOutputTypes[outputType]
	if !ok {
		return nil, errors.New("unsupported output type")
	}
	return parser.Unmarshal(rawOutput)
}

func expandOutputForMarshaling(output Output) (map[string]any, error) {
	if !IsOutputSupported(output.GetType()) {
		return nil, errors.New("unsupported output type")
	}
	result := map[string]any{
		"type": output.GetType(),
	}

	v := reflect.ValueOf(output)
	t := reflect.TypeOf(output)

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		fieldName, omit := keyNameFromJSONTag(field, value)
		if omit {
			continue
		}

		result[fieldName] = value.Interface()
	}

	return result, nil
}

func isZeroOfUnderlyingType(rValue reflect.Value) bool {
	return rValue.Interface() == reflect.Zero(rValue.Type()).Interface()
}

func keyNameFromJSONTag(field reflect.StructField, value reflect.Value) (name string, omit bool) {
	tag := field.Tag.Get("json")
	tagParts := strings.Split(tag, ",")

	name = tagParts[0]
	omitEmpty := len(tagParts) > 1 && tagParts[1] == "omitempty"

	if name == "-" || name == "" {
		name = field.Name
	}

	omit = omitEmpty && isZeroOfUnderlyingType(value)
	return
}

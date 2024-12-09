package request

import (
	"encoding/json"
)

// TransactionSpecification represents a request with specification for making a transaction outline.
type TransactionSpecification struct {
	Outputs []Output `json:"-"`
}

// Output represents an output in a transaction outline request.
type Output interface {
	GetType() string
}

// UnmarshalJSON custom unmarshall logic for TransactionSpecification
func (dt *TransactionSpecification) UnmarshalJSON(data []byte) error {
	rawOutputs, err := dt.unmarshalPartials(data)
	if err != nil {
		return err
	}

	// Unmarshal each output based on the type field.
	outputs, err := unmarshalOutputs(rawOutputs)
	if err != nil {
		return err
	}
	dt.Outputs = outputs

	return nil
}

// unmarshalPartials unmarshalls the data into the TransactionSpecification
// and returns also raw parts that couldn't be unmarshalled out of the box.
func (dt *TransactionSpecification) unmarshalPartials(data []byte) (rawOutputs []json.RawMessage, err error) {
	// Define a temporary struct to unmarshal the struct without unmarshalling outputs.
	// We're defining it here, to not publish Alias type.
	type Alias TransactionSpecification
	temp := &struct {
		Outputs []json.RawMessage `json:"outputs"`
		*Alias
	}{
		Alias: (*Alias)(dt),
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return nil, err //nolint:wrapcheck // unmarshalPartials is run internally by json.Unmarshal, so we don't want to wrap the error
	}

	return temp.Outputs, nil
}

func unmarshalOutputs(outputs []json.RawMessage) ([]Output, error) {
	result := make([]Output, len(outputs))
	for i, rawOutput := range outputs {
		var typeField struct {
			Type string `json:"type"`
		}
		if err := json.Unmarshal(rawOutput, &typeField); err != nil {
			return nil, err //nolint:wrapcheck // unmarshalOutputs is run internally by json.Unmarshal, so we don't want to wrap the error
		}

		output, err := unmarshalOutput(rawOutput, typeField.Type)
		if err != nil {
			return nil, err
		}
		result[i] = output
	}
	return result, nil
}

// MarshalJSON custom marshaller for TransactionSpecification
func (dt *TransactionSpecification) MarshalJSON() ([]byte, error) {
	type Alias TransactionSpecification
	temp := &struct {
		Outputs []any `json:"outputs"`
		*Alias
	}{
		Alias: (*Alias)(dt),
	}

	for _, output := range dt.Outputs {
		out, err := expandOutputForMarshaling(output)
		if err != nil {
			return nil, err
		}
		temp.Outputs = append(temp.Outputs, out)
	}

	return json.Marshal(temp) //nolint:wrapcheck // MarshalJSON is run internally by json.Marshal, so we don't want to wrap the error
}

package request

import (
	"encoding/json"
)

// DraftTransaction represents a request with specification for making a draft transaction.
type DraftTransaction struct {
	Outputs []Output `json:"-"`
}

// Output represents an output in a draft transaction request.
type Output interface {
	GetType() string
}

// UnmarshalJSON custom unmarshall logic for DraftTransaction
func (dt *DraftTransaction) UnmarshalJSON(data []byte) error {
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

// unmarshalPartials unmarshalls the data into the DraftTransaction
// and returns also raw parts that couldn't be unmarshalled out of the box.
func (dt *DraftTransaction) unmarshalPartials(data []byte) (rawOutputs []json.RawMessage, err error) {
	// Define a temporary struct to unmarshal the struct without unmarshalling outputs.
	// We're defining it here, to not publish Alias type.
	type Alias DraftTransaction
	temp := &struct {
		Outputs []json.RawMessage `json:"outputs"`
		*Alias
	}{
		Alias: (*Alias)(dt),
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return nil, err
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
			return nil, err
		}

		output, err := unmarshalOutput(rawOutput, typeField.Type)
		if err != nil {
			return nil, err
		}
		result[i] = output
	}
	return result, nil
}

// MarshalJSON custom marshaller for DraftTransaction
func (dt *DraftTransaction) MarshalJSON() ([]byte, error) {
	type Alias DraftTransaction
	temp := &struct {
		Outputs []interface{} `json:"outputs"`
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

	return json.Marshal(temp)
}

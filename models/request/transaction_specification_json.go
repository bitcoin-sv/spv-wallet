package request

import (
	"encoding/json"
	"errors"

	"github.com/bitcoin-sv/spv-wallet/models/request/opreturn"
	paymailreq "github.com/bitcoin-sv/spv-wallet/models/request/paymail"
)

// unmarshalOutput used by TransactionSpecification unmarshalling to get Output object by type
// IMPORTANT: Every time a new output type is added, it must be handled here also.
func unmarshalOutput(rawOutput json.RawMessage, outputType string) (Output, error) {
	switch outputType {
	case "op_return":
		var out opreturn.Output
		if err := json.Unmarshal(rawOutput, &out); err != nil {
			return nil, err
		}
		return out, nil
	case "paymail":
		var out paymailreq.Output
		if err := json.Unmarshal(rawOutput, &out); err != nil {
			return nil, err
		}
		return out, nil
	default:
		return nil, errors.New("unsupported output type")
	}
}

// expandOutputForMarshaling used by TransactionSpecification marshalling to expand Output object before marshalling.
// IMPORTANT: Every time a new output type is added, it must be handled here also.
func expandOutputForMarshaling(output Output) (any, error) {
	switch o := output.(type) {
	// unfortunately we must do the same for each and every type,
	// because go json is not handling unwrapping embedded type when using just Output interface.
	case opreturn.Output:
		return struct {
			Type string `json:"type"`
			*opreturn.Output
		}{
			Type:   o.GetType(),
			Output: &o,
		}, nil
	case paymailreq.Output:
		return struct {
			Type string `json:"type"`
			*paymailreq.Output
		}{
			Type:   o.GetType(),
			Output: &o,
		}, nil
	default:
		return nil, errors.New("unsupported output type")
	}
}

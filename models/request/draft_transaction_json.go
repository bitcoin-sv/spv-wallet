package request

import (
	"encoding/json"
	"errors"

	"github.com/bitcoin-sv/spv-wallet/models/request/opreturn"
)

// unmarshalOutput used by DraftTransaction unmarshalling to get Output object by type
// IMPORTANT: Every time a new output type is added, it must be handled here also.
func unmarshalOutput(rawOutput json.RawMessage, outputType string) (Output, error) {
	switch outputType {
	case "op_return":
		var opReturnOutput opreturn.Output
		if err := json.Unmarshal(rawOutput, &opReturnOutput); err != nil {
			return nil, err
		}
		return opReturnOutput, nil
	default:
		return nil, errors.New("unsupported output type")
	}
}

// expandOutputForMarshaling used by DraftTransaction marshalling to expand Output object before marshalling.
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
	default:
		return nil, errors.New("unsupported output type")
	}
}

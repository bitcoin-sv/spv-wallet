package internal

import "github.com/bitcoin-sv/spv-wallet/models"

var (
	// ErrorUnmarshal is an error for unmarshalling model.
	ErrorUnmarshal = models.SPVError{Code: "error-model-unmarshal", StatusCode: 400, Message: "Error unmarshalling model"}

	// ErrorMarshal is an error for marshaling model.
	ErrorMarshal = models.SPVError{Code: "error-model-marshal", StatusCode: 500, Message: "Error marshaling model"}

	// ErrorUnsupportedOutputType is an error for unsupported output type.
	ErrorUnsupportedOutputType = models.SPVError{Code: "error-unsupported-output-type", StatusCode: 400, Message: "Unsupported output type"}

	// ErrorInvalidDataType is when unsupported data type is provided
	ErrorInvalidDataType = models.SPVError{Code: "error-invalid-data-type", StatusCode: 400, Message: "Invalid data type"}
)

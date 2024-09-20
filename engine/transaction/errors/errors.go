package txerrors

import "github.com/bitcoin-sv/spv-wallet/models"

var (
	// ErrDraftSpecificationRequired is returned when a draft is created with no specification.
	ErrDraftSpecificationRequired = models.SPVError{Code: "draft-spec-required", Message: "draft requires a specification", StatusCode: 400}

	// ErrDraftRequiresAtLeastOneOutput is returned when a draft is created with no outputs.
	ErrDraftRequiresAtLeastOneOutput = models.SPVError{Code: "draft-output-required", Message: "draft requires at least one output", StatusCode: 400}

	// ErrDraftOpReturnDataRequired is returned when an OP_RETURN output is created with no data.
	ErrDraftOpReturnDataRequired = models.SPVError{Code: "draft-op-return-data-required", Message: "data is required for OP_RETURN output", StatusCode: 400}

	// ErrDraftOpReturnDataTooLarge is returned when OP_RETURN data part is too big to add to transaction.
	ErrDraftOpReturnDataTooLarge = models.SPVError{Code: "draft-op-return-data-too-large", Message: "OP_RETURN data is too large", StatusCode: 400}

	// ErrOutputValueTooLow is when the satoshis output is too low for a given type of output.
	ErrOutputValueTooLow = models.SPVError{Message: "output value is too low", StatusCode: 400, Code: "error-transaction-output-value-too-low"}
)

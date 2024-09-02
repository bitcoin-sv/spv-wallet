package txerrors

import "github.com/bitcoin-sv/spv-wallet/models"

var (
	// ErrDraftSpecificationRequired is returned when a draft is created with no specification.
	ErrDraftSpecificationRequired = models.SPVError{Code: "draft-spec-required", Message: "draft requires a specification", StatusCode: 400}

	// ErrDraftRequiresAtLeastOneOutput is returned when a draft is created with no outputs.
	ErrDraftRequiresAtLeastOneOutput = models.SPVError{Code: "draft-output-required", Message: "draft requires at least one output", StatusCode: 400}

	// ErrDraftOpReturnDataRequired is returned when an OP_RETURN output is created with no data.
	ErrDraftOpReturnDataRequired = models.SPVError{Code: "draft-op-return-data-required", Message: "data is required for OP_RETURN output", StatusCode: 400}
)

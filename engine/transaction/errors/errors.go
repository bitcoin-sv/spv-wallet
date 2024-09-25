package txerrors

import "github.com/bitcoin-sv/spv-wallet/models"

var (
	// ErrDraftSpecificationRequired is returned when a draft is created with no specification.
	ErrDraftSpecificationRequired = models.SPVError{Code: "draft-spec-required", Message: "draft requires a specification", StatusCode: 400}

	// ErrDraftSpecificationXPubIDRequired is returned when a draft is created without xPubID.
	ErrDraftSpecificationXPubIDRequired = models.SPVError{Code: "draft-spec-xpub-id-required", Message: "cannot create draft without knowledge about xPubID", StatusCode: 500}

	// ErrDraftRequiresAtLeastOneOutput is returned when a draft is created with no outputs.
	ErrDraftRequiresAtLeastOneOutput = models.SPVError{Code: "draft-output-required", Message: "draft requires at least one output", StatusCode: 400}

	// ErrDraftOpReturnDataRequired is returned when an OP_RETURN output is created with no data.
	ErrDraftOpReturnDataRequired = models.SPVError{Code: "draft-op-return-data-required", Message: "data is required for OP_RETURN output", StatusCode: 400}

	// ErrDraftOpReturnDataTooLarge is returned when OP_RETURN data part is too big to add to transaction.
	ErrDraftOpReturnDataTooLarge = models.SPVError{Code: "draft-op-return-data-too-large", Message: "OP_RETURN data is too large", StatusCode: 400}

	// ErrDraftOpReturnUnsupportedDataType is returned when the data type for an OP_RETURN output is unsupported.
	ErrDraftOpReturnUnsupportedDataType = models.SPVError{Code: "draft-op-return-data-type-unsupported", Message: "unsupported data type for OP_RETURN output", StatusCode: 400}

	// ErrDraftSenderPaymailAddressNoDefault is when it is not possible to determine the default address for the sender.
	ErrDraftSenderPaymailAddressNoDefault = models.SPVError{Message: "cannot choose paymail address of the sender", StatusCode: 400, Code: "error-draft-paymail-address-no-default"}

	// ErrFailedToDecodeHex is returned when hex decoding fails.
	ErrFailedToDecodeHex = models.SPVError{Code: "failed-to-decode-hex", Message: "failed to decode hex", StatusCode: 400}

	// ErrReceiverPaymailAddressIsInvalid is when the receiver paymail address is NOT alias@domain.com
	ErrReceiverPaymailAddressIsInvalid = models.SPVError{Message: "receiver paymail address is invalid", StatusCode: 400, Code: "error-paymail-address-invalid-receiver"}

	// ErrSenderPaymailAddressIsInvalid is when the sender paymail address is NOT alias@domain.com
	ErrSenderPaymailAddressIsInvalid = models.SPVError{Message: "sender paymail address is invalid", StatusCode: 400, Code: "error-paymail-address-invalid-sender"}

	// ErrOutputValueTooLow is when the satoshis output is too low for a given type of output.
	ErrOutputValueTooLow = models.SPVError{Message: "output value is too low", StatusCode: 400, Code: "error-transaction-output-value-too-low"}
)

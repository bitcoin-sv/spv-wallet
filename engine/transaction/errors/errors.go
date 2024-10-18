package txerrors

import "github.com/bitcoin-sv/spv-wallet/models"

var (
	// ErrTxOutlineSpecificationRequired is returned when a transaction outline is created with no specification.
	ErrTxOutlineSpecificationRequired = models.SPVError{Code: "tx-spec-spec-required", Message: "transaction outline requires a specification", StatusCode: 400}

	// ErrTxOutlineSpecificationXPubIDRequired is returned when a transaction outline is created without xPubID.
	ErrTxOutlineSpecificationXPubIDRequired = models.SPVError{Code: "tx-spec-spec-xpub-id-required", Message: "cannot create transaction outline without knowledge about xPubID", StatusCode: 500}

	// ErrTxOutlineRequiresAtLeastOneOutput is returned when a transaction outline is created with no outputs.
	ErrTxOutlineRequiresAtLeastOneOutput = models.SPVError{Code: "tx-spec-output-required", Message: "transaction outline requires at least one output", StatusCode: 400}

	// ErrTxOutlineOpReturnDataRequired is returned when an OP_RETURN output is created with no data.
	ErrTxOutlineOpReturnDataRequired = models.SPVError{Code: "tx-spec-op-return-data-required", Message: "data is required for OP_RETURN output", StatusCode: 400}

	// ErrTxOutlineOpReturnDataTooLarge is returned when OP_RETURN data part is too big to add to transaction.
	ErrTxOutlineOpReturnDataTooLarge = models.SPVError{Code: "tx-spec-op-return-data-too-large", Message: "OP_RETURN data is too large", StatusCode: 400}

	// ErrTxOutlineOpReturnUnsupportedDataType is returned when the data type for an OP_RETURN output is unsupported.
	ErrTxOutlineOpReturnUnsupportedDataType = models.SPVError{Code: "tx-spec-op-return-data-type-unsupported", Message: "unsupported data type for OP_RETURN output", StatusCode: 400}

	// ErrTxOutlineSenderPaymailAddressNoDefault is when it is not possible to determine the default address for the sender.
	ErrTxOutlineSenderPaymailAddressNoDefault = models.SPVError{Code: "error-tx-spec-paymail-address-no-default", Message: "cannot choose paymail address of the sender", StatusCode: 400}

	// ErrFailedToDecodeHex is returned when hex decoding fails.
	ErrFailedToDecodeHex = models.SPVError{Code: "failed-to-decode-hex", Message: "failed to decode hex", StatusCode: 400}

	// ErrReceiverPaymailAddressIsInvalid is when the receiver paymail address is NOT alias@domain.com
	ErrReceiverPaymailAddressIsInvalid = models.SPVError{Code: "error-paymail-address-invalid-receiver", Message: "receiver paymail address is invalid", StatusCode: 400}

	// ErrSenderPaymailAddressIsInvalid is when the sender paymail address is NOT alias@domain.com
	ErrSenderPaymailAddressIsInvalid = models.SPVError{Code: "error-paymail-address-invalid-sender", Message: "sender paymail address is invalid", StatusCode: 400}

	// ErrOutputValueTooLow is when the satoshis output is too low for a given type of output.
	ErrOutputValueTooLow = models.SPVError{Code: "error-transaction-output-value-too-low", Message: "output value is too low", StatusCode: 400}
)

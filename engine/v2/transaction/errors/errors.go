package txerrors

import "github.com/bitcoin-sv/spv-wallet/models"

var (
	// ErrTxOutlineSpecificationRequired is returned when a transaction outline is created with no specification.
	ErrTxOutlineSpecificationRequired = models.SPVError{Code: "tx-spec-spec-required", Message: "transaction outline requires a specification", StatusCode: 400}

	// ErrTxOutlineSpecificationUserIDRequired is returned when a transaction outline is created without UserID.
	ErrTxOutlineSpecificationUserIDRequired = models.SPVError{Code: "tx-spec-spec-user-id-required", Message: "cannot create transaction outline without knowledge about userID", StatusCode: 500}

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

	// ErrTxOutlineInsufficientFunds is returned when user has not enough BSV in UTXOs to fund the transaction.
	ErrTxOutlineInsufficientFunds = models.SPVError{Code: "tx-outline-not-enough-funds", Message: "not enough funds to make the transaction", StatusCode: 422}

	// ErrFailedToDecodeHex is returned when hex decoding fails.
	ErrFailedToDecodeHex = models.SPVError{Code: "failed-to-decode-hex", Message: "failed to decode hex", StatusCode: 400}

	// ErrReceiverPaymailAddressIsInvalid is when the receiver paymail address is NOT alias@domain.com
	ErrReceiverPaymailAddressIsInvalid = models.SPVError{Code: "error-paymail-address-invalid-receiver", Message: "receiver paymail address is invalid", StatusCode: 400}

	// ErrSenderPaymailAddressIsInvalid is when the sender paymail address is NOT alias@domain.com
	ErrSenderPaymailAddressIsInvalid = models.SPVError{Code: "error-paymail-address-invalid-sender", Message: "sender paymail address is invalid", StatusCode: 400}

	// ErrOutputValueTooLow is when the satoshis output is too low for a given type of output.
	ErrOutputValueTooLow = models.SPVError{Code: "error-transaction-output-value-too-low", Message: "output value is too low", StatusCode: 400}

	// ErrTxValidation is when the transaction validation fails.
	ErrTxValidation = models.SPVError{Code: "error-transaction-validation", Message: "transaction validation failed", StatusCode: 400}

	// ErrUTXOSpent is when the UTXO is already spent.
	ErrUTXOSpent = models.SPVError{Code: "error-utxo-spent", Message: "UTXO is already spent", StatusCode: 400}

	// ErrParsingScript is when the script parsing fails.
	ErrParsingScript = models.SPVError{Code: "error-parsing-script", Message: "failed to parse script", StatusCode: 400}

	// ErrSavingData is when the data saving fails.
	ErrSavingData = models.SPVError{Code: "error-saving-data", Message: "failed to save data", StatusCode: 400}

	// ErrTxBroadcast is when the transaction broadcast fails.
	ErrTxBroadcast = models.SPVError{Code: "error-tx-broadcast", Message: "failed to broadcast transaction", StatusCode: 500}

	// ErrAnnotationIndexOutOfRange is when the annotation index is out of range.
	ErrAnnotationIndexOutOfRange = models.SPVError{Code: "error-annotation-index-out-of-range", Message: "annotation index is out of range", StatusCode: 400}

	// ErrGettingOutputs is when getting outputs fails.
	ErrGettingOutputs = models.SPVError{Code: "error-getting-outputs", Message: "failed to get outputs", StatusCode: 500}

	// ErrAnnotationMismatch is when the annotation does not match to actual output content.
	ErrAnnotationMismatch = models.SPVError{Code: "error-annotation-mismatch", Message: "annotation mismatch", StatusCode: 400}

	// ErrAnnotationIndexConversion is when the annotation index conversion fails.
	ErrAnnotationIndexConversion = models.SPVError{Code: "error-annotation-index-conversion", Message: "failed to convert annotation index", StatusCode: 400}

	// ErrOnlyPushDataAllowed is when only PUSHDATA operations are allowed in OP_RETURN script.
	ErrOnlyPushDataAllowed = models.SPVError{Code: "error-only-push-data-allowed", Message: "Only PUSHDATA operations are allowed in OP_RETURN script", StatusCode: 400}

	// ErrUnexpectedErrorDuringInputsSelection is when an unexpected error occurs during inputs selection for transaction outline.
	ErrUnexpectedErrorDuringInputsSelection = models.SPVError{Code: "error-input-selection", Message: "unexpected error during inputs selection", StatusCode: 500}

	// ErrNoOperations is when there are no operations to save.
	ErrNoOperations = models.SPVError{Code: "error-no-operations", Message: "no operations to save", StatusCode: 400}

	// ErrGettingAddresses is when getting addresses fails.
	ErrGettingAddresses = models.SPVError{Code: "error-getting-addresses", Message: "failed to get addresses", StatusCode: 500}

	// ErrMultiPaymailRecipientsNotSupported is when the transaction has multiple paymail recipients.
	ErrMultiPaymailRecipientsNotSupported = models.SPVError{Code: "error-multi-paymail-recipients", Message: "paymail transaction with multiple recipients is not supported", StatusCode: 400}
)

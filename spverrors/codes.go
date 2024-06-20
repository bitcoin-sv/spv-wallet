package spverrors

const (
	AuthorizationError               = "authorization-error"
	BindingError                     = "binding-error"
	ContactIncorrectStatusError      = "incorrect-status-error"
	ContactRequesterInvalidXpubError = "requester-invalid-xpub-error"
	ContactMoreThanOnePaymailError   = "more-than-one-paymail-registered-error"
	CreateError                      = "add-new-model-error"
	MissingFieldError                = "missing-field-error"
	NotFoundError                    = "not-found-error"
	RecordTransactionError           = "record-transaction-error"
	SaveError                        = "save-error"
	AccessKeyValidationError         = "access-key-validation-error"
	DestinationValidationError       = "destination-validation-error"
	PaymailValidationError           = "paymail-validation-error"
	XpubValidationError              = "xpub-validation-error"
	TransactionValidationError       = "transaction-validation-error"
	UtxoValidationError              = "utxo-validation-error"
	UnknownError                     = "unknown-error"
)

package spverrors

type SPVError struct {
	Code       string
	Message    string
	StatusCode int
}

var SPVErrorResponses = map[error]SPVError{
	// ACCESS KEY ERRORS
	ErrCouldNotFindAccessKey: {
		Code:       NotFoundError,
		Message:    ErrCouldNotFindAccessKey.Error(),
		StatusCode: 404,
	},
	ErrAccessKeyRevoked: {
		Code:       AccessKeyValidationError,
		Message:    ErrAccessKeyRevoked.Error(),
		StatusCode: 400,
	},
	// AUTHORIZATION ERRORS
	ErrAuthorization: {
		Code:       AuthorizationError,
		Message:    ErrAuthorization.Error(),
		StatusCode: 401,
	},
	ErrInvalidFilterOption: {
		Code:       BindingError,
		Message:    ErrInvalidFilterOption.Error(),
		StatusCode: 400,
	},
	// BINDING ERRORS
	ErrCannotBindRequest: {
		Code:       BindingError,
		Message:    ErrCannotBindRequest.Error(),
		StatusCode: 400,
	},
	ErrInvalidConditions: {
		Code:       BindingError,
		Message:    ErrInvalidConditions.Error(),
		StatusCode: 400,
	},
	// DESTINATION ERRORS
	ErrCouldNotFindDestination: {
		Code:       NotFoundError,
		Message:    ErrCouldNotFindDestination.Error(),
		StatusCode: 404,
	},
	ErrUnsupportedDestinationType: {
		Code:       DestinationValidationError,
		Message:    ErrUnsupportedDestinationType.Error(),
		StatusCode: 400,
	},
	ErrUnknownLockingScript: {
		Code:       DestinationValidationError,
		Message:    ErrUnknownLockingScript.Error(),
		StatusCode: 400,
	},
	// CONTACT ERRORS
	ErrContactNotFound: {
		Code:       NotFoundError,
		Message:    ErrContactNotFound.Error(),
		StatusCode: 404,
	},
	ErrContactIncorrectStatus: {
		Code:       ContactIncorrectStatusError,
		Message:    ErrContactIncorrectStatus.Error(),
		StatusCode: 422,
	},
	ErrInvalidRequesterXpub: {
		Code:       ContactRequesterInvalidXpubError,
		Message:    ErrInvalidRequesterXpub.Error(),
		StatusCode: 400,
	},
	ErrMoreThanOnePaymailRegistered: {
		Code:       ContactMoreThanOnePaymailError,
		Message:    ErrMoreThanOnePaymailRegistered.Error(),
		StatusCode: 400,
	},
	ErrAddingContactRequest: {
		Code:       CreateError,
		Message:    ErrAddingContactRequest.Error(),
		StatusCode: 400,
	},
	ErrMissingContactID: {
		Code:       MissingFieldError,
		Message:    ErrMissingContactID.Error(),
		StatusCode: 400,
	},
	ErrMissingContactPaymail: {
		Code:       MissingFieldError,
		Message:    ErrMissingContactPaymail.Error(),
		StatusCode: 400,
	},
	ErrMissingContactXPubKey: {
		Code:       MissingFieldError,
		Message:    ErrMissingContactXPubKey.Error(),
		StatusCode: 400,
	},
	ErrMissingContactStatus: {
		Code:       MissingFieldError,
		Message:    ErrMissingContactStatus.Error(),
		StatusCode: 400,
	},
	ErrMissingContactOwnerXPubId: {
		Code:       MissingFieldError,
		Message:    ErrMissingContactOwnerXPubId.Error(),
		StatusCode: 400,
	},
	ErrEmptyContactPubKey: {
		Code:       MissingFieldError,
		Message:    ErrEmptyContactPubKey.Error(),
		StatusCode: 400,
	},
	ErrEmptyContactPaymail: {
		Code:       MissingFieldError,
		Message:    ErrEmptyContactPaymail.Error(),
		StatusCode: 400,
	},
	ErrMissingContactFullName: {
		Code:       MissingFieldError,
		Message:    ErrMissingContactFullName.Error(),
		StatusCode: 400,
	},
	// PAYMAIL ERRORS
	ErrCouldNotFindPaymail: {
		Code:       NotFoundError,
		Message:    ErrCouldNotFindPaymail.Error(),
		StatusCode: 404,
	},
	ErrPaymailAddressIsInvalid: {
		Code:       PaymailValidationError,
		Message:    ErrPaymailAddressIsInvalid.Error(),
		StatusCode: 400,
	},
	ErrMissingPaymailID: {
		Code:       PaymailValidationError,
		Message:    ErrMissingPaymailID.Error(),
		StatusCode: 400,
	},
	ErrMissingPaymailAddress: {
		Code:       PaymailValidationError,
		Message:    ErrMissingPaymailAddress.Error(),
		StatusCode: 400,
	},
	ErrMissingPaymailDomain: {
		Code:       PaymailValidationError,
		Message:    ErrMissingPaymailDomain.Error(),
		StatusCode: 400,
	},
	ErrMissingPaymailExternalXPub: {
		Code:       PaymailValidationError,
		Message:    ErrMissingPaymailExternalXPub.Error(),
		StatusCode: 400,
	},
	ErrMissingPaymailXPubID: {
		Code:       PaymailValidationError,
		Message:    ErrMissingPaymailXPubID.Error(),
		StatusCode: 400,
	},
	ErrPaymailAlreadyExists: {
		Code:       PaymailValidationError,
		Message:    ErrPaymailAlreadyExists.Error(),
		StatusCode: 400,
	},
	// TRANSACTION ERRORS
	ErrCouldNotFindTransaction: {
		Code:       NotFoundError,
		Message:    ErrCouldNotFindTransaction.Error(),
		StatusCode: 404,
	},
	ErrCouldNotFindSyncTx: {
		Code:       NotFoundError,
		Message:    ErrCouldNotFindSyncTx.Error(),
		StatusCode: 404,
	},
	ErrCouldNotFindDraftTx: {
		Code:       NotFoundError,
		Message:    ErrCouldNotFindDraftTx.Error(),
		StatusCode: 404,
	},
	ErrInvalidTransactionID: {
		Code:       TransactionValidationError,
		Message:    ErrInvalidTransactionID.Error(),
		StatusCode: 400,
	},
	ErrInvalidRequirements: {
		Code:       TransactionValidationError,
		Message:    ErrInvalidRequirements.Error(),
		StatusCode: 400,
	},
	ErrTransactionIDMismatch: {
		Code:       TransactionValidationError,
		Message:    ErrTransactionIDMismatch.Error(),
		StatusCode: 400,
	},
	ErrMissingTransactionOutputs: {
		Code:       TransactionValidationError,
		Message:    ErrMissingTransactionOutputs.Error(),
		StatusCode: 400,
	},
	ErrOutputValueTooLow: {
		Code:       TransactionValidationError,
		Message:    ErrOutputValueTooLow.Error(),
		StatusCode: 400,
	},
	ErrOutputValueTooHigh: {
		Code:       TransactionValidationError,
		Message:    ErrOutputValueTooHigh.Error(),
		StatusCode: 400,
	},
	ErrTransactionFeeInvalid: {
		Code:       TransactionValidationError,
		Message:    ErrTransactionFeeInvalid.Error(),
		StatusCode: 400,
	},
	ErrInvalidOpReturnOutput: {
		Code:       TransactionValidationError,
		Message:    ErrInvalidOpReturnOutput.Error(),
		StatusCode: 400,
	},
	ErrChangeStrategyNotImplemented: {
		Code:       TransactionValidationError,
		Message:    ErrChangeStrategyNotImplemented.Error(),
		StatusCode: 400,
	},
	ErrInvalidLockingScript: {
		Code:       TransactionValidationError,
		Message:    ErrInvalidLockingScript.Error(),
		StatusCode: 400,
	},
	ErrOutputValueNotRecognized: {
		Code:       TransactionValidationError,
		Message:    ErrOutputValueNotRecognized.Error(),
		StatusCode: 400,
	},
	ErrInvalidScriptOutput: {
		Code:       TransactionValidationError,
		Message:    ErrInvalidScriptOutput.Error(),
		StatusCode: 400,
	},
	ErrDraftIDMismatch: {
		Code:       TransactionValidationError,
		Message:    ErrDraftIDMismatch.Error(),
		StatusCode: 400,
	},
	ErrMissingTxHex: {
		Code:       TransactionValidationError,
		Message:    ErrMissingTxHex.Error(),
		StatusCode: 400,
	},
	ErrNoMatchingOutputs: {
		Code:       TransactionValidationError,
		Message:    ErrNoMatchingOutputs.Error(),
		StatusCode: 400,
	},
	ErrCreateOutgoingTxFailed: {
		Code:       RecordTransactionError,
		Message:    ErrCreateOutgoingTxFailed.Error(),
		StatusCode: 400,
	},
	ErrDuringSaveTx: {
		Code:       RecordTransactionError,
		Message:    ErrDuringSaveTx.Error(),
		StatusCode: 400,
	},
	ErrTransactionRejectedByP2PProvider: {
		Code:       RecordTransactionError,
		Message:    ErrTransactionRejectedByP2PProvider.Error(),
		StatusCode: 400,
	},
	ErrDraftTxHasNoOutputs: {
		Code:       RecordTransactionError,
		Message:    ErrDraftTxHasNoOutputs.Error(),
		StatusCode: 400,
	},
	ErrProcessP2PTx: {
		Code:       RecordTransactionError,
		Message:    ErrProcessP2PTx.Error(),
		StatusCode: 400,
	},
	// UTXO ERRORS
	ErrCouldNotFindUtxo: {
		Code:       NotFoundError,
		Message:    ErrCouldNotFindUtxo.Error(),
		StatusCode: 404,
	},
	ErrUtxoAlreadySpent: {
		Code:       UtxoValidationError,
		Message:    ErrUtxoAlreadySpent.Error(),
		StatusCode: 400,
	},
	ErrMissingUTXOsSpendable: {
		Code:       UtxoValidationError,
		Message:    ErrMissingUTXOsSpendable.Error(),
		StatusCode: 400,
	},
	ErrNotEnoughUtxos: {
		Code:       UtxoValidationError,
		Message:    ErrNotEnoughUtxos.Error(),
		StatusCode: 400,
	},
	ErrDuplicateUTXOs: {
		Code:       UtxoValidationError,
		Message:    ErrDuplicateUTXOs.Error(),
		StatusCode: 400,
	},
	ErrUtxoNotReserved: {
		Code:       UtxoValidationError,
		Message:    ErrUtxoNotReserved.Error(),
		StatusCode: 400,
	},
	// XPUB ERRORS
	ErrXpubInvalidLength: {
		Code:       XpubValidationError,
		Message:    ErrXpubInvalidLength.Error(),
		StatusCode: 400,
	},
	ErrXpubNoMatch: {
		Code:       XpubValidationError,
		Message:    ErrXpubNoMatch.Error(),
		StatusCode: 400,
	},
	ErrCouldNotFindXpub: {
		Code:       NotFoundError,
		Message:    ErrCouldNotFindXpub.Error(),
		StatusCode: 404,
	},
	ErrXpubIDMisMatch: {
		Code:       XpubValidationError,
		Message:    ErrXpubIDMisMatch.Error(),
		StatusCode: 404,
	},
	// MISSING FIELD ERROR
	ErrMissingFieldID: {
		Code:       MissingFieldError,
		Message:    ErrMissingFieldID.Error(),
		StatusCode: 400,
	},
	ErrMissingFieldXpubID: {
		Code:       MissingFieldError,
		Message:    ErrMissingFieldXpubID.Error(),
		StatusCode: 400,
	},
	ErrMissingFieldXpub: {
		Code:       MissingFieldError,
		Message:    ErrMissingFieldXpub.Error(),
		StatusCode: 400,
	},
	ErrMissingAccessKey: {
		Code:       MissingFieldError,
		Message:    ErrMissingAccessKey.Error(),
		StatusCode: 400,
	},
	ErrMissingAddress: {
		Code:       MissingFieldError,
		Message:    ErrMissingAddress.Error(),
		StatusCode: 400,
	},
	ErrOneOfTheFieldsIsRequired: {
		Code:       MissingFieldError,
		Message:    ErrOneOfTheFieldsIsRequired.Error(),
		StatusCode: 400,
	},
	ErrMissingFieldScriptPubKey: {
		Code:       MissingFieldError,
		Message:    ErrMissingFieldScriptPubKey.Error(),
		StatusCode: 400,
	},
	ErrMissingFieldSatoshis: {
		Code:       MissingFieldError,
		Message:    ErrMissingFieldSatoshis.Error(),
		StatusCode: 400,
	},
	ErrMissingFieldTransactionID: {
		Code:       MissingFieldError,
		Message:    ErrMissingFieldTransactionID.Error(),
		StatusCode: 400,
	},
	ErrMissingLockingScript: {
		Code:       MissingFieldError,
		Message:    ErrMissingLockingScript.Error(),
		StatusCode: 400,
	},
	// SAVE ERROR
	ErrMissingClient: {
		Code:       SaveError,
		Message:    ErrMissingClient.Error(),
		StatusCode: 400,
	},
	ErrDatastoreRequired: {
		Code:       SaveError,
		Message:    ErrDatastoreRequired.Error(),
		StatusCode: 400,
	},
}
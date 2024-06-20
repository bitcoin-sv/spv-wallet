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
	// AUTHORIZATION ERRORS
	ErrAuthorization: {
		Code:       AuthorizationError,
		Message:    ErrAuthorization.Error(),
		StatusCode: 401,
	},
	ErrInvalidFilterOption: {
		Code:       AuthorizationError,
		Message:    ErrInvalidFilterOption.Error(),
		StatusCode: 401,
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
	// PAYMAIL ERRORS
	ErrCouldNotFindPaymail: {
		Code:       NotFoundError,
		Message:    ErrCouldNotFindPaymail.Error(),
		StatusCode: 404,
	},
	// TRANSACTION ERRORS
	ErrCouldNotFindTransaction: {
		Code:       NotFoundError,
		Message:    ErrCouldNotFindTransaction.Error(),
		StatusCode: 404,
	},
	ErrInvalidTransactionID: {
		Code:       XpubValidationError,
		Message:    ErrInvalidTransactionID.Error(),
		StatusCode: 400,
	},
	ErrInvalidRequirements: {
		Code:       XpubValidationError,
		Message:    ErrInvalidRequirements.Error(),
		StatusCode: 400,
	},
	ErrTransactionIDMismatch: {
		Code:       XpubValidationError,
		Message:    ErrTransactionIDMismatch.Error(),
		StatusCode: 400,
	},
	// UTXO ERRORS
	ErrCouldNotFindUtxo: {
		Code:       NotFoundError,
		Message:    ErrCouldNotFindUtxo.Error(),
		StatusCode: 404,
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
	ErrMissingXpub: {
		Code:       MissingFieldError,
		Message:    ErrMissingXpub.Error(),
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

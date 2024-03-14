package engine

import (
	"errors"
)

// ErrCannotConvertToIDs is the error when the conversion fails from interface into type IDs
var ErrCannotConvertToIDs = errors.New("cannot convert value to type IDs")

// ErrMissingDestination is an error when a destination could not be found
var ErrMissingDestination = errors.New("destination could not be found")

// ErrMissingTransaction is an error when a transaction could not be found
var ErrMissingTransaction = errors.New("transaction could not be found")

// ErrMissingUtxo is an error when a given utxo could not be found
var ErrMissingUtxo = errors.New("utxo could not be found")

// ErrMissingFieldID is an error when missing the id field
var ErrMissingFieldID = errors.New("missing required field: id")

// ErrMissingFieldHex is an error when missing the hex field of a transaction
var ErrMissingFieldHex = errors.New("missing required field: hex")

// ErrMissingFieldHash is an error when missing the hex field of a transaction
var ErrMissingFieldHash = errors.New("missing required field: hash")

// ErrMissingFieldScriptPubKey is when the field is required but missing
var ErrMissingFieldScriptPubKey = errors.New("missing required field: script_pub_key")

// ErrMissingFieldSatoshis is when the field is required but missing
var ErrMissingFieldSatoshis = errors.New("missing required field: satoshis")

// ErrMissingFieldTransactionID is when the field is required but missing
var ErrMissingFieldTransactionID = errors.New("missing required field: transaction_id")

// ErrMissingFieldXpubID is when the field is required but missing
var ErrMissingFieldXpubID = errors.New("missing required field: xpub_id")

// ErrXpubIDMisMatch is when the xPubID does not match
var ErrXpubIDMisMatch = errors.New("xpub_id mismatch")

// ErrMissingXpub is when the field is required but missing
var ErrMissingXpub = errors.New("could not find xpub")

// ErrAccessKeyNotFound is when the access key not found
var ErrAccessKeyNotFound = errors.New("access key not found")

// ErrMissingLockingScript is when the field is required but missing
var ErrMissingLockingScript = errors.New("could not find locking script")

// ErrUnknownLockingScript is when the field is unknown
var ErrUnknownLockingScript = errors.New("could not recognize locking script")

// ErrMissingRequiredXpub is when the xpub should exist but was not found
var ErrMissingRequiredXpub = errors.New("xpub was not found but was expected")

// ErrDatastoreRequired is when a datastore function is called without a datastore present
var ErrDatastoreRequired = errors.New("datastore is required")

// ErrMissingTransactionOutputs is when the draft transaction has not outputs
var ErrMissingTransactionOutputs = errors.New("draft transaction configuration has no outputs")

// ErrNotEnoughUtxos is when a draft transaction cannot be created because of lack of utxos
var ErrNotEnoughUtxos = errors.New("could not select enough outputs to satisfy transaction")

// ErrInvalidLockingScript is when a locking script cannot be decoded
var ErrInvalidLockingScript = errors.New("invalid locking script")

// ErrInvalidOpReturnOutput is when a locking script is not a valid op_return
var ErrInvalidOpReturnOutput = errors.New("invalid op_return output")

// ErrInvalidScriptOutput is when a locking script is not a valid bitcoin script
var ErrInvalidScriptOutput = errors.New("invalid script output")

// ErrInvalidTransactionID is when a transaction id cannot be decoded
var ErrInvalidTransactionID = errors.New("invalid transaction id")

// ErrOutputValueNotRecognized is when there is an invalid output value given, or missing value
var ErrOutputValueNotRecognized = errors.New("output value is unrecognized")

// ErrOutputValueTooLow is when the satoshis output is too low on a transaction
var ErrOutputValueTooLow = errors.New("output value is too low")

// ErrOutputValueTooHigh is when the satoshis output is too high on a transaction
var ErrOutputValueTooHigh = errors.New("output value is too high")

// ErrTransactionFeeInvalid is when the fee on the transaction is not the difference between inputs and outputs
var ErrTransactionFeeInvalid = errors.New("transaction fee is invalid")

// ErrMissingUTXOsSpendable is when there are no utxos found from the "spendable utxos"
var ErrMissingUTXOsSpendable = errors.New("no utxos found using spendable")

// ErrDuplicateUTXOs is when a transaction is created using the same utxo more than once
var ErrDuplicateUTXOs = errors.New("duplicate utxos found")

// ErrPaymailAddressIsInvalid is when the paymail address is NOT alias@domain.com
var ErrPaymailAddressIsInvalid = errors.New("paymail address is invalid")

// ErrPaymailNotFound is when paymaail could not be found
var ErrPaymailNotFound = errors.New("paymail could not be found")

// ErrUtxoNotReserved is when the utxo is not reserved, but a transaction tries to spend it
var ErrUtxoNotReserved = errors.New("transaction utxo has not been reserved for spending")

// ErrDraftIDMismatch is when the reference ID does not match the reservation id
var ErrDraftIDMismatch = errors.New("transaction draft id does not match utxo draft reservation id")

// ErrMissingTxHex is when the hex is missing or invalid and creates an empty id
var ErrMissingTxHex = errors.New("transaction hex is empty or id is missing")

// ErrUtxoAlreadySpent is when the utxo is already spent, but is trying to be used
var ErrUtxoAlreadySpent = errors.New("utxo has already been spent")

// ErrDraftNotFound is when the requested draft transaction was not found
var ErrDraftNotFound = errors.New("corresponding draft transaction not found")

// ErrTransactionNotParsed is when the transaction is not parsed but was expected
var ErrTransactionNotParsed = errors.New("transaction is not parsed")

// ErrNoMatchingOutputs is when the transaction does not match any known destinations
var ErrNoMatchingOutputs = errors.New("transaction outputs do not match any known destinations")

// ErrResolutionFailed is when the paymail resolution failed unexpectedly
var ErrResolutionFailed = errors.New("failed to return a resolution for paymail address")

// ErrMissingAddressResolutionURL is when the paymail resolution url is missing from capabilities
var ErrMissingAddressResolutionURL = errors.New("missing address resolution url from capabilities")

// ErrChangeStrategyNotImplemented is a temporary error until the feature is supported
var ErrChangeStrategyNotImplemented = errors.New("change strategy nominations not implemented yet")

// ErrUnsupportedDestinationType is a destination type that is not currently supported
var ErrUnsupportedDestinationType = errors.New("unsupported destination type")

// ErrMissingAuthHeader is when the authentication header is missing from the request
var ErrMissingAuthHeader = errors.New("missing authentication header")

// ErrMissingSignature is when the signature is missing from the request
var ErrMissingSignature = errors.New("signature missing")

// ErrAuhHashMismatch is when the auth hash does not match the body hash
var ErrAuhHashMismatch = errors.New("auth hash and body hash do not match")

// ErrAuthAccessKeyNotFound is when the auth access key could not be found in the database
var ErrAuthAccessKeyNotFound = errors.New("auth access key could not be found")

// ErrSignatureExpired is when the signature TTL expired
var ErrSignatureExpired = errors.New("signature has expired")

// ErrNotAdminKey is when the xpub being used is not considered an admin key
var ErrNotAdminKey = errors.New("xpub provided is not an admin key")

// ErrMissingXPriv is when the xPriv is missing
var ErrMissingXPriv = errors.New("missing xPriv key")

// ErrMissingAccessKey is when the access key is missing
var ErrMissingAccessKey = errors.New("missing access key")

// ErrMissingBody is when the body is missing
var ErrMissingBody = errors.New("missing body")

// ErrSignatureInvalid is when the signature failed to be valid
var ErrSignatureInvalid = errors.New("signature invalid")

// ErrUnknownAccessKey is when the access key is unknown or not found
var ErrUnknownAccessKey = errors.New("unknown access key")

// ErrAccessKeyRevoked is when the access key has been revoked
var ErrAccessKeyRevoked = errors.New("access key has been revoked")

// ErrMissingPaymailID missing id in paymail
var ErrMissingPaymailID = errors.New("missing id in paymail")

// ErrMissingPaymailAddress missing alias in paymail
var ErrMissingPaymailAddress = errors.New("missing alias in paymail")

// ErrMissingPaymailDomain missing domain in paymail
var ErrMissingPaymailDomain = errors.New("missing domain in paymail")

// ErrMissingPaymailExternalXPub missing external xPub in paymail
var ErrMissingPaymailExternalXPub = errors.New("missing external xPub in paymail")

// ErrMissingPaymailIdentityXPub missing identity xPub in paymail
// var ErrMissingPaymailIdentityXPub = errors.New("missing identity xPub in paymail")

// ErrMissingPaymailXPubID missing xpub_id in paymail
var ErrMissingPaymailXPubID = errors.New("missing xpub_id in paymail")

// ErrMissingClient missing client from model
var ErrMissingClient = errors.New("client is missing from model, cannot save")

// ErrMissingContactID missing id in contact
var ErrMissingContactID = errors.New("missing id in contact")

// ErrMissingContactFullName missing full name in contact
var ErrMissingContactFullName = errors.New("missing full_name in contact")

// ErrMissingContactPaymail missing paymail in contact
var ErrMissingContactPaymail = errors.New("missing paymail in contact")

// ErrMissingContactXPubKey missing XPubKey in contact
var ErrMissingContactXPubKey = errors.New("missing pubKey in contact")

// ErrMissingContactStatus missing status in contact
var ErrMissingContactStatus = errors.New("status is required")

// ErrMissingContactStatus missing status in contact
var ErrMissingContactOwnerXPubId = errors.New("contact must have owner")

// ErrEmptyContactPubKey when pubKey is empty
var ErrEmptyContactPubKey = errors.New("pubKey is empty")

// ErrEmptyContactPaymail when paymail is empty
var ErrEmptyContactPaymail = errors.New("paymail is empty")

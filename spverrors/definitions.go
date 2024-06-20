package spverrors

import (
	"errors"
)

// NOTE: All errors should be implemented in SPVErrorResponses in errors.go

//////////////////////////////////// AUTHORIZATION ERRORS

// ErrAuthorization is basic auth error
var ErrAuthorization = errors.New("unauthorized")

// ErrMissingAuthHeader is when request does not have auth header
var ErrMissingAuthHeader = errors.New("missing auth header")

// ErrNotAnAdminKey is when xpub from auth header is not an admin key
var ErrNotAnAdminKey = errors.New("xpub provided is not an admin key")

// ErrMissingBody is when request is missing body
var ErrMissingBody = errors.New("missing body")

// ErrCheckSignature is when error occurred during checking signature
var ErrCheckSignature = errors.New("error occurred during checking signature")

// ErrInvalidOrMissingToken is when callback token from headers is invalid or missing
var ErrInvalidOrMissingToken = errors.New("invalid or missing bearer token")

// ErrInvalidToken is when callback token from headers is invalid
var ErrInvalidToken = errors.New("invalid authorization token")

//////////////////////////////////// BINDING ERRORS

// ErrCannotBindRequest is when request body cannot be bind into struct
var ErrCannotBindRequest = errors.New("cannot bin request body")

// ErrInvalidFilterOption is when filter has invalid option
var ErrInvalidFilterOption = errors.New("invalid filter option")

// ErrInvalidConditions is when request has invalid conditions
var ErrInvalidConditions = errors.New("invalid conditions")

//////////////////////////////////// ACCESS KEY ERRORS

// ErrCouldNotFindAccessKey is when could not find xpub
var ErrCouldNotFindAccessKey = errors.New("access key not found")

// ErrAccessKeyRevoked is when the access key has been revoked
var ErrAccessKeyRevoked = errors.New("access key has been revoked")

//////////////////////////////////// DESTINATION ERRORS

// ErrCouldNotFindDestination is an error when a destination could not be found
var ErrCouldNotFindDestination = errors.New("destination not found")

// ErrUnsupportedDestinationType is a destination type that is not currently supported
var ErrUnsupportedDestinationType = errors.New("unsupported destination type")

// ErrUnknownLockingScript is when the field is unknown
var ErrUnknownLockingScript = errors.New("could not recognize locking script")

//////////////////////////////////// CONTACT ERRORS

// ErrInvalidRequesterXpub is when requester xpub is not connected with given paymail
var ErrInvalidRequesterXpub = errors.New("invalid requester xpub")

// ErrAddingContactRequest is when error occurred while adding contact
var ErrAddingContactRequest = errors.New("adding contact request failed")

// ErrMoreThanOnePaymailRegistered is when user who want to add contact has more than one paymail address
var ErrMoreThanOnePaymailRegistered = errors.New("there are more than one paymail assigned to the xpub")

// ErrContactNotFound is when contact cannot be found
var ErrContactNotFound = errors.New("contact not found")

// ErrContactIncorrectStatus is when contact is in incorrect status to make a change
var ErrContactIncorrectStatus = errors.New("contact is in incorrect status to proceed")

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

// ErrMissingContactOwnerXPubId missing status in contact
var ErrMissingContactOwnerXPubId = errors.New("contact must have owner")

// ErrEmptyContactPubKey when pubKey is empty
var ErrEmptyContactPubKey = errors.New("pubKey is empty")

// ErrEmptyContactPaymail when paymail is empty
var ErrEmptyContactPaymail = errors.New("paymail is empty")

//////////////////////////////////// PAYMAIL ERRORS

// ErrCouldNotFindPaymail is when paymail could not be found
var ErrCouldNotFindPaymail = errors.New("paymail not found")

// ErrPaymailAddressIsInvalid is when the paymail address is NOT alias@domain.com
var ErrPaymailAddressIsInvalid = errors.New("paymail address is invalid")

// ErrMissingPaymailID missing id in paymail
var ErrMissingPaymailID = errors.New("missing id in paymail")

// ErrMissingPaymailAddress missing alias in paymail
var ErrMissingPaymailAddress = errors.New("missing alias in paymail")

// ErrMissingPaymailDomain missing domain in paymail
var ErrMissingPaymailDomain = errors.New("missing domain in paymail")

// ErrMissingPaymailExternalXPub missing external xPub in paymail
var ErrMissingPaymailExternalXPub = errors.New("missing external xPub in paymail")

// ErrMissingPaymailXPubID missing xpub_id in paymail
var ErrMissingPaymailXPubID = errors.New("missing xpub_id in paymail")

// ErrPaymailAlreadyExists is when paymail with given data already exists in db
var ErrPaymailAlreadyExists = errors.New("paymail already exists")

//////////////////////////////////// TRANSACTION ERRORS

// ErrCouldNotFindTransaction is an error when a transaction could not be found
var ErrCouldNotFindTransaction = errors.New("transaction not be found")

// ErrInvalidTransactionID is when a transaction id cannot be decoded
var ErrInvalidTransactionID = errors.New("invalid transaction id")

// ErrInvalidRequirements is when an invalid requirement was given
var ErrInvalidRequirements = errors.New("requirements are invalid or missing")

// ErrTransactionIDMismatch is when the returned tx does not match the expected given tx id
var ErrTransactionIDMismatch = errors.New("result tx id did not match provided tx id")

// ErrMissingTransactionOutputs is when the draft transaction has not outputs
var ErrMissingTransactionOutputs = errors.New("draft transaction configuration has no outputs")

// ErrOutputValueTooLow is when the satoshis output is too low on a transaction
var ErrOutputValueTooLow = errors.New("output value is too low")

// ErrOutputValueTooHigh is when the satoshis output is too high on a transaction
var ErrOutputValueTooHigh = errors.New("output value is too high")

// ErrInvalidOpReturnOutput is when a locking script is not a valid op_return
var ErrInvalidOpReturnOutput = errors.New("invalid op_return output")

// ErrInvalidLockingScript is when a locking script cannot be decoded
var ErrInvalidLockingScript = errors.New("invalid locking script")

// ErrOutputValueNotRecognized is when there is an invalid output value given, or missing value
var ErrOutputValueNotRecognized = errors.New("output value is unrecognized")

// ErrInvalidScriptOutput is when a locking script is not a valid bitcoin script
var ErrInvalidScriptOutput = errors.New("invalid script output")

// ErrDraftIDMismatch is when the reference ID does not match the reservation id
var ErrDraftIDMismatch = errors.New("transaction draft id does not match utxo draft reservation id")

// ErrMissingTxHex is when the hex is missing or invalid and creates an empty id
var ErrMissingTxHex = errors.New("transaction hex is empty or id is missing")

// ErrNoMatchingOutputs is when the transaction does not match any known destinations
var ErrNoMatchingOutputs = errors.New("transaction outputs do not match any known destinations")

// ErrCouldNotFindSyncTx is an error when a given utxo could not be found
var ErrCouldNotFindSyncTx = errors.New("sync tx not found")

// ErrCouldNotFindDraftTx is an error when a given draft tx could not be found
var ErrCouldNotFindDraftTx = errors.New("draft tx not found")

// ErrCreateOutgoingTxFailed is when error occurred during creation of outgoing tx
var ErrCreateOutgoingTxFailed = errors.New("creation of outgoing tx failed")

// ErrDuringSaveTx is when error occurred during save tx
var ErrDuringSaveTx = errors.New("error during saving tx")

// ErrTransactionRejectedByP2PProvider is an error when a tx was rejected by P2P Provider
var ErrTransactionRejectedByP2PProvider = errors.New("transaction rejected by P2P provider")

// ErrDraftTxHasNoOutputs is when draft transaction has no outputs
var ErrDraftTxHasNoOutputs = errors.New("corresponding draft transaction has no outputs")

// ErrProcessP2PTx is when error occurred during processing p2p tx
var ErrProcessP2PTx = errors.New("error during processing p2p transaction")

//////////////////////////////////// UTXO ERRORS

// ErrCouldNotFindUtxo is an error when a given utxo could not be found
var ErrCouldNotFindUtxo = errors.New("utxo could not be found")

// ErrUtxoAlreadySpent is when the utxo is already spent, but is trying to be used
var ErrUtxoAlreadySpent = errors.New("utxo has already been spent")

// ErrMissingUTXOsSpendable is when there are no utxos found from the "spendable utxos"
var ErrMissingUTXOsSpendable = errors.New("no utxos found using spendable")

// ErrNotEnoughUtxos is when a draft transaction cannot be created because of lack of utxos
var ErrNotEnoughUtxos = errors.New("could not select enough outputs to satisfy transaction")

// ErrDuplicateUTXOs is when a transaction is created using the same utxo more than once
var ErrDuplicateUTXOs = errors.New("duplicate utxos found")

// ErrTransactionFeeInvalid is when the fee on the transaction is not the difference between inputs and outputs
var ErrTransactionFeeInvalid = errors.New("transaction fee is invalid")

// ErrChangeStrategyNotImplemented is a temporary error until the feature is supported
var ErrChangeStrategyNotImplemented = errors.New("change strategy nominations not implemented yet")

// ErrUtxoNotReserved is when the utxo is not reserved, but a transaction tries to spend it
var ErrUtxoNotReserved = errors.New("transaction utxo has not been reserved for spending")

//////////////////////////////////// XPUB ERRORS

// ErrXpubInvalidLength is when the length of the xpub does not match the desired length
var ErrXpubInvalidLength = errors.New("xpub is an invalid length")

// ErrXpubNoMatch is when the derived xpub key does not match the key given
var ErrXpubNoMatch = errors.New("xpub key does not match raw key")

// ErrCouldNotFindXpub is when could not find xpub
var ErrCouldNotFindXpub = errors.New("xpub not found")

// ErrXpubIDMisMatch is when the xPubID does not match
var ErrXpubIDMisMatch = errors.New("xpub_id mismatch")

//////////////////////////////////// MISSING FIELDS

// ErrOneOfTheFieldsIsRequired is when all of required fields are missing
var ErrOneOfTheFieldsIsRequired = errors.New("missing all of the fields, one of them is required")

// ErrMissingAccessKey is when the access key field is required but missing
var ErrMissingAccessKey = errors.New("missing required field: access key")

// ErrMissingFieldID is when the id field is required but missing
var ErrMissingFieldID = errors.New("missing required field: id")

// ErrMissingFieldXpubID is when the xpub_id field is required but missing
var ErrMissingFieldXpubID = errors.New("missing required field: xpub_id")

// ErrMissingFieldXpub is when the xpub field is required but missing
var ErrMissingFieldXpub = errors.New("missing required field: xpub")

// ErrMissingAddress  is when the address field address is required but missing
var ErrMissingAddress = errors.New("missing required field: address")

// ErrMissingFieldScriptPubKey is when the field is required but missing
var ErrMissingFieldScriptPubKey = errors.New("missing required field: script_pub_key")

// ErrMissingFieldSatoshis is when the field satoshis is required but missing
var ErrMissingFieldSatoshis = errors.New("missing required field: satoshis")

// ErrMissingFieldTransactionID is when the field transaction id is required but missing
var ErrMissingFieldTransactionID = errors.New("missing required field: transaction_id")

// ErrMissingLockingScript is when the field locking script is required but missing
var ErrMissingLockingScript = errors.New("missing required field: locking script")

//////////////////////////////////// SAVE ERROR

// ErrMissingClient missing client from model
var ErrMissingClient = errors.New("client is missing from model, cannot save")

// ErrDatastoreRequired is when a datastore function is called without a datastore present
var ErrDatastoreRequired = errors.New("datastore is required")

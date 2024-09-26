package spverrors

import "github.com/bitcoin-sv/spv-wallet/models"

// How the Codes are generated?
// 1. "error" - like mandatory prefix for all error codes
// 2. (optional) {error group} - e.g. "unauthorized", "bind"
// 3. (optional) {subject} - name of model (with or without specific field) or some noun e.g. "body", "auth-header", "transaction", "paymail-address"
// 4. (optional) {reason} - what happened, e.g. "not-found", "missing", "invalid"

// ////////////////////////////////// AUTHORIZATION ERRORS

// ErrAuthorization is basic auth error
var ErrAuthorization = models.SPVError{Message: "unauthorized", StatusCode: 401, Code: "error-unauthorized"}

// ErrMissingAuthHeader is when request does not have auth header
var ErrMissingAuthHeader = models.SPVError{Message: "missing auth header", StatusCode: 401, Code: "error-unauthorized-auth-header-missing"}

// ErrNotAnAdminKey is when xpub from auth header is not an admin key
var ErrNotAnAdminKey = models.SPVError{Message: "xpub provided is not an admin key", StatusCode: 401, Code: "error-unauthorized-xpub-not-an-admin-key"}

// ErrInternal is a generic error that something weird went wrong
var ErrInternal = models.SPVError{Message: "internal server error", StatusCode: 500, Code: "error-internal-server-error"}

// ErrInvalidOrMissingToken is when callback token from headers is invalid or missing
var ErrInvalidOrMissingToken = models.SPVError{Message: "invalid or missing bearer token", StatusCode: 401, Code: "error-unauthorized-token-invalid-or-missing"}

// ErrInvalidToken is when callback token from headers is invalid
var ErrInvalidToken = models.SPVError{Message: "invalid authorization token", StatusCode: 401, Code: "error-unauthorized-token-invalid"}

// ErrInvalidSignature is when signature is invalid
var ErrInvalidSignature = models.SPVError{Message: "invalid signature", StatusCode: 401, Code: "error-unauthorized-signature-invalid"}

// ErrMissingSignature is when signature is missing in authorization process
var ErrMissingSignature = models.SPVError{Message: "missing signature", StatusCode: 401, Code: "error-unauthorized-signature-missing"}

// ErrSignatureExpired is when given signature is expired
var ErrSignatureExpired = models.SPVError{Message: "signature has expired", StatusCode: 401, Code: "error-unauthorized-signature-expired"}

// ErrDeriveChildKey is when error occurred during deriving child key
var ErrDeriveChildKey = models.SPVError{Message: "error deriving child key", StatusCode: 401, Code: "error-unauthorized-derive-child-key"}

// ErrGettingAddressFromHdKey is when error occurred during getting address from hd key
var ErrGettingAddressFromHdKey = models.SPVError{Message: "error getting address from hd key", StatusCode: 401, Code: "error-unauthorized-address-failed-to-get-from-hd-key"}

// ErrValidateXPub is when validation xpub
var ErrValidateXPub = models.SPVError{Message: "error validation xpub", StatusCode: 401, Code: "error-unauthorized-xpub-invalid"}

// ErrAdminAuthOnUserEndpoint is when admin authorization is used on user's endpoint
var ErrAdminAuthOnUserEndpoint = models.SPVError{Message: "cannot call user's endpoints with admin authorization ", StatusCode: 401, Code: "error-admin-auth-on-user-endpoint"}

// ErrXPubAuthRequired is when xpub authorization is required
var ErrXPubAuthRequired = models.SPVError{Message: "xpub authorization required", StatusCode: 401, Code: "error-xpub-authorization-required"}

// ////////////////////////////////// BINDING ERRORS

// ErrCannotBindRequest is when request body cannot be bind into struct
var ErrCannotBindRequest = models.SPVError{Message: "cannot bind request body", StatusCode: 400, Code: "error-bind-body-invalid"}

// ErrCannotParseQueryParams is when query params cannot be parsed into expected struct.
var ErrCannotParseQueryParams = models.SPVError{Message: "cannot parse request query params", StatusCode: 400, Code: "error-query-params-invalid"}

// ErrInvalidConditions is when request has invalid conditions
var ErrInvalidConditions = models.SPVError{Message: "invalid conditions", StatusCode: 400, Code: "error-bind-conditions-invalid"}

// ////////////////////////////////// MAPPING ERRORS

// ErrCannotMapFromModel is when request body model cannot be mapped into domain model.
var ErrCannotMapFromModel = models.SPVError{Message: "error during reading request body", StatusCode: 500, Code: "error-request-read"}

// ////////////////////////////////// ACCESS KEY ERRORS

// ErrCouldNotFindAccessKey is when could not find xpub
var ErrCouldNotFindAccessKey = models.SPVError{Message: "access key not found", StatusCode: 404, Code: "error-access-key-not-found"}

// ErrAccessKeyRevoked is when the access key has been revoked
var ErrAccessKeyRevoked = models.SPVError{Message: "access key has been revoked", StatusCode: 400, Code: "error-access-key-revoked"}

// ////////////////////////////////// DESTINATION ERRORS

// ErrCouldNotFindDestination is an error when a destination could not be found
var ErrCouldNotFindDestination = models.SPVError{Message: "destination not found", StatusCode: 404, Code: "error-destination-not-found"}

// ErrUnsupportedDestinationType is a destination type that is not currently supported
var ErrUnsupportedDestinationType = models.SPVError{Message: "unsupported destination type", StatusCode: 400, Code: "error-destination-unsupported-type"}

// ErrUnknownLockingScript is when the field is unknown
var ErrUnknownLockingScript = models.SPVError{Message: "could not recognize locking script", StatusCode: 400, Code: "error-destination-unknown-locking-script"}

// ////////////////////////////////// CONTACT ERRORS

// ErrContactNotFound is when contact cannot be found
var ErrContactNotFound = models.SPVError{Message: "contact not found", StatusCode: 404, Code: "error-contact-not-found"}

// ErrInvalidRequesterXpub is when requester xpub is not connected with given paymail
var ErrInvalidRequesterXpub = models.SPVError{Message: "invalid requester xpub", StatusCode: 400, Code: "error-contact-invalid-requester-xpub"}

// ErrAddingContactRequest is when error occurred while adding contact
var ErrAddingContactRequest = models.SPVError{Message: "adding contact request failed", StatusCode: 500, Code: "error-contact-request-failed"}

// ErrMoreThanOnePaymailRegistered is when user who want to add contact has more than one paymail address
var ErrMoreThanOnePaymailRegistered = models.SPVError{Message: "there are more than one paymail assigned to the xpub", StatusCode: 400, Code: "error-contact-more-than-one-paymail-registered"}

// ErrContactIncorrectStatus is when contact is in incorrect status to make a change
var ErrContactIncorrectStatus = models.SPVError{Message: "contact is in incorrect status to proceed", StatusCode: 400, Code: "error-contact-status-incorrect"}

// ErrMissingContactID is when id is missing in contact
var ErrMissingContactID = models.SPVError{Message: "missing id in contact", StatusCode: 400, Code: "error-contact-id-missing"}

// ErrMissingContactFullName is when full name is missing in contact
var ErrMissingContactFullName = models.SPVError{Message: "missing full name in contact", StatusCode: 400, Code: "error-contact-full-name-missing"}

// ErrInvalidContactPaymail is when paymail is missing in contact
var ErrInvalidContactPaymail = models.SPVError{Message: "invalid paymail in contact", StatusCode: 400, Code: "error-contact-paymail-invalid"}

// ErrMissingContactXPubKey is when XPubKey is missing in contact
var ErrMissingContactXPubKey = models.SPVError{Message: "missing pubKey in contact", StatusCode: 400, Code: "error-contact-xpub-missing"}

// ErrMissingContactStatus is when status is missing in contact
var ErrMissingContactStatus = models.SPVError{Message: "status is required", StatusCode: 400, Code: "error-contact-status-missing"}

// ErrMissingContactOwnerXPubID is when owner XPubID is missing in contact
var ErrMissingContactOwnerXPubID = models.SPVError{Message: "contact must have owner", StatusCode: 400, Code: "error-contact-owner-xpub-id-missing"}

// ErrRequestedContactInvalid is when the requested contact is invalid
var ErrRequestedContactInvalid = models.SPVError{Message: "requested contact paymail is invalid", StatusCode: 400, Code: "error-contact-requested-contact-invalid"}

// ErrGettingPKIFailed is when getting PKI for contact paymail failed
var ErrGettingPKIFailed = models.SPVError{Message: "getting PKI for contact failed", StatusCode: 400, Code: "error-contact-getting-pki-failed"}

// ErrSaveContact is when saving new contact failed
var ErrSaveContact = models.SPVError{Message: "adding contact failed", StatusCode: 400, Code: "error-contact-adding-contact-failed"}

// ////////////////////////////////// PAYMAIL ERRORS

// ErrCouldNotFindPaymail is when paymail could not be found
var ErrCouldNotFindPaymail = models.SPVError{Message: "paymail not found", StatusCode: 404, Code: "error-paymail-not-found"}

// ErrPaymailAddressIsInvalid is when the paymail address is NOT alias@domain.com
var ErrPaymailAddressIsInvalid = models.SPVError{Message: "paymail address is invalid", StatusCode: 400, Code: "error-paymail-address-invalid"}

// ErrMissingPaymailID is when id is missing in paymail
var ErrMissingPaymailID = models.SPVError{Message: "missing id in paymail", StatusCode: 400, Code: "error-paymail-id-missing"}

// ErrMissingPaymailAddress is when alias is missing in paymail
var ErrMissingPaymailAddress = models.SPVError{Message: "missing alias in paymail", StatusCode: 400, Code: "error-paymail-address-missing"}

// ErrMissingPaymailDomain is when domain is missing in paymail
var ErrMissingPaymailDomain = models.SPVError{Message: "missing domain in paymail", StatusCode: 400, Code: "error-paymail-domain-missing"}

// ErrMissingPaymailExternalXPub is when external xPub is missing in paymail
var ErrMissingPaymailExternalXPub = models.SPVError{Message: "missing external xPub in paymail", StatusCode: 400, Code: "error-paymail-external-xpub-missing"}

// ErrMissingPaymailXPubID is when xpub_id is missing in paymail
var ErrMissingPaymailXPubID = models.SPVError{Message: "missing xpub_id in paymail", StatusCode: 400, Code: "error-paymail-xpub-id-missing"}

// ErrPaymailAlreadyExists is when paymail with given data already exists in db
var ErrPaymailAlreadyExists = models.SPVError{Message: "paymail already exists", StatusCode: 409, Code: "error-paymail-already-exists"}

// ////////////////////////////////// CAPABILITIES ERRORS

// ErrCapabilitiesPkiUnsupported is when PKI is not supported for given paymail domain
var ErrCapabilitiesPkiUnsupported = models.SPVError{Message: "server doesn't support PKI", StatusCode: 400, Code: "error-capabilities-pki-unsupported"}

// ErrCapabilitiesPikeUnsupported is when PIKE is not supported for given paymail domain
var ErrCapabilitiesPikeUnsupported = models.SPVError{Message: "server doesn't support PIKE", StatusCode: 400, Code: "error-capabilities-pike-unsupported"}

// ErrGetCapabilities is when getting capabilities failed
var ErrGetCapabilities = models.SPVError{Message: "failed to get paymail capabilities", StatusCode: 400, Code: "error-capabilities-failed-to-get"}

// ////////////////////////////////// TRANSACTION ERRORS

// ErrCouldNotFindTransaction is an error when a transaction could not be found
var ErrCouldNotFindTransaction = models.SPVError{Message: "transaction not found", StatusCode: 404, Code: "error-transaction-not-found"}

// ErrCouldNotFindSyncTx is an error when a given utxo could not be found
var ErrCouldNotFindSyncTx = models.SPVError{Message: "sync tx not found", StatusCode: 404, Code: "error-transaction-sync-tx-not-found"}

// ErrCouldNotFindDraftTx is an error when a given draft tx could not be found
var ErrCouldNotFindDraftTx = models.SPVError{Message: "draft tx not found", StatusCode: 404, Code: "error-transaction-draft-tx-not-found"}

// ErrInvalidTransactionID is when a transaction id cannot be decoded
var ErrInvalidTransactionID = models.SPVError{Message: "invalid transaction id", StatusCode: 400, Code: "error-transaction-id-invalid"}

// ErrInvalidRequirements is when an invalid requirement was given
var ErrInvalidRequirements = models.SPVError{Message: "requirements are invalid or missing", StatusCode: 400, Code: "error-transaction-requirements-invalid"}

// ErrTransactionIDMismatch is when the returned tx does not match the expected given tx id
var ErrTransactionIDMismatch = models.SPVError{Message: "result tx id did not match provided tx id", StatusCode: 400, Code: "error-transaction-id-mismatch"}

// ErrMissingTransactionOutputs is when the draft transaction has no outputs
var ErrMissingTransactionOutputs = models.SPVError{Message: "draft transaction configuration has no outputs", StatusCode: 400, Code: "error-transaction-outputs-missing"}

// ErrOutputValueTooLow is when the satoshis output is too low on a transaction
var ErrOutputValueTooLow = models.SPVError{Message: "output value is too low", StatusCode: 400, Code: "error-transaction-output-value-too-low"}

// ErrOutputValueTooHigh is when the satoshis output is too high on a transaction
var ErrOutputValueTooHigh = models.SPVError{Message: "output value is too high", StatusCode: 400, Code: "error-transaction-output-value-too-high"}

// ErrInvalidOpReturnOutput is when a locking script is not a valid op_return
var ErrInvalidOpReturnOutput = models.SPVError{Message: "invalid op_return output", StatusCode: 400, Code: "error-transaction-invalid-op-return-output"}

// ErrInvalidLockingScript is when a locking script cannot be decoded
var ErrInvalidLockingScript = models.SPVError{Message: "invalid locking script", StatusCode: 400, Code: "error-transaction-locking-script-invalid"}

// ErrOutputValueNotRecognized is when there is an invalid output value given, or missing value
var ErrOutputValueNotRecognized = models.SPVError{Message: "output value is unrecognized", StatusCode: 400, Code: "error-transaction-output-value-unrecognized"}

// ErrInvalidScriptOutput is when a locking script is not a valid bitcoin script
var ErrInvalidScriptOutput = models.SPVError{Message: "invalid script output", StatusCode: 400, Code: "error-transaction-script-output-invalid"}

// ErrDraftIDMismatch is when the reference ID does not match the reservation id
var ErrDraftIDMismatch = models.SPVError{Message: "transaction draft id does not match utxo draft reservation id", StatusCode: 400, Code: "error-transaction-draft-id-mismatch"}

// ErrMissingTxHex is when the hex is missing or invalid and creates an empty id
var ErrMissingTxHex = models.SPVError{Message: "transaction hex is empty or id is missing", StatusCode: 400, Code: "error-transaction-hex-missing"}

// ErrNoMatchingOutputs is when the transaction does not match any known destinations
var ErrNoMatchingOutputs = models.SPVError{Message: "transaction outputs do not match any known destinations", StatusCode: 400, Code: "error-transaction-outputs-no-matching"}

// ErrCreateOutgoingTxFailed is when error occurred during creation of outgoing tx
var ErrCreateOutgoingTxFailed = models.SPVError{Message: "creation of outgoing tx failed", StatusCode: 500, Code: "error-transaction-create-outgoing-tx-failed"}

// ErrDuringSaveTx is when error occurred during save tx
var ErrDuringSaveTx = models.SPVError{Message: "error during saving tx", StatusCode: 500, Code: "error-transaction-save-failed"}

// ErrTransactionRejectedByP2PProvider is an error when a tx was rejected by P2P Provider
var ErrTransactionRejectedByP2PProvider = models.SPVError{Message: "transaction rejected by P2P provider", StatusCode: 400, Code: "error-transaction-rejected"}

// ErrDraftTxHasNoOutputs is when draft transaction has no outputs
var ErrDraftTxHasNoOutputs = models.SPVError{Message: "corresponding draft transaction has no outputs", StatusCode: 400, Code: "error-transaction-draft-has-no-outputs"}

// ErrProcessP2PTx is when error occurred during processing p2p tx
var ErrProcessP2PTx = models.SPVError{Message: "error during processing p2p transaction", StatusCode: 500, Code: "error-transaction-process-p2p"}

// ErrInvalidHex is when cannot create tx from hex
var ErrInvalidHex = models.SPVError{Message: "invalid hex", StatusCode: 400, Code: "error-transaction-hex-invalid"}

// ErrEmptyRelatedDraftID is when related draft id is empty
var ErrEmptyRelatedDraftID = models.SPVError{Message: "empty RelatedDraftID", StatusCode: 400, Code: "error-transaction-related-draft-id-empty"}

// ErrEmptyXpubKey is when xpub key is empty
var ErrEmptyXpubKey = models.SPVError{Message: "empty xPubKey", StatusCode: 400, Code: "error-transaction-xpub-key-empty"}

// ErrEmptyTx is when tx is empty
var ErrEmptyTx = models.SPVError{Message: "empty tx", StatusCode: 400, Code: "error-transaction-empty"}

// ErrTxRevertEmptyDraftID is when draft id is empty this means that tx is not from spv-wallet
var ErrTxRevertEmptyDraftID = models.SPVError{Message: "not a spv wallet engine originating transaction, cannot revert", StatusCode: 400, Code: "error-transaction-revert-draft-id-empty"}

// ErrTxRevertCouldNotFindDraftTx is when draft tx could not be found
var ErrTxRevertCouldNotFindDraftTx = models.SPVError{Message: "could not find the draft transaction for this transaction, cannot revert", StatusCode: 400, Code: "error-transaction-revert-draft-tx-not-found"}

// ErrTxRevertNotFoundOnChain is when tx was not found on chain
var ErrTxRevertNotFoundOnChain = models.SPVError{Message: "transaction was found on-chain, cannot revert", StatusCode: 400, Code: "error-transaction-revert-not-found-on-chain"}

// ErrTxRevertUtxoAlreadySpent is when utxo from tx was already spent
var ErrTxRevertUtxoAlreadySpent = models.SPVError{Message: "utxo of this transaction has been spent, cannot revert", StatusCode: 400, Code: "error-transaction-revert-utxo-already-spent"}

// ////////////////////////////////// UTXO ERRORS

// ErrCouldNotFindUtxo is an error when a given utxo could not be found
var ErrCouldNotFindUtxo = models.SPVError{Message: "utxo could not be found", StatusCode: 404, Code: "error-utxo-not-found"}

// ErrUtxoAlreadySpent is when the utxo is already spent, but is trying to be used
var ErrUtxoAlreadySpent = models.SPVError{Message: "utxo has already been spent", StatusCode: 400, Code: "error-utxo-already-spent"}

// ErrMissingUTXOsSpendable is when there are no utxos found from the "spendable utxos"
var ErrMissingUTXOsSpendable = models.SPVError{Message: "no utxos found using spendable", StatusCode: 404, Code: "error-utxo-spendable-missing"}

// ErrNotEnoughUtxos is when a draft transaction cannot be created because of lack of utxos
var ErrNotEnoughUtxos = models.SPVError{Message: "could not select enough outputs to satisfy transaction", StatusCode: 400, Code: "error-utxo-not-enough"}

// ErrDuplicateUTXOs is when a transaction is created using the same utxo more than once
var ErrDuplicateUTXOs = models.SPVError{Message: "duplicate utxos found", StatusCode: 400, Code: "error-utxo-duplicate"}

// ErrTransactionFeeInvalid is when the fee on the transaction is not the difference between inputs and outputs
var ErrTransactionFeeInvalid = models.SPVError{Message: "transaction fee is invalid", StatusCode: 400, Code: "error-utxo-transaction-fee-invalid"}

// ErrChangeStrategyNotImplemented is a temporary error until the feature is supported
var ErrChangeStrategyNotImplemented = models.SPVError{Message: "change strategy nominations not implemented yet", StatusCode: 501, Code: "error-utxo-change-strategy-not-implemented"}

// ErrUtxoNotReserved is when the utxo is not reserved, but a transaction tries to spend it
var ErrUtxoNotReserved = models.SPVError{Message: "transaction utxo has not been reserved for spending", StatusCode: 400, Code: "error-utxo-not-reserved"}

// ////////////////////////////////// XPUB ERRORS

// ErrCouldNotFindXpub is when could not find xpub
var ErrCouldNotFindXpub = models.SPVError{Message: "xpub not found", StatusCode: 404, Code: "error-xpub-not-found"}

// ErrXpubInvalidLength is when the length of the xpub does not match the desired length
var ErrXpubInvalidLength = models.SPVError{Message: "xpub is an invalid length", StatusCode: 400, Code: "error-xpub-length-invalid"}

// ErrXpubNoMatch is when the derived xpub key does not match the key given
var ErrXpubNoMatch = models.SPVError{Message: "xpub key does not match raw key", StatusCode: 400, Code: "error-xpub-key-no-match"}

// ErrXpubIDMisMatch is when the xPubID does not match
var ErrXpubIDMisMatch = models.SPVError{Message: "xpub_id mismatch", StatusCode: 400, Code: "error-xpub-id-mismatch"}

// ////////////////////////////////// MISSING FIELDS

// ErrXPubAlreadyExists is when xpub already exists
var ErrXPubAlreadyExists = models.SPVError{Message: "xpub already exists", StatusCode: 409, Code: "error-xpub-already-exists"}

// ErrOneOfTheFieldsIsRequired is when all of required fields are missing
var ErrOneOfTheFieldsIsRequired = models.SPVError{Message: "missing all of the fields, one of them is required", StatusCode: 400, Code: "error-missing-field-all-required"}

// ErrMissingAccessKey is when the access key field is required but missing
var ErrMissingAccessKey = models.SPVError{Message: "missing required field: access key", StatusCode: 400, Code: "error-missing-field-access-key"}

// ErrMissingFieldID is when the id field is required but missing
var ErrMissingFieldID = models.SPVError{Message: "missing required field: id", StatusCode: 400, Code: "error-missing-field-id"}

// ErrMissingFieldXpubID is when the xpub_id field is required but missing
var ErrMissingFieldXpubID = models.SPVError{Message: "missing required field: xpub_id", StatusCode: 400, Code: "error-missing-field-xpub-id"}

// ErrMissingFieldXpub is when the xpub field is required but missing
var ErrMissingFieldXpub = models.SPVError{Message: "missing required field: xpub", StatusCode: 400, Code: "error-missing-field-xpub"}

// ErrMissingAddress is when the address field address is required but missing
var ErrMissingAddress = models.SPVError{Message: "missing required field: address", StatusCode: 400, Code: "error-missing-field-address"}

// ErrMissingFieldScriptPubKey is when the field is required but missing
var ErrMissingFieldScriptPubKey = models.SPVError{Message: "missing required field: script_pub_key", StatusCode: 400, Code: "error-missing-field-script-pub-key"}

// ErrMissingFieldSatoshis is when the field satoshis is required but missing
var ErrMissingFieldSatoshis = models.SPVError{Message: "missing required field: satoshis", StatusCode: 400, Code: "error-missing-field-satoshis"}

// ErrMissingFieldTransactionID is when the field transaction id is required but missing
var ErrMissingFieldTransactionID = models.SPVError{Message: "missing required field: transaction_id", StatusCode: 400, Code: "error-missing-field-transaction-id"}

// ErrMissingLockingScript is when the field locking script is required but missing
var ErrMissingLockingScript = models.SPVError{Message: "missing required field: locking script", StatusCode: 400, Code: "error-missing-field-locking-script"}

// ////////////////////////////////// SAVE ERROR

// ErrMissingClient is when client is missing from model, cannot save
var ErrMissingClient = models.SPVError{Message: "client is missing from model, cannot save", StatusCode: 400, Code: "error-client-missing"}

// ErrDatastoreRequired is when a datastore function is called without a datastore present
var ErrDatastoreRequired = models.SPVError{Message: "datastore is required", StatusCode: 500, Code: "error-datastore-required"}

// ////////////////////////////////// NOTIFICATION ERRORS

// ErrWebhookSubscriptionFailed is when webhook subscription failed
var ErrWebhookSubscriptionFailed = models.SPVError{Message: "webhook subscription failed", StatusCode: 500, Code: "error-webhook-subscription-failed"}

// ErrWebhookUnsubscriptionFailed is when webhook unsubscription failed
var ErrWebhookUnsubscriptionFailed = models.SPVError{Message: "webhook unsubscription failed", StatusCode: 500, Code: "error-webhook-unsubscription-failed"}

// ErrWebhookSubscriptionNotFound is when cannot find webhook to unsubscribe
var ErrWebhookSubscriptionNotFound = models.SPVError{Message: "webhook subscription not found", StatusCode: 404, Code: "error-webhook-subscription-not-found"}

// ErrWebhookGetAll is when cannot get all the stored webhooks
var ErrWebhookGetAll = models.SPVError{Message: "cannot get all the stored webhooks", StatusCode: 500, Code: "error-webhook-get-all"}

// ErrNotificationsDisabled happens when the notifications are not enabled in the config
var ErrNotificationsDisabled = models.SPVError{Message: "notifications are disabled", StatusCode: 404, Code: "error-notifications-disabled"}

// ////////////////////////////////// ROUTES ERRORS

// ErrRouteNotFound is when route is not found
var ErrRouteNotFound = models.SPVError{Message: "route not found", StatusCode: 404, Code: "error-route-not-found"}

// ErrRouteMethodNotAllowed is when route method is not allowed
var ErrRouteMethodNotAllowed = models.SPVError{Message: "method not allowed", StatusCode: 405, Code: "error-route-method-not-allowed"}

// ////////////////////////////////// BROADCAST ERRORS

// ErrBroadcastUnreachable is when broadcast server cannot be requested
var ErrBroadcastUnreachable = models.SPVError{Message: "broadcast server cannot be requested", StatusCode: 404, Code: "error-broadcast-unreachable"}

// ErrBroadcastWrongBUMPResponse is when broadcast server returned wrong BUMP response
var ErrBroadcastWrongBUMPResponse = models.SPVError{Message: "broadcast server returned wrong BUMP response", StatusCode: 400, Code: "error-broadcast-wrong-bump-response"}

// ErrBroadcastRejectedTransaction is when broadcast server rejected transaction
var ErrBroadcastRejectedTransaction = models.SPVError{Message: "broadcast rejected transaction", StatusCode: 400, Code: "error-broadcast-rejected-transaction"}

// ////////////////////////////////// BLOCK HEADER SERVICE ERRORS

// ErrBHSUnreachable is when Block Header Service (BHS) doesn't respond to status check
var ErrBHSUnreachable = models.SPVError{Message: "Block Header Service cannot be requested", StatusCode: 503, Code: "error-bhs-unreachable"}

// ErrBHSBadURL is when creation of Block Header Service URL doesn't succeed. Probably a fault with the config file
var ErrBHSBadURL = models.SPVError{Message: "cannot create Block Header Service url. Please check your configuration", StatusCode: 500, Code: "error-bhs-bad-url"}

// ErrBHSParsingResponse is when creation of Block Header Service URL doesn't succeed. Probably a fault with the config file
var ErrBHSParsingResponse = models.SPVError{Message: "cannot parse Block Header Service response", StatusCode: 500, Code: "error-bhs-parse-response"}

// ErrBHSInvalidBatchSize is when Block Header Service request contains incorrect batch size query param
var ErrBHSInvalidBatchSize = models.SPVError{Message: "batchSize must be 0 or a positive integer", StatusCode: 400, Code: "error-bhs-invalid-batch-size"}

// ErrBHSMerkleRootNotFound is when Block Header Service cannot find requested merkleroot in lastEvaluatedKey query param
var ErrBHSMerkleRootNotFound = models.SPVError{Message: "No block with provided merkleroot was found", StatusCode: 404, Code: "error-bhs-merkleroot-not-found"}

// ErrBHSMerkleRootNotINLC is when Block Header Service finds merkleroot in lastEvaluateKey query param but it is not in longest chain
var ErrBHSMerkleRootNotInLC = models.SPVError{Message: "Provided merkleroot is not part of the longest chain", StatusCode: 409, Code: "error-bhs-merkleroot-not-part-of-lc"}

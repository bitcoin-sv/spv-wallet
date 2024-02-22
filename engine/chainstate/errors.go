package chainstate

import "errors"

// ErrInvalidTransactionID is when the transaction id is missing or invalid
var ErrInvalidTransactionID = errors.New("invalid transaction id")

// ErrInvalidTransactionHex is when the transaction hex is missing or invalid
var ErrInvalidTransactionHex = errors.New("invalid transaction hex")

// ErrTransactionIDMismatch is when the returned tx does not match the expected given tx id
var ErrTransactionIDMismatch = errors.New("result tx id did not match provided tx id")

// ErrTransactionNotFound is when a transaction was not found in any on-chain provider
var ErrTransactionNotFound = errors.New("transaction not found using all chain providers")

// ErrInvalidRequirements is when an invalid requirement was given
var ErrInvalidRequirements = errors.New("requirements are invalid or missing")

// ErrMissingBroadcastMiners is when broadcasting miners are missing
var ErrMissingBroadcastMiners = errors.New("missing: broadcasting miners")

// ErrMissingQueryMiners is when query miners are missing
var ErrMissingQueryMiners = errors.New("missing: query miners")

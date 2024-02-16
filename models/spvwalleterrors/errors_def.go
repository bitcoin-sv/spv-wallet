// Package spvwalleterrors contains errors that can be returned by spv-wallet api
package spvwalleterrors

import "errors"

// ErrDraftNotFound is when the requested draft transaction was not found
var ErrDraftNotFound = errors.New("corresponding draft transaction not found")

// ErrMissingXPriv is when the xPriv is missing
var ErrMissingXPriv = errors.New("missing xPriv key")

// ErrMissingAccessKey is when the access key is missing
var ErrMissingAccessKey = errors.New("missing access key")

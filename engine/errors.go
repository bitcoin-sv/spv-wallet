package engine

import (
	"errors"
)

// ErrMissingFieldHex is an error when missing the hex field of a transaction
var ErrMissingFieldHex = errors.New("missing required field: hex")

// ErrNoMatchingOutputs is when the transaction does not match any known destinations
var ErrNoMatchingOutputs = errors.New("transaction outputs do not match any known destinations")

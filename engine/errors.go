package engine

import (
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
)

// ErrMissingFieldHex is an error when missing the hex field of a transaction
var ErrMissingFieldHex = spverrors.Newf("missing required field: hex")

// ErrNoMatchingOutputs is when the transaction does not match any known destinations
var ErrNoMatchingOutputs = spverrors.Newf("transaction outputs do not match any known destinations")

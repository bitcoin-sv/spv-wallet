package utils

import (
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
)

// ErrHDKeyNil is when the HD Key is nil
var ErrHDKeyNil = spverrors.Newf("hd key is nil")

// ErrDeriveFailed is when the address derivation failed
var ErrDeriveFailed = spverrors.Newf("derive addresses failed, missing addresses")

// ErrCouldNotDetermineDestinationOutput error when token output could not be determined
var ErrCouldNotDetermineDestinationOutput = spverrors.Newf("could not determine token output destination")

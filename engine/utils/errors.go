package utils

import (
	"errors"
)

// ErrHDKeyNil is when the HD Key is nil
var ErrHDKeyNil = errors.New("hd key is nil")

// ErrDeriveFailed is when the address derivation failed
var ErrDeriveFailed = errors.New("derive addresses failed, missing addresses")

// ErrCouldNotDetermineDestinationOutput error when token output could not be determined
var ErrCouldNotDetermineDestinationOutput = errors.New("could not determine token output destination")

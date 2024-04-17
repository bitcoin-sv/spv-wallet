package utils

import (
	"errors"
)

// ErrXpubInvalidLength is when the length of the xpub does not match the desired length
var ErrXpubInvalidLength = errors.New("xpub is an invalid length")

// ErrXpubNoMatch is when the derived xpub key does not match the key given
var ErrXpubNoMatch = errors.New("xpub key does not match raw key")

// ErrHDKeyNil is when the HD Key is nil
var ErrHDKeyNil = errors.New("hd key is nil")

// ErrDeriveFailed is when the address derivation failed
var ErrDeriveFailed = errors.New("derive addresses failed, missing addresses")

// ErrCouldNotDetermineDestinationOutput error when token output could not be determined
var ErrCouldNotDetermineDestinationOutput = errors.New("could not determine token output destination")

package errors

import "github.com/bitcoin-sv/spv-wallet/models"

// ErrType42DerivationFailed is returned when a type42 public key cannot be derived
var ErrType42DerivationFailed = models.SPVError{
	Code:       "error-custom-instructions-derivation-failed",
	Message:    "Failed to derive type42 public key for given instruction",
	StatusCode: 500,
}

// ErrUnknownInstructionType is returned when an unknown instruction type is encountered
var ErrUnknownInstructionType = models.SPVError{
	Code:       "error-unknown-instruction-type",
	Message:    "Unknown instruction type",
	StatusCode: 500,
}

// ErrGettingAddressFromPublicKey is returned when an address cannot be derived from a public key
var ErrGettingAddressFromPublicKey = models.SPVError{
	Code:       "error-getting-address-from-public-key",
	Message:    "Failed to get address from public key",
	StatusCode: 500,
}

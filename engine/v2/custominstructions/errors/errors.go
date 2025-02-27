package errors

import "github.com/bitcoin-sv/spv-wallet/models"

// ErrProcessingCustomInstructions is returned when custom instructions cannot be processed
var ErrProcessingCustomInstructions = models.SPVError{
	Code:       "error-custom-instructions-processing",
	Message:    "Failed to process custom instructions",
	StatusCode: 422,
}

// ErrFinalizingCustomInstructions is returned when custom instructions cannot be finalized
var ErrFinalizingCustomInstructions = models.SPVError{
	Code:       "error-custom-instructions-finalize",
	Message:    "Failed to finalize processing custom instructions",
	StatusCode: 500,
}

// ErrType42DerivationFailed is returned when a type42 public key cannot be derived
var ErrType42DerivationFailed = models.SPVError{
	Code:       "error-custom-instructions-derivation-failed",
	Message:    "Failed to derive type42 public key for given instruction",
	StatusCode: 422,
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
	StatusCode: 422,
}

// ErrGettingLockingScript is returned when a locking script cannot be derived from an address
var ErrGettingLockingScript = models.SPVError{
	Code:       "error-getting-locking-script",
	Message:    "Failed to get locking script from address",
	StatusCode: 422,
}

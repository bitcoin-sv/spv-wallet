package errors

import "github.com/bitcoin-sv/spv-wallet/models"

var ErrType42DerivationFailed = models.SPVError{
	Code:       "error-custom-instructions-derivation-failed",
	Message:    "Failed to derive type42 public key for given instruction",
	StatusCode: 500,
}

var ErrUnknownInstructionType = models.SPVError{
	Code:       "error-unknown-instruction-type",
	Message:    "Unknown instruction type",
	StatusCode: 500,
}

var ErrGettingAddressFromPublicKey = models.SPVError{
	Code:       "error-getting-address-from-public-key",
	Message:    "Failed to get address from public key",
	StatusCode: 500,
}

var ErrCreatingLockingScript = models.SPVError{
	Code:       "error-creating-locking-script",
	Message:    "Failed to create locking script",
	StatusCode: 500,
}

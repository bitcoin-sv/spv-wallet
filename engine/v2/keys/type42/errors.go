package type42

import "github.com/bitcoin-sv/spv-wallet/models"

// ErrDeriveKey is an error that occurs when a child key cannot be derived from a public key.
var ErrDeriveKey = models.SPVError{Message: "Failed to derive a child key for provided public key", StatusCode: 500, Code: "error-derive-key"}

// ErrRandomReferenceID is an error that occurs when a random reference ID cannot be generated.
var ErrRandomReferenceID = models.SPVError{Message: "Failed to generate a random reference ID", StatusCode: 500, Code: "error-random-reference-id"}

package errors

import "github.com/bitcoin-sv/spv-wallet/models"

// ErrInvalidPublicKey is returned when the public key, provided by the requester, is invalid
var ErrInvalidPublicKey = models.SPVError{Message: "invalid requester public key", StatusCode: 400, Code: "error-user-invalid-requester-pubkey"}

// ErrCreatingUser is returned when the user creation fails
var ErrCreatingUser = models.SPVError{Message: "error creating user", StatusCode: 500, Code: "error-user-creating"}

// ErrInvalidPaymail is returned when the paymail is invalid
var ErrInvalidPaymail = models.SPVError{Message: "invalid paymail", StatusCode: 400, Code: "error-user-invalid-paymail"}

// ErrAddingPaymail is returned when the paymail addition fails
var ErrAddingPaymail = models.SPVError{Message: "error adding paymail", StatusCode: 500, Code: "error-user-adding-paymail"}

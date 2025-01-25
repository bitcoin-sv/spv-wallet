package paymailerrors

import "github.com/bitcoin-sv/spv-wallet/models"

// ErrNoDefaultPaymailAddress is when the user has no default paymail - it actually means that the user has no paymail addresses at all.
var ErrNoDefaultPaymailAddress = models.SPVError{Message: "no default paymail address for user", StatusCode: 400, Code: "error-no-default-paymail-address"}

// ErrInvalidPaymailAddress is when the paymail address is invalid.
var ErrInvalidPaymailAddress = models.SPVError{Message: "invalid paymail address", StatusCode: 400, Code: "error-invalid-paymail-address"}

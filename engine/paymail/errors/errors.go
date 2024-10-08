package pmerrors

import "github.com/bitcoin-sv/spv-wallet/models"

// ErrPaymailHostResponseError is when the paymail host is responding with errors.
var ErrPaymailHostResponseError = models.SPVError{Message: "paymail host is responding with error", StatusCode: 500, Code: "error-paymail-host-error"}

// ErrPaymailHostNotSupportingP2P is when the paymail host is not supporting P2P capabilities.
var ErrPaymailHostNotSupportingP2P = models.SPVError{Message: "paymail host is not supporting P2P capabilities", StatusCode: 400, Code: "error-paymail-host-not-supporting-p2p"}

// ErrPaymailHostInvalidResponse is when the paymail host is responding with invalid response.
var ErrPaymailHostInvalidResponse = models.SPVError{Message: "paymail host invalid response", StatusCode: 500, Code: "error-paymail-host-invalid-response"}

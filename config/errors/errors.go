package errors

import "github.com/bitcoin-sv/spv-wallet/engine/spverrors"

// ErrUnsupportedDomain is returned when the domain is not supported
var ErrUnsupportedDomain = spverrors.Newf("unsupported domain")

// ErrPaymailNotConfigured is returned when the paymail is not configured
var ErrPaymailNotConfigured = spverrors.Newf("paymail not configured")

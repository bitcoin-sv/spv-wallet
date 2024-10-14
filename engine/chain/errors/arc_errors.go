package chainerrors

import "github.com/bitcoin-sv/spv-wallet/models"

// ErrARCUnreachable is when ARC cannot be requested
var ErrARCUnreachable = models.SPVError{Message: "ARC cannot be requested", StatusCode: 500, Code: "error-arc-unreachable"}

// ErrARCUnauthorized is when ARC returns unauthorized
var ErrARCUnauthorized = models.SPVError{Message: "ARC returned unauthorized", StatusCode: 500, Code: "error-arc-unauthorized"}

// ErrARCGenericError is when ARC returns generic error (according to documentation - status code: 409)
var ErrARCGenericError = models.SPVError{Message: "ARC returned generic error", StatusCode: 500, Code: "error-arc-generic-error"}

// ErrARCUnsupportedStatusCode is when ARC returns unsupported status code
var ErrARCUnsupportedStatusCode = models.SPVError{Message: "ARC returned unsupported status code", StatusCode: 500, Code: "error-arc-unsupported-status-code"}

// ErrARCUnprocessable is when ARC rejects because provided tx cannot be processed
var ErrARCUnprocessable = models.SPVError{Message: "ARC cannot process provided transaction", StatusCode: 500, Code: "error-arc-unprocessable-tx"}

// ErrARCNotExtendedFormat is when ARC rejects transaction which is not in extended format
var ErrARCNotExtendedFormat = models.SPVError{Message: "ARC expects transaction in extended format", StatusCode: 500, Code: "error-arc-not-extended-format"}

// ErrARCWrongFee is when ARC rejects transaction because of wrong fee
var ErrARCWrongFee = models.SPVError{Message: "ARC rejected transaction because of wrong fee", StatusCode: 500, Code: "error-arc-wrong-fee"}

// ErrARCProblematicStatus is when ARC returns problematic status
var ErrARCProblematicStatus = models.SPVError{Message: "ARC returned problematic status", StatusCode: 500, Code: "error-arc-problematic-status"}

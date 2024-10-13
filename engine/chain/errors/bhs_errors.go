package chainerrors

import "github.com/bitcoin-sv/spv-wallet/models"

// ErrBHSUnreachable is when Block Header Service (BHS) doesn't respond to status check
var ErrBHSUnreachable = models.SPVError{Message: "Block Header Service cannot be requested", StatusCode: 500, Code: "error-bhs-unreachable"}

// ErrBHSNoSuccessResponse is when Block Header Service request doesn't return a success response
var ErrBHSNoSuccessResponse = models.SPVError{Message: "Block Header Service request did not return a success response", StatusCode: 500, Code: "error-bhs-no-success-response"}

// ErrBHSUnauthorized is when BHS returns unauthorized
var ErrBHSUnauthorized = models.SPVError{Message: "BHS returned unauthorized", StatusCode: 500, Code: "error-bhs-unauthorized"}

// ErrBHSBadRequest is when BHS returns bad request
var ErrBHSBadRequest = models.SPVError{Message: "BHS bad request", StatusCode: 500, Code: "error-bhs-bad-request"}

// ErrBHSUnhealthy is when BHS Healthcheck fails
var ErrBHSUnhealthy = models.SPVError{Message: "Block Header Service is unhealthy", StatusCode: 500, Code: "error-bhs-unhealthy"}

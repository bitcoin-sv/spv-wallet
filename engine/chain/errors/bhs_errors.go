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

// ErrBHSBadURL is when creation of Block Header Service URL doesn't succeed. Probably a fault with the config file
var ErrBHSBadURL = models.SPVError{Message: "cannot create Block Header Service url. Please check your configuration", StatusCode: 500, Code: "error-bhs-bad-url"}

// ErrBHSParsingResponse is when creation of Block Header Service URL doesn't succeed. Probably a fault with the config file
var ErrBHSParsingResponse = models.SPVError{Message: "cannot parse Block Header Service response", StatusCode: 500, Code: "error-bhs-parse-response"}

// ErrInvalidBatchSize is when Block Header Service request contains incorrect batch size query param
var ErrInvalidBatchSize = models.SPVError{Message: "batchSize must be 0 or a positive integer", StatusCode: 400, Code: "error-invalid-batch-size"}

// ErrMerkleRootNotFound is when Block Header Service cannot find requested merkleroot in lastEvaluatedKey query param
var ErrMerkleRootNotFound = models.SPVError{Message: "No block with provided merkleroot was found", StatusCode: 404, Code: "error-merkleroot-not-found"}

// ErrMerkleRootNotInLongestChain is when Block Header Service finds merkleroot in lastEvaluateKey query param but it is not in longest chain
var ErrMerkleRootNotInLongestChain = models.SPVError{Message: "Provided merkleroot is not part of the longest chain", StatusCode: 409, Code: "error-merkleroot-not-part-of-longest-chain"}

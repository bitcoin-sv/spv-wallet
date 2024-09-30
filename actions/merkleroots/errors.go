package merkleroots

import "github.com/bitcoin-sv/spv-wallet/models"

// ////////////////////////////////// BLOCK HEADER SERVICE ERRORS

// ErrBHSUnreachable is when Block Header Service (BHS) doesn't respond to status check
var ErrBHSUnreachable = models.SPVError{Message: "Block Header Service cannot be requested", StatusCode: 500, Code: "error-bhs-unreachable"}

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

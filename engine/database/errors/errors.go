package dberrors

import "github.com/bitcoin-sv/spv-wallet/models"

// ErrDBFailed is when the database operation failed.
var ErrDBFailed = models.SPVError{Message: "database operation failed", StatusCode: 500, Code: "error-db-failed"}

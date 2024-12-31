package txtracker

import "github.com/bitcoin-sv/spv-wallet/models"

// ErrCannotCheckMissingTransactions is when database cannot check for missing transactions
var ErrCannotCheckMissingTransactions = models.SPVError{Message: "Cannot check for missing transactions", StatusCode: 500, Code: "error-check-missing-txs"}

// ErrCannotSaveTransactions is when database cannot save transactions
var ErrCannotSaveTransactions = models.SPVError{Message: "Cannot save transactions", StatusCode: 500, Code: "error-save-txs"}

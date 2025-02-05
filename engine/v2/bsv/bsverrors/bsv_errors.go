package bsverrors

import (
	"github.com/bitcoin-sv/spv-wallet/models"
)

// ErrUnknownTransactionFormat is returned when an unknown transaction format is provided
var ErrUnknownTransactionFormat = models.SPVError{Message: "unknown transaction format provided", StatusCode: 400, Code: "error-unknown-transaction-format"}

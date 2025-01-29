package bsverrors

import (
	"github.com/bitcoin-sv/spv-wallet/models"
)

var ErrUnknownTransactionFormat = models.SPVError{Message: "unknown transaction format provided", StatusCode: 400, Code: "error-unknown-transaction-format"}

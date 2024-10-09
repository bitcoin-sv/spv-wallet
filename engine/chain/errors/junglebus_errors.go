package chainerrors

import "github.com/bitcoin-sv/spv-wallet/models"

// ErrJunglebusFailure is when we can't get transaction from junglebus
var ErrJunglebusFailure = models.SPVError{Message: "junglebus failed to return transaction", StatusCode: 500, Code: "error-junglebus-failure"}

// ErrJunglebusParseTransaction is when we can't parse transaction from junglebus response
var ErrJunglebusParseTransaction = models.SPVError{Message: "failed to parse transaction from junglebus response", StatusCode: 500, Code: "error-junglebus-parse-transaction"}

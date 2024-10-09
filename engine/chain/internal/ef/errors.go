package ef

import "github.com/bitcoin-sv/spv-wallet/models"

var ErrMissingSourceTXID = models.SPVError{Message: "missing source txid", StatusCode: 400, Code: "error-ef-converter-missing-source-txid"}

var ErrGetTransactions = models.SPVError{Message: "error getting transactions", StatusCode: 500, Code: "error-ef-converter-get-transactions"}

var ErrEFHexGeneration = models.SPVError{Message: "error generating ef hex", StatusCode: 500, Code: "error-ef-converter-hex-generation"}

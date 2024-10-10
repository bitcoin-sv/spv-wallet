package ef

import "github.com/bitcoin-sv/spv-wallet/models"

// ErrMissingSourceTXID is returned when SourceTXID field from go-sdk's TransactionInput is nil
var ErrMissingSourceTXID = models.SPVError{Message: "missing source txid", StatusCode: 400, Code: "error-ef-converter-missing-source-txid"}

// ErrGetTransactions is returned when TransactionsGetter fails to get requested transactions
var ErrGetTransactions = models.SPVError{Message: "error getting transactions", StatusCode: 500, Code: "error-ef-converter-get-transactions"}

// ErrEFHexGeneration is returned when EFHex generation fails
var ErrEFHexGeneration = models.SPVError{Message: "error generating ef hex", StatusCode: 500, Code: "error-ef-converter-hex-generation"}

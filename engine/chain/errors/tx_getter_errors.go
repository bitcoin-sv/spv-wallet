package chainerrors

import "github.com/bitcoin-sv/spv-wallet/models"

// ErrGetTransactionsByTxsGetter is when error occurred during getting transactions
var ErrGetTransactionsByTxsGetter = models.SPVError{Message: "error getting transactions during collecting transactions for Txs getter", StatusCode: 500, Code: "error-get-transactions-txs-getter"}

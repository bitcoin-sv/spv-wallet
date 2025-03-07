package usererrors

import "github.com/bitcoin-sv/spv-wallet/models"

// ErrUserHasUnspentUTXOs is when an attempt to delete a user with existing UTXOs occurred
var ErrUserHasUnspentUTXOs = models.SPVError{Message: "cannot delete user with existing UTXOs", StatusCode: 400, Code: "error-user-has-existing-utxos"}

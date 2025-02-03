package testabilities

import (
	"testing"

	"github.com/bitcoin-sv/go-sdk/spv"
	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/stretchr/testify/require"
)

// AssertTxBEEFHex verifies that the given transaction can be successfully encoded into BEEF hex format.
// It ensures that the transaction is not nil, the encoding does not produce an error, and the resulting hex string is non-empty.
func AssertTxBEEFHex(t *testing.T, subjectTx *sdk.Transaction) {
	require.NotNil(t, subjectTx, "subjectTx must not be nil")

	hex, err := subjectTx.BEEFHex()
	require.NoError(t, err, "failed to generate BEEF hex encoding")
	require.NotEmpty(t, hex, "BEEF hex encoding must not be empty")
}

// AssertTxInputs verifies that the provided transaction inputs are valid by checking their SPV script verification.
// It ensures that each input is not nil, its source transaction exists, and script verification passes without errors.
func AssertTxInputs(t *testing.T, inputs ...*sdk.TransactionInput) {
	require.NotNil(t, inputs, "inputs must not be nil")
	for _, input := range inputs {
		ok, err := spv.VerifyScripts(input.SourceTransaction)
		require.True(t, ok, "SPV script verification failed")
		require.NoError(t, err, "unexpected error during SPV script verification")
	}
}

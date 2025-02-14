package outlines

import (
	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"testing"
)

const estimatedInputSizeForP2PKH = 148

// TestEstimateSizeAgainstRealTxSize verifies that our internal transaction size estimation function
// remains in sync with the actual signed transaction sizes produced by the SDK.
//
// Since transaction fees depend on the byte-size of the transaction,
// it is essential that our estimation closely approximates the real size.
//
// This is a test for internal function `estimateSize` which we try to avoid. However, because accurate fee calculation
// depends critically on transaction size, it is essential to directly verify the behavior of this internal estimation logic.
//
// Note: A slight overestimation (up to 1 byte per input) is tolerated to accommodate minor differences
// arising from fee calculations and edge cases in the SDK.
func TestEstimateSizeAgainstRealTxSie(t *testing.T) {
	tests := []*sdk.Transaction{
		fixtures.GivenTX(t).
			WithInput(10).
			WithP2PKHOutput(5).
			WithOPReturn("test").
			TX(),
		fixtures.GivenTX(t).
			WithInput(10).
			WithInput(15).
			WithP2PKHOutput(5).
			WithP2PKHOutput(10).
			WithOPReturn("test").
			WithOPReturn("foo").
			TX(),
		fixtures.GivenTX(t).
			WithInput(10).
			WithInput(15).
			WithInput(20).
			WithP2PKHOutput(5).
			WithP2PKHOutput(10).
			WithOPReturn("test").
			WithOPReturn("foo").
			TX(),
	}
	for _, tx := range tests {
		t.Run("", func(t *testing.T) {
			// given:

			aIn := lo.Map(tx.Inputs, func(input *sdk.TransactionInput, _ int) *annotatedInput {
				return &annotatedInput{
					TransactionInput: input,
					estimatedSize:    estimatedInputSizeForP2PKH,
					utxoSatoshis:     bsv.Satoshis(*input.SourceTxSatoshis()),
				}
			})

			aOut := lo.Map(tx.Outputs, func(output *sdk.TransactionOutput, _ int) *annotatedOutput {
				return &annotatedOutput{
					TransactionOutput: output,
				}
			})

			// when:
			realSignedTxSize := tx.Size()
			internalSizeAlg := int(estimatedSize(aIn, aOut))

			// then:
			if realSignedTxSize != internalSizeAlg {
				if internalSizeAlg >= realSignedTxSize && internalSizeAlg <= realSignedTxSize+len(tx.Inputs)*1 {
					t.Logf("Internal estimation alg returned size %d, more than real tx size %d", internalSizeAlg, realSignedTxSize)
				} else {
					require.Failf(t, "size mismatch", "size mismatch. SDK: %d, internal: %d", realSignedTxSize, internalSizeAlg)
				}
			}
		})
	}
}

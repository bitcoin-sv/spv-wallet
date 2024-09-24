package draft_test

import (
	"context"
	"encoding/hex"
	"strings"
	"testing"

	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction/draft"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction/draft/outputs"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction/draft/testabilities"
	txerrors "github.com/bitcoin-sv/spv-wallet/engine/transaction/errors"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/request/opreturn"
	"github.com/stretchr/testify/require"
)

func TestCreateOpReturnDraft(t *testing.T) {
	// maxOpPushDataSize is the maximum size of the data chunk that can be pushed to transaction with OP_PUSH operation.
	const maxOpPushDataSize = 0xFFFFFFFF

	successTests := map[string]struct {
		opReturn      *outputs.OpReturn
		lockingScript string
	}{
		"for single string": {
			opReturn: &outputs.OpReturn{
				DataType: opreturn.DataTypeStrings,
				Data:     []string{"Example data"},
			},
			lockingScript: "006a0c4578616d706c652064617461",
		},
		"for multiple strings": {
			opReturn: &outputs.OpReturn{
				DataType: opreturn.DataTypeStrings,
				Data:     []string{"Example", " ", "data"},
			},
			lockingScript: "006a074578616d706c6501200464617461",
		},
		"for single hex": {
			opReturn: &outputs.OpReturn{
				DataType: opreturn.DataTypeHexes,
				Data:     []string{toHex("Example data")},
			},
			lockingScript: "006a0c4578616d706c652064617461",
		},
		"for multiple hexes": {
			opReturn: &outputs.OpReturn{
				DataType: opreturn.DataTypeHexes,
				Data:     []string{toHex("Example"), toHex(" "), toHex("data")},
			},
			lockingScript: "006a074578616d706c6501200464617461",
		},
	}
	for name, test := range successTests {
		t.Run("return draft "+name, func(t *testing.T) {
			given := testabilities.New(t)

			// given:
			draftService := given.NewDraftTransactionService()

			// and:
			spec := &draft.TransactionSpec{
				XPubID:  fixtures.Sender.XPubID,
				Outputs: outputs.NewSpecifications(test.opReturn),
			}

			// when:
			draftTx, err := draftService.Create(context.Background(), spec)

			// then:
			require.NoError(t, err)
			require.NotNil(t, draftTx)

			// and:
			annotations := draftTx.Annotations
			require.Len(t, annotations.Outputs, 1)
			require.Equal(t, transaction.BucketData, annotations.Outputs[0].Bucket)

			// debug:
			t.Logf("BEEF: %s", draftTx.BEEF)

			// when:
			tx, err := sdk.NewTransactionFromBEEFHex(draftTx.BEEF)

			// then:
			require.NoError(t, err)
			require.Len(t, tx.Outputs, 1)
			require.EqualValues(t, tx.Outputs[0].Satoshis, 0)
			require.Equal(t, tx.Outputs[0].LockingScript.IsData(), true)
			require.Equal(t, test.lockingScript, tx.Outputs[0].LockingScriptHex())
		})
	}

	errorTests := map[string]struct {
		spec          *outputs.OpReturn
		expectedError models.SPVError
	}{
		"for no data in default type": {
			spec:          &outputs.OpReturn{},
			expectedError: txerrors.ErrDraftOpReturnDataRequired,
		},
		"for no data string type": {
			spec: &outputs.OpReturn{
				DataType: opreturn.DataTypeStrings,
			},
			expectedError: txerrors.ErrDraftOpReturnDataRequired,
		},
		"for invalid hex": {
			spec: &outputs.OpReturn{
				DataType: opreturn.DataTypeHexes,
				Data:     []string{"invalid hex"},
			},
			expectedError: txerrors.ErrFailedToDecodeHex,
		},
		"for unknown data type": {
			spec: &outputs.OpReturn{
				DataType: 123,
				Data:     []string{"Example", " ", "data"},
			},
			expectedError: txerrors.ErrDraftOpReturnUnsupportedDataType,
		},
		"for to big string": {
			spec: &outputs.OpReturn{
				DataType: opreturn.DataTypeStrings,
				Data:     []string{strings.Repeat("1", maxOpPushDataSize+1)},
			},
			expectedError: txerrors.ErrDraftOpReturnDataTooLarge,
		},
	}
	for name, test := range errorTests {
		t.Run("return error "+name, func(t *testing.T) {
			given := testabilities.New(t)

			// given:
			draftService := given.NewDraftTransactionService()

			// and:
			spec := &draft.TransactionSpec{
				XPubID:  fixtures.Sender.XPubID,
				Outputs: outputs.NewSpecifications(test.spec),
			}

			// when:
			_, err := draftService.Create(context.Background(), spec)

			// then:
			require.Error(t, err)
			require.ErrorIs(t, err, test.expectedError)
		})
	}
}

func toHex(data string) string {
	return hex.EncodeToString([]byte(data))
}

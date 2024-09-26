package draft_test

import (
	"context"
	"encoding/hex"
	"strings"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction/draft"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction/draft/outputs"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction/draft/testabilities"
	txerrors "github.com/bitcoin-sv/spv-wallet/engine/transaction/errors"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/request/opreturn"
)

func TestCreateOpReturnDraft(t *testing.T) {
	// maxOpPushDataSize is the maximum size of the data chunk that can be pushed to transaction with OP_PUSH operation.
	const maxOpPushDataSize = 0xFFFFFFFF

	successTests := map[string]struct {
		opReturn      *outputs.OpReturn
		lockingScript string
	}{
		"return draft for single string": {
			opReturn: &outputs.OpReturn{
				DataType: opreturn.DataTypeStrings,
				Data:     []string{"Example data"},
			},
			lockingScript: "006a0c4578616d706c652064617461",
		},
		"return draft for multiple strings": {
			opReturn: &outputs.OpReturn{
				DataType: opreturn.DataTypeStrings,
				Data:     []string{"Example", " ", "data"},
			},
			lockingScript: "006a074578616d706c6501200464617461",
		},
		"return draft for single hex": {
			opReturn: &outputs.OpReturn{
				DataType: opreturn.DataTypeHexes,
				Data:     []string{toHex("Example data")},
			},
			lockingScript: "006a0c4578616d706c652064617461",
		},
		"return draft for multiple hexes": {
			opReturn: &outputs.OpReturn{
				DataType: opreturn.DataTypeHexes,
				Data:     []string{toHex("Example"), toHex(" "), toHex("data")},
			},
			lockingScript: "006a074578616d706c6501200464617461",
		},
	}
	for name, test := range successTests {
		t.Run(name, func(t *testing.T) {
			given, then := testabilities.New(t)

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
			then.Created(draftTx).WithNoError(err).WithParseableBEEFHex().
				HasOutputs(1).
				Output(0).
				HasBucket(transaction.BucketData).
				IsDataOnly().
				HasLockingScript(test.lockingScript)
		})
	}

	errorTests := map[string]struct {
		spec          *outputs.OpReturn
		expectedError models.SPVError
	}{
		"return error for no data in default type": {
			spec:          &outputs.OpReturn{},
			expectedError: txerrors.ErrDraftOpReturnDataRequired,
		},
		"return error for no data string type": {
			spec: &outputs.OpReturn{
				DataType: opreturn.DataTypeStrings,
			},
			expectedError: txerrors.ErrDraftOpReturnDataRequired,
		},
		"return error for invalid hex": {
			spec: &outputs.OpReturn{
				DataType: opreturn.DataTypeHexes,
				Data:     []string{"invalid hex"},
			},
			expectedError: txerrors.ErrFailedToDecodeHex,
		},
		"return error for unknown data type": {
			spec: &outputs.OpReturn{
				DataType: 123,
				Data:     []string{"Example", " ", "data"},
			},
			expectedError: txerrors.ErrDraftOpReturnUnsupportedDataType,
		},
		"return error for to big string": {
			spec: &outputs.OpReturn{
				DataType: opreturn.DataTypeStrings,
				Data:     []string{strings.Repeat("1", maxOpPushDataSize+1)},
			},
			expectedError: txerrors.ErrDraftOpReturnDataTooLarge,
		},
	}
	for name, test := range errorTests {
		t.Run(name, func(t *testing.T) {
			given, then := testabilities.New(t)

			// given:
			draftService := given.NewDraftTransactionService()

			// and:
			spec := &draft.TransactionSpec{
				XPubID:  fixtures.Sender.XPubID,
				Outputs: outputs.NewSpecifications(test.spec),
			}

			// when:
			tx, err := draftService.Create(context.Background(), spec)

			// then:
			then.Created(tx).WithError(err).ThatIs(test.expectedError)
		})
	}
}

func toHex(data string) string {
	return hex.EncodeToString([]byte(data))
}

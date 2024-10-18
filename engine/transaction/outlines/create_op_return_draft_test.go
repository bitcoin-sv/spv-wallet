package outlines_test

import (
	"context"
	"encoding/hex"
	"strings"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
	txerrors "github.com/bitcoin-sv/spv-wallet/engine/transaction/errors"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction/outlines"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction/outlines/outputs"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction/outlines/testabilities"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/request/opreturn"
	"github.com/bitcoin-sv/spv-wallet/models/transaction/bucket"
)

func TestCreateOpReturnTransactionOutline(t *testing.T) {
	// maxOpPushDataSize is the maximum size of the data chunk that can be pushed to transaction with OP_PUSH operation.
	const maxOpPushDataSize = 0xFFFFFFFF

	successTests := map[string]struct {
		opReturn      *outputs.OpReturn
		lockingScript string
	}{
		"return transaction outline for single string": {
			opReturn: &outputs.OpReturn{
				DataType: opreturn.DataTypeStrings,
				Data:     []string{"Example data"},
			},
			lockingScript: "006a0c4578616d706c652064617461",
		},
		"return transaction outline for multiple strings": {
			opReturn: &outputs.OpReturn{
				DataType: opreturn.DataTypeStrings,
				Data:     []string{"Example", " ", "data"},
			},
			lockingScript: "006a074578616d706c6501200464617461",
		},
		"return transaction outline for single hex": {
			opReturn: &outputs.OpReturn{
				DataType: opreturn.DataTypeHexes,
				Data:     []string{toHex("Example data")},
			},
			lockingScript: "006a0c4578616d706c652064617461",
		},
		"return transaction outline for multiple hexes": {
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
			service := given.NewTransactionOutlinesService()

			// and:
			spec := &outlines.TransactionSpec{
				XPubID:  fixtures.Sender.XPubID(),
				Outputs: outputs.NewSpecifications(test.opReturn),
			}

			// when:
			tx, err := service.Create(context.Background(), spec)

			// then:
			thenTx := then.Created(tx).WithNoError(err).WithParseableBEEFHex()

			thenTx.HasOutputs(1)

			thenTx.Output(0).
				HasBucket(bucket.Data).
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
			expectedError: txerrors.ErrTxOutlineOpReturnDataRequired,
		},
		"return error for no data string type": {
			spec: &outputs.OpReturn{
				DataType: opreturn.DataTypeStrings,
			},
			expectedError: txerrors.ErrTxOutlineOpReturnDataRequired,
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
			expectedError: txerrors.ErrTxOutlineOpReturnUnsupportedDataType,
		},
		"return error for to big string": {
			spec: &outputs.OpReturn{
				DataType: opreturn.DataTypeStrings,
				Data:     []string{strings.Repeat("1", maxOpPushDataSize+1)},
			},
			expectedError: txerrors.ErrTxOutlineOpReturnDataTooLarge,
		},
	}
	for name, test := range errorTests {
		t.Run(name, func(t *testing.T) {
			given, then := testabilities.New(t)

			// given:
			service := given.NewTransactionOutlinesService()

			// and:
			spec := &outlines.TransactionSpec{
				XPubID:  fixtures.Sender.XPubID(),
				Outputs: outputs.NewSpecifications(test.spec),
			}

			// when:
			tx, err := service.Create(context.Background(), spec)

			// then:
			then.Created(tx).WithError(err).ThatIs(test.expectedError)
		})
	}
}

func toHex(data string) string {
	return hex.EncodeToString([]byte(data))
}

package outlines_test

import (
	"context"
	"encoding/hex"
	"strings"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/errors"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/outlines"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/outlines/testabilities"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/request/opreturn"
	"github.com/bitcoin-sv/spv-wallet/models/transaction/bucket"
)

func TestCreateOpReturnTransactionOutlineBEEF(t *testing.T) {
	// maxOpPushDataSize is the maximum size of the data chunk that can be pushed to transaction with OP_PUSH operation.
	const maxOpPushDataSize = 0xFFFFFFFF

	successTests := map[string]struct {
		opReturn      *outlines.OpReturn
		lockingScript string
	}{
		"return transaction outline for single string": {
			opReturn: &outlines.OpReturn{
				DataType: opreturn.DataTypeStrings,
				Data:     []string{"Example data"},
			},
			lockingScript: "006a0c4578616d706c652064617461",
		},
		"return transaction outline for multiple strings": {
			opReturn: &outlines.OpReturn{
				DataType: opreturn.DataTypeStrings,
				Data:     []string{"Example", " ", "data"},
			},
			lockingScript: "006a074578616d706c6501200464617461",
		},
		"return transaction outline for single hex": {
			opReturn: &outlines.OpReturn{
				DataType: opreturn.DataTypeHexes,
				Data:     []string{toHex("Example data")},
			},
			lockingScript: "006a0c4578616d706c652064617461",
		},
		"return transaction outline for multiple hexes": {
			opReturn: &outlines.OpReturn{
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
				UserID:  fixtures.Sender.ID(),
				Outputs: outlines.NewOutputsSpecs(test.opReturn),
			}

			// when:
			tx, err := service.CreateBEEF(context.Background(), spec)

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
		spec          *outlines.OpReturn
		expectedError models.SPVError
	}{
		"return error for no data in default type": {
			spec:          &outlines.OpReturn{},
			expectedError: txerrors.ErrTxOutlineOpReturnDataRequired,
		},
		"return error for no data string type": {
			spec: &outlines.OpReturn{
				DataType: opreturn.DataTypeStrings,
			},
			expectedError: txerrors.ErrTxOutlineOpReturnDataRequired,
		},
		"return error for invalid hex": {
			spec: &outlines.OpReturn{
				DataType: opreturn.DataTypeHexes,
				Data:     []string{"invalid hex"},
			},
			expectedError: txerrors.ErrFailedToDecodeHex,
		},
		"return error for unknown data type": {
			spec: &outlines.OpReturn{
				DataType: 123,
				Data:     []string{"Example", " ", "data"},
			},
			expectedError: txerrors.ErrTxOutlineOpReturnUnsupportedDataType,
		},
		"return error for to big string": {
			spec: &outlines.OpReturn{
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
				UserID:  fixtures.Sender.ID(),
				Outputs: outlines.NewOutputsSpecs(test.spec),
			}

			// when:
			tx, err := service.CreateBEEF(context.Background(), spec)

			// then:
			then.Created(tx).WithError(err).ThatIs(test.expectedError)
		})
	}
}

func TestCreateOpReturnTransactionOutlineRAW(t *testing.T) {
	// maxOpPushDataSize is the maximum size of the data chunk that can be pushed to transaction with OP_PUSH operation.
	const maxOpPushDataSize = 0xFFFFFFFF

	successTests := map[string]struct {
		opReturn      *outlines.OpReturn
		lockingScript string
	}{
		"return transaction outline for default data type (strings)": {
			opReturn: &outlines.OpReturn{
				Data: []string{"Example data"},
			},
			lockingScript: "006a0c4578616d706c652064617461",
		},
		"return transaction outline for single string": {
			opReturn: &outlines.OpReturn{
				DataType: opreturn.DataTypeStrings,
				Data:     []string{"Example data"},
			},
			lockingScript: "006a0c4578616d706c652064617461",
		},
		"return transaction outline for multiple strings": {
			opReturn: &outlines.OpReturn{
				DataType: opreturn.DataTypeStrings,
				Data:     []string{"Example", " ", "data"},
			},
			lockingScript: "006a074578616d706c6501200464617461",
		},
		"return transaction outline for single hex": {
			opReturn: &outlines.OpReturn{
				DataType: opreturn.DataTypeHexes,
				Data:     []string{toHex("Example data")},
			},
			lockingScript: "006a0c4578616d706c652064617461",
		},
		"return transaction outline for multiple hexes": {
			opReturn: &outlines.OpReturn{
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
				UserID:  fixtures.Sender.ID(),
				Outputs: outlines.NewOutputsSpecs(test.opReturn),
			}

			// when:
			tx, err := service.CreateRawTx(context.Background(), spec)

			// then:
			thenTx := then.Created(tx).WithNoError(err).WithParseableRawHex()

			thenTx.HasInputs(1)

			thenTx.Input(0).
				HasOutpoint(testabilities.UserFundsTransactionOutpoint)

			thenTx.HasOutputs(1)

			thenTx.Output(0).
				HasBucket(bucket.Data).
				IsDataOnly().
				HasLockingScript(test.lockingScript)
		})
	}

	errorTests := map[string]struct {
		spec          *outlines.OpReturn
		expectedError models.SPVError
	}{
		"return error for no data in default type": {
			spec:          &outlines.OpReturn{},
			expectedError: txerrors.ErrTxOutlineOpReturnDataRequired,
		},
		"return error for no data string type": {
			spec: &outlines.OpReturn{
				DataType: opreturn.DataTypeStrings,
			},
			expectedError: txerrors.ErrTxOutlineOpReturnDataRequired,
		},
		"return error for invalid hex": {
			spec: &outlines.OpReturn{
				DataType: opreturn.DataTypeHexes,
				Data:     []string{"invalid hex"},
			},
			expectedError: txerrors.ErrFailedToDecodeHex,
		},
		"return error for unknown data type": {
			spec: &outlines.OpReturn{
				DataType: 123,
				Data:     []string{"Example", " ", "data"},
			},
			expectedError: txerrors.ErrTxOutlineOpReturnUnsupportedDataType,
		},
		"return error for to big string": {
			spec: &outlines.OpReturn{
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
				UserID:  fixtures.Sender.ID(),
				Outputs: outlines.NewOutputsSpecs(test.spec),
			}

			// when:
			tx, err := service.CreateBEEF(context.Background(), spec)

			// then:
			then.Created(tx).WithError(err).ThatIs(test.expectedError)
		})
	}
}

func toHex(data string) string {
	return hex.EncodeToString([]byte(data))
}

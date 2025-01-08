package record_test

import (
	"context"
	"testing"

	"github.com/bitcoin-sv/go-sdk/script"
	"github.com/bitcoin-sv/spv-wallet/engine/database"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction"
	txerrors "github.com/bitcoin-sv/spv-wallet/engine/transaction/errors"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction/outlines"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction/record/testabilities"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"github.com/bitcoin-sv/spv-wallet/models/transaction/bucket"
	"github.com/pkg/errors"
)

const (
	dataOfOpReturnTx = "hello world"
	notABeefHex      = "0100000001b23c7c47320b3818c665bf28a46c290f3fb379ea8d357625bbff3117ae14b09b0000000000ffffffff0100000000000000000e006a0b68656c6c6f20776f726c6400000000"
)

func givenTXWithOpReturn(t *testing.T) fixtures.GivenTXSpec {
	return fixtures.GivenTX(t).
		WithInput(1).
		WithOPReturn(dataOfOpReturnTx)
}

func givenTxWithOpReturnWithoutOPFalse(t *testing.T) fixtures.GivenTXSpec {
	return fixtures.GivenTX(t).
		WithInput(1).
		WithOutputScriptParts(
			fixtures.OpCode(script.OpRETURN),
			fixtures.PushData(dataOfOpReturnTx),
		)
}

func givenStandardOpReturnOutline(t *testing.T) *outlines.Transaction {
	return &outlines.Transaction{
		BEEF: givenTXWithOpReturn(t).BEEF(),
		Annotations: transaction.Annotations{
			Outputs: transaction.OutputsAnnotations{
				0: &transaction.OutputAnnotation{
					Bucket: bucket.Data,
				},
			},
		},
	}
}

func TestRecordOutlineOpReturn(t *testing.T) {
	tests := map[string]struct {
		storedUTXOs   []bsv.Outpoint
		outline       *outlines.Transaction
		expectTxID    string
		expectOutputs []database.Output
		expectData    []database.Data
	}{
		"RecordTransactionOutline for op_return": {
			storedUTXOs: []bsv.Outpoint{givenTXWithOpReturn(t).InputUTXO(0)},
			outline:     givenStandardOpReturnOutline(t),
			expectTxID:  givenTXWithOpReturn(t).ID(),
			expectOutputs: []database.Output{
				{
					TxID:       givenTXWithOpReturn(t).InputUTXO(0).TxID,
					Vout:       givenTXWithOpReturn(t).InputUTXO(0).Vout,
					SpendingTX: givenTXWithOpReturn(t).ID(),
				},
				{
					TxID:       givenTXWithOpReturn(t).ID(),
					Vout:       0,
					SpendingTX: "",
				},
			},
			expectData: []database.Data{
				{
					TxID: givenTXWithOpReturn(t).ID(),
					Vout: 0,
					Blob: []byte(dataOfOpReturnTx),
				},
			},
		},
		"RecordTransactionOutline for op_return without leading OP_FALSE": {
			storedUTXOs: []bsv.Outpoint{givenTxWithOpReturnWithoutOPFalse(t).InputUTXO(0)},
			outline: &outlines.Transaction{
				BEEF: givenTxWithOpReturnWithoutOPFalse(t).BEEF(),
				Annotations: transaction.Annotations{
					Outputs: transaction.OutputsAnnotations{
						0: &transaction.OutputAnnotation{
							Bucket: bucket.Data,
						},
					},
				},
			},
			expectTxID: givenTxWithOpReturnWithoutOPFalse(t).ID(),
			expectOutputs: []database.Output{
				{
					TxID:       givenTxWithOpReturnWithoutOPFalse(t).InputUTXO(0).TxID,
					Vout:       givenTxWithOpReturnWithoutOPFalse(t).InputUTXO(0).Vout,
					SpendingTX: givenTxWithOpReturnWithoutOPFalse(t).ID(),
				},
				{
					TxID:       givenTxWithOpReturnWithoutOPFalse(t).ID(),
					Vout:       0,
					SpendingTX: "",
				},
			},
			expectData: []database.Data{
				{
					TxID: givenTxWithOpReturnWithoutOPFalse(t).ID(),
					Vout: 0,
					Blob: []byte(dataOfOpReturnTx),
				},
			},
		},
		"RecordTransactionOutline for op_return with untracked utxo as inputs": {
			storedUTXOs: []bsv.Outpoint{},
			outline:     givenStandardOpReturnOutline(t),
			expectTxID:  givenTXWithOpReturn(t).ID(),
			expectOutputs: []database.Output{{
				TxID: givenTXWithOpReturn(t).ID(),
				Vout: 0,
			}},
			expectData: []database.Data{
				{
					TxID: givenTXWithOpReturn(t).ID(),
					Vout: 0,
					Blob: []byte(dataOfOpReturnTx),
				},
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			given, then := testabilities.New(t)
			given.Repository().WithUTXOs(test.storedUTXOs...)

			service := given.NewRecordService()

			// when:
			err := service.RecordTransactionOutline(context.Background(), test.outline)

			// then:
			then.NoError(err).
				Broadcasted(test.expectTxID).
				StoredAsBroadcasted(test.expectTxID)

			then.
				StoredOutputs(test.expectOutputs).
				StoredData(test.expectData)
		})
	}
}

func TestRecordOutlineOpReturnErrorCases(t *testing.T) {
	givenUnsignedTX := fixtures.GivenTX(t).
		WithoutSigning().
		WithInput(1).
		WithOPReturn(dataOfOpReturnTx)

	givenTxWithOpZeroAfterOpReturn := fixtures.GivenTX(t).
		WithInput(1).
		WithOutputScriptParts(
			fixtures.OpCode(script.OpFALSE),
			fixtures.OpCode(script.OpRETURN),
			fixtures.PushData(dataOfOpReturnTx),
			fixtures.OpCode(script.OpZERO),
			fixtures.OpCode(script.OpZERO),
			fixtures.PushData(dataOfOpReturnTx),
		)

	givenTxWithP2PKHOutput := fixtures.GivenTX(t).
		WithInput(2).
		WithP2PKHOutput(1)

	tests := map[string]struct {
		storedOutputs []database.Output
		outline       *outlines.Transaction
		expectErr     error
	}{
		"RecordTransactionOutline for not signed transaction": {
			storedOutputs: []database.Output{},
			outline: &outlines.Transaction{
				BEEF: givenUnsignedTX.BEEF(),
			},
			expectErr: txerrors.ErrTxValidation,
		},
		"RecordTransactionOutline for not a BEEF hex": {
			storedOutputs: []database.Output{},
			outline: &outlines.Transaction{
				BEEF: notABeefHex,
			},
			expectErr: txerrors.ErrTxValidation,
		},
		"RecordTransactionOutline for invalid OP_ZERO after OP_RETURN": {
			storedOutputs: []database.Output{},
			outline: &outlines.Transaction{
				BEEF: givenTxWithOpZeroAfterOpReturn.BEEF(),
				Annotations: transaction.Annotations{
					Outputs: transaction.OutputsAnnotations{
						0: &transaction.OutputAnnotation{
							Bucket: bucket.Data,
						},
					},
				},
			},
			expectErr: txerrors.ErrOnlyPushDataAllowed,
		},
		"Tx with already spent utxo": {
			storedOutputs: []database.Output{{
				TxID:       givenTXWithOpReturn(t).InputUTXO(0).TxID,
				Vout:       givenTXWithOpReturn(t).InputUTXO(0).Vout,
				SpendingTX: "05aa91319c773db18071310ecd5ddc15d3aa4242b55705a13a66f7fefe2b80a1",
			}},
			outline: &outlines.Transaction{
				BEEF: givenTXWithOpReturn(t).BEEF(),
			},
			expectErr: txerrors.ErrUTXOSpent,
		},
		"Vout out of range in annotation": {
			storedOutputs: []database.Output{},
			outline: &outlines.Transaction{
				BEEF: givenTXWithOpReturn(t).BEEF(),
				Annotations: transaction.Annotations{
					Outputs: transaction.OutputsAnnotations{
						1: &transaction.OutputAnnotation{
							Bucket: bucket.Data,
						},
					},
				},
			},
			expectErr: txerrors.ErrAnnotationIndexOutOfRange,
		},
		"Vout as negative value in annotation": {
			storedOutputs: []database.Output{},
			outline: &outlines.Transaction{
				BEEF: givenTXWithOpReturn(t).BEEF(),
				Annotations: transaction.Annotations{
					Outputs: transaction.OutputsAnnotations{
						-1: &transaction.OutputAnnotation{
							Bucket: bucket.Data,
						},
					},
				},
			},
			expectErr: txerrors.ErrAnnotationIndexConversion,
		},
		"no-op_return output annotated as data": {
			storedOutputs: []database.Output{},
			outline: &outlines.Transaction{
				BEEF: givenTxWithP2PKHOutput.BEEF(),
				Annotations: transaction.Annotations{
					Outputs: transaction.OutputsAnnotations{
						0: &transaction.OutputAnnotation{
							Bucket: bucket.Data,
						},
					},
				},
			},
			expectErr: txerrors.ErrAnnotationMismatch,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			given, then := testabilities.New(t)
			given.Repository().WithOutputs(test.storedOutputs...)

			service := given.NewRecordService()

			// when:
			err := service.RecordTransactionOutline(context.Background(), test.outline)

			// then:
			then.ErrorIs(err, test.expectErr).NothingChanged()
		})
	}
}

func TestOnBroadcastErr(t *testing.T) {
	// given:
	given, then := testabilities.New(t)
	given.Repository().
		WithOutputs(database.Output{
			TxID:       givenTXWithOpReturn(t).InputUTXO(0).TxID,
			Vout:       givenTXWithOpReturn(t).InputUTXO(0).Vout,
			SpendingTX: "",
		})
	given.Broadcaster().
		WillFailOnBroadcast(errors.New("broadcast error"))

	service := given.NewRecordService()

	// and:
	outline := givenStandardOpReturnOutline(t)

	// when:
	err := service.RecordTransactionOutline(context.Background(), outline)

	// then:
	then.ErrorIs(err, txerrors.ErrTxBroadcast).NothingChanged()
}

func TestOnSaveTXErr(t *testing.T) {
	// given:
	given, then := testabilities.New(t)
	given.Repository().
		WillFailOnSaveTX(errors.New("saveTX error"))

	service := given.NewRecordService()

	// and:
	outline := givenStandardOpReturnOutline(t)

	// when:
	err := service.RecordTransactionOutline(context.Background(), outline)

	// then:
	then.ErrorIs(err, txerrors.ErrSavingData).NothingChanged()
}

func TestOnGetOutputsErr(t *testing.T) {
	// given:
	given, then := testabilities.New(t)
	given.Repository().
		WillFailOnGetOutputs(errors.New("getOutputs error"))
	service := given.NewRecordService()

	// and:
	outline := givenStandardOpReturnOutline(t)

	// when:
	err := service.RecordTransactionOutline(context.Background(), outline)

	// then:
	then.ErrorIs(err, txerrors.ErrGettingOutputs).NothingChanged()
}

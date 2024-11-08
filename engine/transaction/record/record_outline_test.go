package record_test

import (
	"context"
	"testing"

	"github.com/bitcoin-sv/go-sdk/script"
	"github.com/bitcoin-sv/spv-wallet/engine/database"
	"github.com/bitcoin-sv/spv-wallet/engine/tester"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction"
	txerrors "github.com/bitcoin-sv/spv-wallet/engine/transaction/errors"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction/outlines"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction/record"
	"github.com/bitcoin-sv/spv-wallet/models/transaction/bucket"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

const (
	dataOfOpReturnTx = "hello world"
	notABeefHex      = "0100000001b23c7c47320b3818c665bf28a46c290f3fb379ea8d357625bbff3117ae14b09b0000000000ffffffff0100000000000000000e006a0b68656c6c6f20776f726c6400000000"
)

func TestRecordOutlineOpReturn(t *testing.T) {
	givenTxWithOpReturn := fixtures.GivenTX(t).
		WithInput(1).
		WithOPReturn(dataOfOpReturnTx)

	givenTxWithOpReturnWithoutOPFalse := fixtures.GivenTX(t).
		WithInput(1).
		WithOutputScript(
			fixtures.OpCode(script.OpRETURN),
			fixtures.PushData(dataOfOpReturnTx),
		)

	tests := map[string]struct {
		repo          *mockRepository
		outline       *outlines.Transaction
		expectTxID    string
		expectOutputs []database.Output
		expectData    []database.Data
	}{
		"RecordTransactionOutline for op_return": {
			repo: newMockRepository().withUTXO(givenTxWithOpReturn.InputUTXO(0)),
			outline: &outlines.Transaction{
				BEEF: givenTxWithOpReturn.BEEF(),
				Annotations: transaction.Annotations{
					Outputs: transaction.OutputsAnnotations{
						0: &transaction.OutputAnnotation{
							Bucket: bucket.Data,
						},
					},
				},
			},
			expectTxID: givenTxWithOpReturn.ID(),
			expectOutputs: []database.Output{
				{
					TxID:       givenTxWithOpReturn.InputUTXO(0).TxID,
					Vout:       givenTxWithOpReturn.InputUTXO(0).Vout,
					SpendingTX: ptr(givenTxWithOpReturn.ID()),
				},
				{
					TxID:       givenTxWithOpReturn.ID(),
					Vout:       0,
					SpendingTX: nil,
				},
			},
			expectData: []database.Data{
				{
					TxID: givenTxWithOpReturn.ID(),
					Vout: 0,
					Blob: []byte(dataOfOpReturnTx),
				},
			},
		},
		"RecordTransactionOutline for op_return without leading OP_FALSE": {
			repo: newMockRepository().withUTXO(givenTxWithOpReturnWithoutOPFalse.InputUTXO(0)),
			outline: &outlines.Transaction{
				BEEF: givenTxWithOpReturnWithoutOPFalse.BEEF(),
				Annotations: transaction.Annotations{
					Outputs: transaction.OutputsAnnotations{
						0: &transaction.OutputAnnotation{
							Bucket: bucket.Data,
						},
					},
				},
			},
			expectTxID: givenTxWithOpReturnWithoutOPFalse.ID(),
			expectOutputs: []database.Output{
				{
					TxID:       givenTxWithOpReturnWithoutOPFalse.InputUTXO(0).TxID,
					Vout:       givenTxWithOpReturnWithoutOPFalse.InputUTXO(0).Vout,
					SpendingTX: ptr(givenTxWithOpReturnWithoutOPFalse.ID()),
				},
				{
					TxID:       givenTxWithOpReturnWithoutOPFalse.ID(),
					Vout:       0,
					SpendingTX: nil,
				},
			},
			expectData: []database.Data{
				{
					TxID: givenTxWithOpReturnWithoutOPFalse.ID(),
					Vout: 0,
					Blob: []byte(dataOfOpReturnTx),
				},
			},
		},
		"RecordTransactionOutline for op_return with untracked utxo as inputs": {
			repo: newMockRepository(),
			outline: &outlines.Transaction{
				BEEF: givenTxWithOpReturn.BEEF(),
				Annotations: transaction.Annotations{
					Outputs: transaction.OutputsAnnotations{
						0: &transaction.OutputAnnotation{
							Bucket: bucket.Data,
						},
					},
				},
			},
			expectTxID: givenTxWithOpReturn.ID(),
			expectOutputs: []database.Output{{
				TxID: givenTxWithOpReturn.ID(),
				Vout: 0,
			}},
			expectData: []database.Data{
				{
					TxID: givenTxWithOpReturn.ID(),
					Vout: 0,
					Blob: []byte(dataOfOpReturnTx),
				},
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			broadcaster := newMockBroadcaster()
			repo := test.repo
			service := record.NewService(tester.Logger(t), repo, broadcaster)

			// when:
			err := service.RecordTransactionOutline(context.Background(), test.outline)

			// then:
			require.NoError(t, err)

			require.Contains(t, broadcaster.broadcastedTxs, test.expectTxID)

			require.Contains(t, repo.transactions, test.expectTxID)
			txEntry := repo.transactions[test.expectTxID]
			require.Equal(t, test.expectTxID, repo.transactions[test.expectTxID].ID)
			require.Equal(t, database.TxStatusBroadcasted, txEntry.TxStatus)

			require.Subset(t, repo.getAllOutputs(), test.expectOutputs)
			require.Subset(t, repo.getAllData(), test.expectData)
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
		WithOutputScript(
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

	givenTxWithOpReturn := fixtures.GivenTX(t).
		WithInput(1).
		WithOPReturn(dataOfOpReturnTx)

	tests := map[string]struct {
		repo        *mockRepository
		outline     *outlines.Transaction
		broadcaster *mockBroadcaster
		expectErr   error
	}{
		"RecordTransactionOutline for not signed transaction": {
			broadcaster: newMockBroadcaster(),
			repo:        newMockRepository(),
			outline: &outlines.Transaction{
				BEEF: givenUnsignedTX.BEEF(),
			},
			expectErr: txerrors.ErrTxValidation,
		},
		"RecordTransactionOutline for not a BEEF hex": {
			broadcaster: newMockBroadcaster(),
			repo:        newMockRepository(),
			outline: &outlines.Transaction{
				BEEF: notABeefHex,
			},
			expectErr: txerrors.ErrTxValidation,
		},
		"RecordTransactionOutline for invalid OP_ZERO after OP_RETURN": {
			broadcaster: newMockBroadcaster(),
			repo:        newMockRepository(),
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
			broadcaster: newMockBroadcaster(),
			repo: newMockRepository().withOutput(database.Output{
				TxID:       givenTxWithOpReturn.InputUTXO(0).TxID,
				Vout:       givenTxWithOpReturn.InputUTXO(0).Vout,
				SpendingTX: ptr("05aa91319c773db18071310ecd5ddc15d3aa4242b55705a13a66f7fefe2b80a1"),
			}),
			outline: &outlines.Transaction{
				BEEF: givenTxWithOpReturn.BEEF(),
			},
			expectErr: txerrors.ErrUTXOSpent,
		},
		"Vout out of range in annotation": {
			broadcaster: newMockBroadcaster(),
			repo:        newMockRepository(),
			outline: &outlines.Transaction{
				BEEF: givenTxWithOpReturn.BEEF(),
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
			broadcaster: newMockBroadcaster(),
			repo:        newMockRepository(),
			outline: &outlines.Transaction{
				BEEF: givenTxWithOpReturn.BEEF(),
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
			broadcaster: newMockBroadcaster(),
			repo:        newMockRepository(),
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
		"error during broadcasting": {
			broadcaster: newMockBroadcaster().withError(errors.New("broadcast error")),
			repo: newMockRepository().withOutput(database.Output{
				TxID: givenTxWithOpReturn.InputUTXO(0).TxID,
				Vout: givenTxWithOpReturn.InputUTXO(0).Vout,
			}),
			outline: &outlines.Transaction{
				BEEF: givenTxWithOpReturn.BEEF(),
				Annotations: transaction.Annotations{
					Outputs: transaction.OutputsAnnotations{
						0: &transaction.OutputAnnotation{
							Bucket: bucket.Data,
						},
					},
				},
			},
			expectErr: txerrors.ErrTxBroadcast,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			service := record.NewService(tester.Logger(t), test.repo, test.broadcaster)
			initialOutputs := test.repo.getAllOutputs()
			initialData := test.repo.getAllData()

			// when:
			err := service.RecordTransactionOutline(context.Background(), test.outline)

			// then:
			require.Error(t, err)
			require.ErrorIs(t, err, test.expectErr)

			// ensure that no changes were made to the repository
			require.ElementsMatch(t, initialOutputs, test.repo.getAllOutputs())
			require.ElementsMatch(t, initialData, test.repo.getAllData())
		})
	}
}

func ptr[T any](value T) *T {
	return &value
}

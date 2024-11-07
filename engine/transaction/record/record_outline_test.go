package record_test

import (
	"context"
	"testing"

	"github.com/bitcoin-sv/go-sdk/script"
	trx "github.com/bitcoin-sv/go-sdk/transaction"
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

	beefWithOnePTPKH = "0100beef01fe47390d000802b5003af2d2162c10310bee0e861e8a8dc94bd896d416e784aa22f667501ac270bcc5b402b23c7c47320b3818c665bf28a46c290f3fb379ea8d357625bbff3117ae14b09b015b003640de1da50be7bcc0a7a6fd56f55cf70c9818394c2f99f71a43a777607b3d38012c0059c7771f2e6b337f2bd826b4f35b3159b801d6a8fe8c485c4eff9e902f341da8011700e24840f82e0d91356217fe8854eb28d70b347d231389a6c978d4c89eb3d034ac010a0064841631ac16075ff068de9bd4f4a4a1c74394284e584e7d3060b594be31d39b0104001097bf4eda72ca3d8092b7c7eb1e17251e377dff79a9104d24ee0a111f7119800103009ebf3de54b3d0380536be7938bb0d9ee9145d4a4a2ad96d65f3011b58424b844010000a09fdac73af0f273b576ac48d330362b7552a0508c987ecd26d4a455bf0a24ab02010000000238163d2fd3cb01c87476437f37428fd0680ec3fe89a7057678b00b0435fdddde010000006b483045022100c141a1b551dc3f9a2d30a4234933d6d6aa86a889a143a1e6d613713dc3f2ed68022038b631deef1b229f2da84536286cf5d5377c148849f54ccb07d7957ce22e732c412102ab62d27fc4692c260b30d856ab68d8ca82d4f1fe3aec805250d8f850faa71827ffffffff06656a05bcd886e933f99088ab038e416a0f1e28fd7d856dfdb4aa31d7335ed7010000006a473044022065025d2c32ae5fb2e53732b1c48d4e12b3cb14bb6e66ae535e779f94932d303e02201906963e4c37d6e73890e65f7944d53a233255f8ad95b8bd1ee420ff3cea87084121027249db5e1c879cdc4ed98dc24d3a94d64f364dbb9f84b039cd44e415ba282305ffffffff0202000000000000001976a914298b970788e84da5f82276638ad6c0204d7622a988ac12000000000000001976a914b69e83df1129f7f442c7e0d54c31681f64ae9d7f88ac0000000001000100000001b23c7c47320b3818c665bf28a46c290f3fb379ea8d357625bbff3117ae14b09b000000006a47304402202a28b2fda3d37411189fb0d539f054e14b805a8a181d4bf7656fd8f8933fba3702207598eab7b1470ee98ae178a4d5a7bacde74e13e3d694295406051e620a409adc412102aaec858e0431eaf8a056d37a55639a18aa189e45eeb8bb94fe33880aa3ddd65effffffff0101000000000000001976a914f4838029c3838a94b5b7736f2fe82d860e0239df88ac0000000000"
)

func TestRecordOutlineOpReturn(t *testing.T) {
	givenTxWithOpReturn := fixtures.GivenTX().WithT(t).
		WithSign().
		WithInput(1).
		WithOPReturn(dataOfOpReturnTx)

	givenTxWithOpReturnWithoutOPFalse := fixtures.GivenTX().WithT(t).
		WithSign().
		WithInput(1).
		WithOutputFunc(func() *trx.TransactionOutput {
			s := &script.Script{}
			_ = s.AppendOpcodes(script.OpRETURN)
			_ = s.AppendPushDataArray([][]byte{[]byte(dataOfOpReturnTx)})
			return &trx.TransactionOutput{LockingScript: s}
		})

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
				BEEF: givenTxWithOpReturn.BeefHex(),
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
				BEEF: givenTxWithOpReturnWithoutOPFalse.BeefHex(),
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
				BEEF: givenTxWithOpReturn.BeefHex(),
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
	givenUnsignedTX := fixtures.GivenTX().WithT(t).
		WithInput(1).
		WithOPReturn(dataOfOpReturnTx)

	givenTxWithOpZeroAfterOpReturn := fixtures.GivenTX().WithT(t).
		WithSign().
		WithInput(1).
		WithOutputFunc(func() *trx.TransactionOutput {
			s := &script.Script{}
			_ = s.AppendOpcodes(script.OpFALSE, script.OpRETURN)
			_ = s.AppendPushDataArray([][]byte{[]byte(dataOfOpReturnTx)})
			_ = s.AppendOpcodes(script.OpZERO)
			_ = s.AppendOpcodes(script.OpZERO)
			_ = s.AppendPushDataArray([][]byte{[]byte(dataOfOpReturnTx)})
			return &trx.TransactionOutput{LockingScript: s}
		})

	givenTxWithOpReturn := fixtures.GivenTX().WithT(t).
		WithSign().
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
				BEEF: givenUnsignedTX.BeefHex(),
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
				BEEF: givenTxWithOpZeroAfterOpReturn.BeefHex(),
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
				BEEF: givenTxWithOpReturn.BeefHex(),
			},
			expectErr: txerrors.ErrUTXOSpent,
		},
		"Vout out of range in annotation": {
			broadcaster: newMockBroadcaster(),
			repo:        newMockRepository(),
			outline: &outlines.Transaction{
				BEEF: givenTxWithOpReturn.BeefHex(),
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
				BEEF: givenTxWithOpReturn.BeefHex(),
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
				BEEF: beefWithOnePTPKH,
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
				BEEF: givenTxWithOpReturn.BeefHex(),
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

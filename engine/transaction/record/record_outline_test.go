package record

import (
	"context"
	"fmt"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction/outlines"
	"github.com/bitcoin-sv/spv-wallet/models/transaction/bucket"
	"github.com/stretchr/testify/require"
	"iter"
	"testing"

	trx "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/database"
	"github.com/bitcoin-sv/spv-wallet/engine/tester"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
)

type mockRepository struct {
	transactions map[string]database.Transaction
	outputs      map[string]database.Output
	data         map[string]database.Data
}

func newMockRepository() *mockRepository {
	return &mockRepository{
		transactions: make(map[string]database.Transaction),
		outputs:      make(map[string]database.Output),
		data:         make(map[string]database.Data),
	}
}

func (m *mockRepository) outpointID(txID string, vout uint32) string {
	return fmt.Sprintf("%s-%d", txID, vout)
}

func (m *mockRepository) SaveTX(_ context.Context, txTable *database.Transaction, outputs []database.Output, data []database.Data) error {
	m.transactions[txTable.ID] = *txTable
	for _, output := range outputs {
		m.outputs[m.outpointID(output.TxID, output.Vout)] = output
	}
	for _, d := range data {
		m.data[m.outpointID(d.TxID, d.Vout)] = d
	}
	return nil
}

func (m *mockRepository) GetOutputs(_ context.Context, outpoints iter.Seq[bsv.Outpoint]) ([]database.Output, error) {
	var outputs []database.Output
	for outpoint := range outpoints {
		key := m.outpointID(outpoint.TxID, outpoint.Vout)
		if output, ok := m.outputs[key]; ok {
			outputs = append(outputs, output)
		}
	}
	return outputs, nil
}

type mockBroadcaster struct {
	broadcastedTxs map[string]*trx.Transaction
}

func (m *mockBroadcaster) Broadcast(_ context.Context, tx *trx.Transaction) error {
	m.broadcastedTxs[tx.TxID().String()] = tx
	return nil
}

func newMockBroadcaster() *mockBroadcaster {
	return &mockBroadcaster{
		broadcastedTxs: make(map[string]*trx.Transaction),
	}
}

const (
	beefWithOpReturn = "0100beef01fe47390d000802b402b23c7c47320b3818c665bf28a46c290f3fb379ea8d357625bbff3117ae14b09bb5003af2d2162c10310bee0e861e8a8dc94bd896d416e784aa22f667501ac270bcc5015b003640de1da50be7bcc0a7a6fd56f55cf70c9818394c2f99f71a43a777607b3d38012c0059c7771f2e6b337f2bd826b4f35b3159b801d6a8fe8c485c4eff9e902f341da8011700e24840f82e0d91356217fe8854eb28d70b347d231389a6c978d4c89eb3d034ac010a0064841631ac16075ff068de9bd4f4a4a1c74394284e584e7d3060b594be31d39b0104001097bf4eda72ca3d8092b7c7eb1e17251e377dff79a9104d24ee0a111f7119800103009ebf3de54b3d0380536be7938bb0d9ee9145d4a4a2ad96d65f3011b58424b844010000a09fdac73af0f273b576ac48d330362b7552a0508c987ecd26d4a455bf0a24ab02010000000238163d2fd3cb01c87476437f37428fd0680ec3fe89a7057678b00b0435fdddde010000006b483045022100c141a1b551dc3f9a2d30a4234933d6d6aa86a889a143a1e6d613713dc3f2ed68022038b631deef1b229f2da84536286cf5d5377c148849f54ccb07d7957ce22e732c412102ab62d27fc4692c260b30d856ab68d8ca82d4f1fe3aec805250d8f850faa71827ffffffff06656a05bcd886e933f99088ab038e416a0f1e28fd7d856dfdb4aa31d7335ed7010000006a473044022065025d2c32ae5fb2e53732b1c48d4e12b3cb14bb6e66ae535e779f94932d303e02201906963e4c37d6e73890e65f7944d53a233255f8ad95b8bd1ee420ff3cea87084121027249db5e1c879cdc4ed98dc24d3a94d64f364dbb9f84b039cd44e415ba282305ffffffff0202000000000000001976a914298b970788e84da5f82276638ad6c0204d7622a988ac12000000000000001976a914b69e83df1129f7f442c7e0d54c31681f64ae9d7f88ac0000000001000100000001b23c7c47320b3818c665bf28a46c290f3fb379ea8d357625bbff3117ae14b09b000000006a47304402205febda774d651f1b15dcb50b450fb10356ae25ab887e9023f35e440138b3f19b022027af299175e8a42649b03d110c1e031f6c20a1d9fb1fa72af813185fec4ae0be412102aaec858e0431eaf8a056d37a55639a18aa189e45eeb8bb94fe33880aa3ddd65effffffff0100000000000000000e006a0b68656c6c6f20776f726c640000000000"
	txIDWithOpReturn = "af5d7b59b2973355a043c70e1fa0738bd0ceecd90af87fd52d4cdb88ac3eb10b"
	dataOfOpReturnTx = "hello world"
)

func TestRecordOutline(t *testing.T) {
	t.Run("RecordTransactionOutline", func(t *testing.T) {
		repo := newMockRepository()
		broadcaster := newMockBroadcaster()
		service := NewService(tester.Logger(t), repo, broadcaster)

		outline := outlines.Transaction{
			BEEF: beefWithOpReturn,
			Annotations: &transaction.Annotations{
				Outputs: transaction.OutputsAnnotations{
					0: &transaction.OutputAnnotation{
						Bucket: bucket.Data,
					},
				},
			},
		}

		err := service.RecordTransactionOutline(context.Background(), &outline)
		require.NoError(t, err)

		require.Contains(t, broadcaster.broadcastedTxs, txIDWithOpReturn)

		require.Contains(t, repo.transactions, txIDWithOpReturn)
		txEntry := repo.transactions[txIDWithOpReturn]
		require.Equal(t, txIDWithOpReturn, txEntry.ID)
		require.Equal(t, database.TxStatusBroadcasted, txEntry.TxStatus)

		require.Empty(t, repo.outputs)

		require.Contains(t, repo.data, repo.outpointID(txIDWithOpReturn, 0))
		data := repo.data[repo.outpointID(txIDWithOpReturn, 0)]
		require.Equal(t, txIDWithOpReturn, data.TxID)
		require.Equal(t, uint32(0), data.Vout)
		require.Equal(t, dataOfOpReturnTx, string(data.Blob))
	})
}

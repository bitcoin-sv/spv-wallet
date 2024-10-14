package internal_test

import (
	"context"
	"errors"
	"testing"
	"time"

	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/chain/internal"
	"github.com/bitcoin-sv/spv-wallet/engine/chain/models"
	"github.com/stretchr/testify/require"
	"iter"
)

type mockTxsGetter struct {
	transactions []*sdk.Transaction
	returnError  error
	applyTimeout bool
}

func (m *mockTxsGetter) GetTransactions(ctx context.Context, ids iter.Seq[string]) ([]*sdk.Transaction, error) {
	if m.applyTimeout {
		<-ctx.Done()
		return nil, ctx.Err()
	}
	if m.returnError != nil {
		return nil, m.returnError
	}
	return m.transactions, nil
}

const (
	tx1Hex = "0100000002c646da06628a5846b3c1ba79936bb0b2ba019ecd8acc2281441bd1aa4408f3cc000000006a4730440220500357d7e85c623405afd767db11d2da20b50c6fcf97c0db61c70867748e3ea802201c5fef8d7d1f9b4f24a5e5ccb84071a9b1f9808b5d475be21dbe475ce167006041210217ff58e102ed361d4946bafb257afc724af3c50ab108a4333d585ae81f230095ffffffff46324e8a83f0bcddd70358d8dd8d42613529dc633a14c001233bc8be959642df000000006b483045022100afd4023dd47cda9e7f6704fffa8740e29a07a41a0c5d4884043477bbb60082d902202b4921eb1ef7d8fe828c8213f98e1cbecea2c6acdd864a516ef29769cc6423b541210321d3a9b13c6c0b3d1e16bdc4c0990749c52cb394f6f06e314b1186c9c53d603affffffff0214000000000000001976a91405b1321ee4478faa0a3d0aae8893c149811ec68e88ac13000000000000001976a91420aeb57d4809d74b7bb78ccbb7eb07fc9c59e1cc88ac00000000"
	tx2Hex = "0100000002779082a94ec3c113337afee5667b69d64d7d7591472c92a4cde37db15719ea45010000006a473044022038ecc2b1482df117871722b634a24cf83d18008a26ba2e98f342ab8b1df4dfed02202e0b1e989726134910013329b244dee9b7b8e48a3ba2eb5926bb375ca67233b2412102d8d04ed1c01919f0973e1091e054edb7f3af469929ae5e12d01af32f5595b5ccffffffffe071ef01d65fcbf6b0882ea2216a67928aa05b6ffa34ecd32fc950e96acad520010000006a47304402205db7b79a3075e7bebbbbacaa9fb583e1b008a82deaafd4512ea8b84bb1d4dca3022044a6a0a102499521333bcb8ca4f34025e8c0c80f3bd79abbe58f1aa0af7c71db4121023ecd8b030eb23991cc5abe718a78ec80abba740a6cd199037b13f1ab4e0c7376ffffffff0314000000000000001976a914975f6eef3229c0896ada2ead550702af3fda512888ac04100000000000001976a914486633d2c0e1eb6f7fa79c525c4e93d0a8b4c20688ac6d000000000000001976a91490b7fc2adfec1d9daa0622f857a06d9bb12ee21888ac00000000"
	tx3Hex = "01000000026d9a56036235bb5b5e39b04b6f188c74bc45189a9122f235cad5f0c4b668817d010000006a47304402205bad758ddad1816827d2f8d683d055c95652fa3e8902c0af1a2009b039e360350220152ac3f8417ef311ff61636fa91e99782a77e4d1d192c33fb2f1f245913b879f412103daf4a6e60ee877fa7e0639e0a4f416e5c80dfa05cd17762d77c62391a3322a52ffffffff08b11a9c1c6534f7d07d5367d4d7b073ca7803139122cbdc06f5e77463746310010000006a47304402204c8d367127e9d68b79c4d40bded75e773194ced7ab07345df966d5cbabca460102204d029f791d80b75d59fd2ba7b9e186e1d449e2e4e5c3d03a0f3e9f98b4bd132041210330278625061e9bb6104e959851e6698f2b96333cbd7262ee3c22410f2e5cf4a1ffffffff0314000000000000001976a9149df7244a4edade5f8d06ad392d568bc62218fe3988ac34120000000000001976a9144d4dec9ae2199860fc6ebe9cd446472f32f6ffd388acb5000000000000001976a914fd8154ad427ca5cce7209c121a29452add41e25d88ac00000000"
)

func TestCombinedTxGetter(t *testing.T) {
	tx1 := fromHex(tx1Hex)
	tx2 := fromHex(tx2Hex)
	tx3 := fromHex(tx3Hex)
	tests := map[string]struct {
		getters      []chainmodels.TransactionsGetter
		requestedTXs []*sdk.Transaction
		expectedTXs  []*sdk.Transaction
	}{
		"Transactions from single getter": {
			getters: []chainmodels.TransactionsGetter{
				&mockTxsGetter{transactions: []*sdk.Transaction{tx1, tx2}},
			},
			requestedTXs: []*sdk.Transaction{tx1, tx2},
			expectedTXs:  []*sdk.Transaction{tx1, tx2},
		},
		"Transactions from two getters": {
			getters: []chainmodels.TransactionsGetter{
				&mockTxsGetter{transactions: []*sdk.Transaction{tx1, tx2}},
				&mockTxsGetter{transactions: []*sdk.Transaction{tx3}},
			},
			requestedTXs: []*sdk.Transaction{tx1, tx2, tx3},
			expectedTXs:  []*sdk.Transaction{tx1, tx2, tx3},
		},
		"GetTransactions with empty requested transaction list": {
			getters: []chainmodels.TransactionsGetter{
				&mockTxsGetter{transactions: []*sdk.Transaction{tx1}},
			},
			requestedTXs: []*sdk.Transaction{},
			expectedTXs:  []*sdk.Transaction{},
		},
		"CombinedTxsGetter with no getters": {
			getters:      []chainmodels.TransactionsGetter{},
			requestedTXs: []*sdk.Transaction{tx1},
			expectedTXs:  []*sdk.Transaction{},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			getter := internal.NewCombinedTxsGetter(test.getters...)
			transactions, err := getter.GetTransactions(context.Background(), ids(test.requestedTXs...))

			require.NoError(t, err)
			require.Equal(t, len(test.expectedTXs), len(transactions))
			shouldAllContain(t, transactions, ids(test.expectedTXs...))
		})
	}
}

func TestCombinedTxGetterErrorCases(t *testing.T) {
	tx1 := fromHex(tx1Hex)

	t.Run("Getter returns error", func(t *testing.T) {
		expectedErr := errors.New("some error")
		getter := internal.NewCombinedTxsGetter(&mockTxsGetter{returnError: expectedErr})

		transactions, err := getter.GetTransactions(context.Background(), ids(tx1))

		require.ErrorIs(t, err, expectedErr)
		require.Nil(t, transactions)
	})

	t.Run("Getter interrupted by ctx timeout ", func(t *testing.T) {
		getter := internal.NewCombinedTxsGetter(&mockTxsGetter{applyTimeout: true})

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
		defer cancel()

		transactions, err := getter.GetTransactions(ctx, ids(tx1))

		require.ErrorIs(t, err, context.DeadlineExceeded)
		require.Nil(t, transactions)
	})
}

func fromHex(hex string) *sdk.Transaction {
	tx, _ := sdk.NewTransactionFromHex(hex)
	return tx
}

func id(tx *sdk.Transaction) string {
	return tx.TxID().String()
}

func ids(txs ...*sdk.Transaction) iter.Seq[string] {
	return func(yield func(string) bool) {
		for _, tx := range txs {
			if !yield(id(tx)) {
				return
			}
		}
	}
}

func shouldAllContain(t *testing.T, transactions []*sdk.Transaction, expectedTxIDs iter.Seq[string]) {
	txs := make(map[string]bool)
	for _, tx := range transactions {
		txs[id(tx)] = true
	}
	for expectedTxID := range expectedTxIDs {
		if _, exists := txs[expectedTxID]; !exists {
			require.Failf(t, "transaction %s not found", expectedTxID)
		}
	}
}

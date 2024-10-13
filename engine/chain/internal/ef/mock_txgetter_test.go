package ef_test

import (
	"context"
	"slices"
	"testing"

	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"iter"
)

type onMissingTxBehavior int

const (
	onMissingTxReturnError onMissingTxBehavior = iota
	onMissingTxAddNil
	onMissingTxSkip
)

type mockTransactionsGetter struct {
	transactions []*sdk.Transaction
	onMissingTx  onMissingTxBehavior
	returnAll    bool
}

func newMockTransactionsGetter(t *testing.T, rawTXBase64 []string) *mockTransactionsGetter {
	transactions := make([]*sdk.Transaction, 0, len(rawTXBase64))
	for _, rawTX := range rawTXBase64 {
		tx := fromHex(t, rawTX)
		transactions = append(transactions, tx)
	}
	return &mockTransactionsGetter{
		transactions: transactions,
		onMissingTx:  onMissingTxReturnError,
		returnAll:    false,
	}
}

func (m *mockTransactionsGetter) WithOnMissingBehavior(behavior onMissingTxBehavior) *mockTransactionsGetter {
	m.onMissingTx = behavior
	return m
}

func (m *mockTransactionsGetter) WithReturnAll(value bool) *mockTransactionsGetter {
	m.returnAll = value
	return m
}

func (m *mockTransactionsGetter) GetTransactions(_ context.Context, ids iter.Seq[string]) ([]*sdk.Transaction, error) {
	var result []*sdk.Transaction
	if m.returnAll {
		return append(result, m.transactions...), nil
	}

	for id := range ids {
		index := slices.IndexFunc(m.transactions, func(tx *sdk.Transaction) bool {
			return tx.TxID().String() == id
		})
		if index == -1 {
			switch m.onMissingTx {
			case onMissingTxReturnError:
				return nil, spverrors.Newf("transaction with ID %s not found", id)
			case onMissingTxAddNil:
				result = append(result, nil)
				continue
			case onMissingTxSkip:
				continue
			}
		} else {
			result = append(result, m.transactions[index])
		}
	}
	return result, nil
}

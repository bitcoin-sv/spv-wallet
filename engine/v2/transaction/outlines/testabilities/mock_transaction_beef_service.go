package testabilities

import (
	"context"
	"testing"

	sdk "github.com/bitcoin-sv/go-sdk/transaction"
)

type transactionBEEFService struct {
	t testing.TB
}

func (t *transactionBEEFService) PrepareBEEF(ctx context.Context, tx *sdk.Transaction) (string, error) {
	tmpTx := &sdk.Transaction{Outputs: tx.Outputs}
	return tmpTx.BEEFHex()
}

func newTransactionBEEFServiceMock(t testing.TB) *transactionBEEFService {
	return &transactionBEEFService{t: t}
}

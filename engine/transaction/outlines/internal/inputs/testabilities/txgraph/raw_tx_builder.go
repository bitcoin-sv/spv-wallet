package txgraph

import (
	"testing"

	"github.com/bitcoin-sv/go-sdk/script"
	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/go-sdk/transaction/template/p2pkh"
	"github.com/stretchr/testify/require"
)

type RawTxBuilder struct {
	T                    *testing.T
	P2PKHLockingScript   *script.Script
	P2PKHUnlockingScript *p2pkh.P2PKH
}

type AscendantTx struct {
	Tx   *sdk.Transaction
	Vout uint32
}

func (t *RawTxBuilder) MakeTx(ascendants ...AscendantTx) *sdk.Transaction {
	t.T.Helper() // todo: add nil handling

	tx := sdk.NewTransaction()
	for _, asc := range ascendants {
		tx.AddInputFromTx(asc.Tx, asc.Vout, t.P2PKHUnlockingScript)
		tx.AddOutput(&sdk.TransactionOutput{LockingScript: t.P2PKHLockingScript})
	}

	require.NoError(t.T, tx.Sign(), "failed to sign transaction")
	return tx
}

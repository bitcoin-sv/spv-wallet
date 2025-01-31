package beef

import (
	"testing"

	sdk "github.com/bitcoin-sv/go-sdk/transaction"

	"github.com/stretchr/testify/require"
)

func TestSourceTransactionResolver_BEEFGrandparentForTx1(t *testing.T) {
	// given:
	graphBuilder := NewTransactionGraphBuilder(t)

	// graph:
	tx6 := graphBuilder.CreateMinedTx("tx6", 1)
	tx1 := graphBuilder.CreateRawTx("tx1", ParentTx{Tx: tx6, Vout: 0})
	tx0 := graphBuilder.CreateRawTx("tx0", ParentTx{Tx: tx1, Vout: 0})
	graphBuilder.EnsureGraphIsValid()

	// subject tx:
	subjectTx := sdk.NewTransaction()
	subjectTx.AddInput(&sdk.TransactionInput{SourceTXID: tx0.Inputs[0].SourceTXID})

	// when:
	resolver, err := NewSourceTransactionResolver(subjectTx, graphBuilder.ToTxQueryResultSlice())
	require.NoError(t, err)

	err = resolver.Resolve()

	// then:
	require.NoError(t, err)
}

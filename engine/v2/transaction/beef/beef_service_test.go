package beef_test

import (
	"context"
	"testing"

	"github.com/bitcoin-sv/go-sdk/chainhash"
	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/beef"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/beef/testabilities"
	txerrors "github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/errors"
	"github.com/stretchr/testify/require"
)

func TestPrepareBEEF_MissingSourceTx_InSubjectTxInput(t *testing.T) {
	// given:
	graphBuilder := testabilities.NewTxGraphBuilder(t)
	repository := testabilities.NewTxRepository(t, graphBuilder)
	service := beef.NewService(repository)

	var hexGen testabilities.HexGen
	ID, err := chainhash.NewHashFromHex(hexGen.Val())
	require.NoError(t, err, "expected to create hash from the hex gen val: %s", hexGen.Val())

	// subject tx:
	subjectTx := sdk.NewTransaction()
	subjectTx.AddInput(&sdk.TransactionInput{SourceTXID: ID})

	// when:
	hexBEEF, err := service.PrepareBEEF(context.Background(), subjectTx)

	// then:
	thenTx := testabilities.Then(t)

	thenTx.
		IsEmpty(hexBEEF).
		HasError(err, txerrors.ErrInvalidTransactionInput)
}

func TestPrepareBEEF_GraphTx_WithBEEFGrandparent(t *testing.T) {
	// given:
	graphBuilder := testabilities.NewTxGraphBuilder(t)
	repository := testabilities.NewTxRepository(t, graphBuilder)
	service := beef.NewService(repository)

	// graph:
	tx6 := graphBuilder.CreateMinedTx("tx6", 1)
	tx1 := graphBuilder.CreateRawTx("tx1", testabilities.ParentTx{Tx: tx6, Vout: 0})
	tx0 := graphBuilder.CreateRawTx("tx0", testabilities.ParentTx{Tx: tx1, Vout: 0})
	graphBuilder.EnsureGraphIsValid()

	subjectTx := sdk.NewTransaction()
	subjectTx.AddInput(&sdk.TransactionInput{SourceTXID: tx0.Inputs[0].SourceTXID})

	// when:
	hexBEEF, err := service.PrepareBEEF(context.Background(), subjectTx)

	// then:
	thenTx := testabilities.Then(t)

	thenTx.
		Created(hexBEEF).
		WithNoError(err).
		WithParseableBEEFHEX().
		WithSourceTransactions()
}

func TestPrepareBEEF_GraphTx_WithBEEFGrandparents_Tx1_Tx2_Tx3(t *testing.T) {
	// given:
	graphBuilder := testabilities.NewTxGraphBuilder(t)
	repository := testabilities.NewTxRepository(t, graphBuilder)
	service := beef.NewService(repository)

	// graph:
	tx4 := graphBuilder.CreateMinedTx("tx4", 1)
	tx1 := graphBuilder.CreateRawTx("tx1", testabilities.ParentTx{Tx: tx4, Vout: 0})

	tx5 := graphBuilder.CreateMinedTx("tx5", 1)
	tx3 := graphBuilder.CreateRawTx("tx3", testabilities.ParentTx{Tx: tx5, Vout: 0})

	tx6 := graphBuilder.CreateMinedTx("tx6", 1)
	tx2 := graphBuilder.CreateRawTx("tx2", testabilities.ParentTx{Tx: tx6, Vout: 0})
	graphBuilder.EnsureGraphIsValid()

	tx0 := graphBuilder.CreateRawTx("tx0",
		testabilities.ParentTx{Tx: tx1, Vout: 0},
		testabilities.ParentTx{Tx: tx3, Vout: 0},
		testabilities.ParentTx{Tx: tx2, Vout: 0},
	)

	subjectTx := sdk.NewTransaction()
	subjectTx.AddInput(&sdk.TransactionInput{SourceTXID: tx0.Inputs[0].SourceTXID})
	subjectTx.AddInput(&sdk.TransactionInput{SourceTXID: tx0.Inputs[1].SourceTXID})
	subjectTx.AddInput(&sdk.TransactionInput{SourceTXID: tx0.Inputs[2].SourceTXID})

	// when:
	hexBEEF, err := service.PrepareBEEF(context.Background(), subjectTx)

	// then:
	thenTx := testabilities.Then(t)

	thenTx.
		Created(hexBEEF).
		WithNoError(err).
		WithParseableBEEFHEX().
		WithSourceTransactions()
}

func TestPrepareBEEF_GraphTx_WithBEEFParents_Tx0(t *testing.T) {
	// given:
	graphBuilder := testabilities.NewTxGraphBuilder(t)
	repository := testabilities.NewTxRepository(t, graphBuilder)
	service := beef.NewService(repository)

	// graph:
	tx1 := graphBuilder.CreateMinedTx("tx4", 1)
	tx3 := graphBuilder.CreateMinedTx("tx3", 1)
	tx2 := graphBuilder.CreateMinedTx("tx2", 1)
	tx0 := graphBuilder.CreateRawTx(
		"tx0",
		testabilities.ParentTx{Tx: tx1, Vout: 0},
		testabilities.ParentTx{Tx: tx3, Vout: 0},
		testabilities.ParentTx{Tx: tx2, Vout: 0},
	)

	subjectTx := sdk.NewTransaction()
	subjectTx.AddInput(&sdk.TransactionInput{SourceTXID: tx0.Inputs[0].SourceTXID})
	subjectTx.AddInput(&sdk.TransactionInput{SourceTXID: tx0.Inputs[1].SourceTXID})
	subjectTx.AddInput(&sdk.TransactionInput{SourceTXID: tx0.Inputs[2].SourceTXID})

	// when:
	hexBEEF, err := service.PrepareBEEF(context.Background(), subjectTx)

	// then:
	thenTx := testabilities.Then(t)

	thenTx.
		Created(hexBEEF).
		WithNoError(err).
		WithParseableBEEFHEX().
		WithSourceTransactions()
}

func TestPrepareBEEF_GraphTx_WithBEEFGreatGrandparents_Tx1_Tx3(t *testing.T) {
	// given:
	graphBuilder := testabilities.NewTxGraphBuilder(t)
	repository := testabilities.NewTxRepository(t, graphBuilder)
	service := beef.NewService(repository)

	// graph:
	tx6 := graphBuilder.CreateMinedTx("tx6", 1)
	tx4 := graphBuilder.CreateRawTx("tx4", testabilities.ParentTx{Tx: tx6, Vout: 0})
	tx1 := graphBuilder.CreateRawTx("tx1", testabilities.ParentTx{Tx: tx4, Vout: 0})

	tx7 := graphBuilder.CreateMinedTx("tx7", 1)
	tx5 := graphBuilder.CreateRawTx("tx5", testabilities.ParentTx{Tx: tx7, Vout: 0})
	tx3 := graphBuilder.CreateRawTx("tx3", testabilities.ParentTx{Tx: tx5, Vout: 0})

	tx9 := graphBuilder.CreateMinedTx("tx9", 1)
	tx8 := graphBuilder.CreateRawTx("tx8", testabilities.ParentTx{Tx: tx9, Vout: 0})
	tx2 := graphBuilder.CreateRawTx("tx2", testabilities.ParentTx{Tx: tx8, Vout: 0})

	tx0 := graphBuilder.CreateRawTx(
		"tx0",
		testabilities.ParentTx{Tx: tx1, Vout: 0},
		testabilities.ParentTx{Tx: tx3, Vout: 0},
		testabilities.ParentTx{Tx: tx2, Vout: 0},
	)

	subjectTx := sdk.NewTransaction()
	subjectTx.AddInput(&sdk.TransactionInput{SourceTXID: tx0.Inputs[0].SourceTXID})
	subjectTx.AddInput(&sdk.TransactionInput{SourceTXID: tx0.Inputs[1].SourceTXID})
	subjectTx.AddInput(&sdk.TransactionInput{SourceTXID: tx0.Inputs[2].SourceTXID})

	// when:
	hexBEEF, err := service.PrepareBEEF(context.Background(), subjectTx)

	// then:
	thenTx := testabilities.Then(t)

	thenTx.
		Created(hexBEEF).
		WithNoError(err).
		WithParseableBEEFHEX().
		WithSourceTransactions()
}

func TestPrepareBEEF_GraphTx_WithSharedBEEFGrandparent_Tx1_Tx3(t *testing.T) {
	// given:
	graphBuilder := testabilities.NewTxGraphBuilder(t)
	repository := testabilities.NewTxRepository(t, graphBuilder)
	service := beef.NewService(repository)

	// graph:
	tx6 := graphBuilder.CreateMinedTx("tx6", 2)
	tx5 := graphBuilder.CreateRawTx("tx5", testabilities.ParentTx{Tx: tx6, Vout: 0})
	tx3 := graphBuilder.CreateRawTx("tx3", testabilities.ParentTx{Tx: tx5, Vout: 0})

	tx4 := graphBuilder.CreateRawTx("tx4", testabilities.ParentTx{Tx: tx6, Vout: 0})
	tx1 := graphBuilder.CreateRawTx("tx1", testabilities.ParentTx{Tx: tx4, Vout: 0})

	tx7 := graphBuilder.CreateMinedTx("tx7", 1)
	tx2 := graphBuilder.CreateRawTx("tx2", testabilities.ParentTx{Tx: tx7, Vout: 0})

	tx0 := graphBuilder.CreateRawTx(
		"tx0",
		testabilities.ParentTx{Tx: tx1, Vout: 0},
		testabilities.ParentTx{Tx: tx3, Vout: 0},
		testabilities.ParentTx{Tx: tx2, Vout: 0},
	)

	subjectTx := sdk.NewTransaction()
	subjectTx.AddInput(&sdk.TransactionInput{SourceTXID: tx0.Inputs[0].SourceTXID})
	subjectTx.AddInput(&sdk.TransactionInput{SourceTXID: tx0.Inputs[1].SourceTXID})
	subjectTx.AddInput(&sdk.TransactionInput{SourceTXID: tx0.Inputs[2].SourceTXID})

	// when:
	hexBEEF, err := service.PrepareBEEF(context.Background(), subjectTx)

	// then:
	thenTx := testabilities.Then(t)

	thenTx.
		Created(hexBEEF).
		WithNoError(err).
		WithParseableBEEFHEX().
		WithSourceTransactions()
}

func TestPrepareBEEF_GraphTx_SharedBEEFGreatGrandparents_Tx1_Tx3(t *testing.T) {
	// given:
	graphBuilder := testabilities.NewTxGraphBuilder(t)
	repository := testabilities.NewTxRepository(t, graphBuilder)
	service := beef.NewService(repository)

	// graph:
	tx6 := graphBuilder.CreateMinedTx("tx6", 1)
	tx4 := graphBuilder.CreateRawTx("tx4", testabilities.ParentTx{Tx: tx6, Vout: 0})
	tx1 := graphBuilder.CreateRawTx("tx1", testabilities.ParentTx{Tx: tx4, Vout: 0})

	tx7 := graphBuilder.CreateMinedTx("tx7", 1)
	tx5 := graphBuilder.CreateRawTx("tx5", testabilities.ParentTx{Tx: tx7, Vout: 0})
	tx3 := graphBuilder.CreateRawTx("tx3", testabilities.ParentTx{Tx: tx5, Vout: 0})

	tx9 := graphBuilder.CreateMinedTx("tx9", 1)
	tx8 := graphBuilder.CreateRawTx("tx8", testabilities.ParentTx{Tx: tx9, Vout: 0})
	tx2 := graphBuilder.CreateRawTx("tx2", testabilities.ParentTx{Tx: tx8, Vout: 0})

	tx0 := graphBuilder.CreateRawTx(
		"tx0",
		testabilities.ParentTx{Tx: tx1, Vout: 0},
		testabilities.ParentTx{Tx: tx3, Vout: 0},
		testabilities.ParentTx{Tx: tx2, Vout: 0},
	)

	subjectTx := sdk.NewTransaction()
	subjectTx.AddInput(&sdk.TransactionInput{SourceTXID: tx0.Inputs[0].SourceTXID})
	subjectTx.AddInput(&sdk.TransactionInput{SourceTXID: tx0.Inputs[1].SourceTXID})
	subjectTx.AddInput(&sdk.TransactionInput{SourceTXID: tx0.Inputs[2].SourceTXID})

	// when:
	hexBEEF, err := service.PrepareBEEF(context.Background(), subjectTx)

	// then:
	thenTx := testabilities.Then(t)

	thenTx.
		Created(hexBEEF).
		WithNoError(err).
		WithParseableBEEFHEX().
		WithSourceTransactions()
}

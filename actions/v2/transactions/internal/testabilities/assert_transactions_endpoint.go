package testabilities

import (
	"testing"

	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/actions/testabilities"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TransactionsEndpointAssertions interface {
	Response(response *resty.Response) TransactionsResponseAssertions
}

type TransactionsResponseAssertions interface {
	testabilities.SPVWalletResponseAssertions
	ContainsValidTransaction(format string) TransactionDetailsAssertions
}

func Then(t testing.TB, app testabilities.SPVWalletApplicationFixture) TransactionsEndpointAssertions {
	return &transactionEndpointAssertions{
		t:                     t,
		applicationAssertions: testabilities.Then(t, app),
	}
}

type transactionEndpointAssertions struct {
	t                     testing.TB
	applicationAssertions testabilities.SPVWalletApplicationAssertions
}

func (a *transactionEndpointAssertions) Response(response *resty.Response) TransactionsResponseAssertions {
	return &transactionResponseAssertions{
		t:                           a.t,
		require:                     require.New(a.t),
		assert:                      assert.New(a.t),
		response:                    response,
		SPVWalletResponseAssertions: a.applicationAssertions.Response(response),
	}
}

type transactionResponseAssertions struct {
	testabilities.SPVWalletResponseAssertions
	t        testing.TB
	response *resty.Response
	require  *require.Assertions
	assert   *assert.Assertions
}

func (a *transactionResponseAssertions) ContainsValidTransaction(format string) TransactionDetailsAssertions {
	a.t.Helper()
	txHex := a.SPVWalletResponseAssertions.JSONValue().GetString("hex")
	var tx *sdk.Transaction
	var err error
	switch format {
	case "BEEF":
		tx, err = sdk.NewTransactionFromBEEFHex(txHex)
	case "RAW":
		tx, err = sdk.NewTransactionFromHex(txHex)
	default:
		a.t.Fatalf("unsupported format: %s", format)
	}

	a.require.NoError(err, "hex is not valid tx in format %s", format)
	a.assert.NotZero(tx.Version, "tx version is 0 which is not acceptable by nodes")
	return &transactionAssertions{
		t:       a.t,
		require: a.require,
		assert:  a.assert,
		tx:      tx,
	}
}

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
	ContainsValidBEEFHexInField(field string) TransactionDetailsAssertions
	ContainsValidRawTxHexInField(field string) TransactionDetailsAssertions
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

func (a *transactionResponseAssertions) ContainsValidBEEFHexInField(field string) TransactionDetailsAssertions {
	txHex := a.SPVWalletResponseAssertions.JSONValue().GetString(field)
	tx, err := sdk.NewTransactionFromBEEFHex(txHex)
	a.require.NoError(err, "hex is not valid BEEF")
	return &transactionAssertions{
		t:       a.t,
		require: a.require,
		tx:      tx,
	}
}

func (a *transactionResponseAssertions) ContainsValidRawTxHexInField(field string) TransactionDetailsAssertions {
	txHex := a.SPVWalletResponseAssertions.JSONValue().GetString(field)
	tx, err := sdk.NewTransactionFromHex(txHex)
	a.require.NoError(err, "hex is not valid raw tx")
	return &transactionAssertions{
		t:       a.t,
		require: a.require,
		tx:      tx,
	}
}

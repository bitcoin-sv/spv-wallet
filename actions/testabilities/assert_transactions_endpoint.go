package testabilities

import (
	"testing"

	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TransactionsEndpointAssertions interface {
	Response(response *resty.Response) TransactionsResponseAssertions
}

type TransactionsResponseAssertions interface {
	SPVWalletResponseAssertions
	ContainsValidBEEFHexInField(field string) TransactionsResponseAssertions
	ContainsValidRawTxHexInField(field string) TransactionsResponseAssertions
}

func NewTransactionsEndpointAssertions(t testing.TB, app SPVWalletApplicationFixture) TransactionsEndpointAssertions {
	return &transactionEndpointAssertions{
		t:                     t,
		applicationAssertions: Then(t, app),
	}
}

type transactionEndpointAssertions struct {
	t                     testing.TB
	applicationAssertions SPVWalletApplicationAssertions
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
	SPVWalletResponseAssertions
	t        testing.TB
	response *resty.Response
	require  *require.Assertions
	assert   *assert.Assertions
}

func (a *transactionResponseAssertions) ContainsValidBEEFHexInField(field string) TransactionsResponseAssertions {
	txHex := a.SPVWalletResponseAssertions.JSONValue().GetString(field)
	_, err := sdk.NewTransactionFromBEEFHex(txHex)
	a.require.NoError(err, "hex is not valid BEEF")
	return a
}

func (a *transactionResponseAssertions) ContainsValidRawTxHexInField(field string) TransactionsResponseAssertions {
	txHex := a.SPVWalletResponseAssertions.JSONValue().GetString(field)
	_, err := sdk.NewTransactionFromHex(txHex)
	a.require.NoError(err, "hex is not valid raw tx")
	return a
}

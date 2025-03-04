package testabilities

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet/api"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/require"
)

type SPVWalletAppUserAssertions interface {
	Balance() BalanceAssertions
	Operations() OperationsAssertions
}

type BalanceAssertions interface {
	IsEqualTo(expected bsv.Satoshis)
	IsGreaterThanOrEqualTo(expected bsv.Satoshis)
	IsZero()
}

type OperationsAssertions interface {
	Last() LastOperationAssertions
}

type LastOperationAssertions interface {
	WithTxID(txID string) LastOperationAssertions
	WithValue(value int64) LastOperationAssertions
	WithType(operationType string) LastOperationAssertions
	WithCounterparty(counterparty string) LastOperationAssertions
	WithNoCounterparty() LastOperationAssertions
	WithTxStatus(txStatus string) LastOperationAssertions
	WithBlockHeight(blockHeight int64) LastOperationAssertions
	WithBlockHash(blockHash string) LastOperationAssertions
}

type userAssertions struct {
	userClient *resty.Client
	t          testing.TB
	require    *require.Assertions
}

func (u *userAssertions) Balance() BalanceAssertions {
	return u
}

func (u *userAssertions) Operations() OperationsAssertions {
	return u
}

func (u *userAssertions) Last() LastOperationAssertions {
	u.t.Helper()

	var result api.ModelsOperationsSearchResult
	_, err := u.userClient.R().SetResult(&result).Get("/api/v2/operations/search")
	u.require.NoError(err)
	u.require.NotEmpty(result.Content, "No operations found")

	return &lastOperationAssertions{
		t:       u.t,
		require: u.require,
		content: result.Content[0],
	}
}

func (u *userAssertions) balance() bsv.Satoshis {
	u.t.Helper()

	var userInfo api.ModelsUserInfo
	_, err := u.userClient.R().SetResult(&userInfo).Get("/api/v2/users/current")
	u.require.NoError(err)

	return bsv.Satoshis(userInfo.CurrentBalance)
}

func (u *userAssertions) IsEqualTo(expected bsv.Satoshis) {
	u.t.Helper()
	actual := u.balance()
	require.Equal(u.t, expected, actual)
}

func (u *userAssertions) IsGreaterThanOrEqualTo(expected bsv.Satoshis) {
	u.t.Helper()
	actual := u.balance()
	require.GreaterOrEqual(u.t, actual, expected)
}

func (u *userAssertions) IsZero() {
	u.t.Helper()
	actual := u.balance()
	require.Zero(u.t, actual)
}

type lastOperationAssertions struct {
	t       testing.TB
	require *require.Assertions
	content api.ModelsOperation
}

func (l *lastOperationAssertions) WithTxID(txID string) LastOperationAssertions {
	l.t.Helper()
	l.require.Equal(txID, l.content.TxID)
	return l
}

func (l *lastOperationAssertions) WithValue(value int64) LastOperationAssertions {
	l.t.Helper()
	l.require.Equal(value, l.content.Value)
	return l
}

func (l *lastOperationAssertions) WithType(operationType string) LastOperationAssertions {
	l.t.Helper()
	l.require.EqualValues(operationType, l.content.Type)
	return l
}

func (l *lastOperationAssertions) WithCounterparty(counterparty string) LastOperationAssertions {
	l.t.Helper()
	l.require.Equal(counterparty, l.content.Counterparty)
	return l
}

func (l *lastOperationAssertions) WithNoCounterparty() LastOperationAssertions {
	l.t.Helper()
	l.require.Empty(l.content.Counterparty)
	return l
}

func (l *lastOperationAssertions) WithTxStatus(txStatus string) LastOperationAssertions {
	l.t.Helper()
	l.require.EqualValues(txStatus, l.content.TxStatus)
	return l
}

func (l *lastOperationAssertions) WithBlockHeight(blockHeight int64) LastOperationAssertions {
	l.t.Helper()
	l.require.NotNil(l.content.BlockHeight)
	l.require.EqualValues(blockHeight, *l.content.BlockHeight)
	return l
}

func (l *lastOperationAssertions) WithBlockHash(blockHash string) LastOperationAssertions {
	l.t.Helper()
	l.require.NotNil(l.content.BlockHash)
	l.require.Equal(blockHash, *l.content.BlockHash)
	return l
}

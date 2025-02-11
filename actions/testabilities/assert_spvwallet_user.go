package testabilities

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"github.com/bitcoin-sv/spv-wallet/models/response"
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
	WithTxStatus(txStatus string) LastOperationAssertions
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

	var result response.PageModel[response.Operation]
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

	var userInfo response.UserInfo
	_, err := u.userClient.R().SetResult(&userInfo).Get("/api/v2/users/current")
	u.require.NoError(err)

	return userInfo.CurrentBalance
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
	content *response.Operation
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
	l.require.Equal(operationType, l.content.Type)
	return l
}

func (l *lastOperationAssertions) WithCounterparty(counterparty string) LastOperationAssertions {
	l.t.Helper()
	l.require.Equal(counterparty, l.content.Counterparty)
	return l
}

func (l *lastOperationAssertions) WithTxStatus(txStatus string) LastOperationAssertions {
	l.t.Helper()
	l.require.Equal(txStatus, l.content.TxStatus)
	return l
}

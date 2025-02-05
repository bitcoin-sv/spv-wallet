package testabilities

import (
	"context"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"github.com/stretchr/testify/require"
)

type EngineAssertions interface {
	User(fixtures.User) UserAssertions
}

type UserAssertions interface {
	Balance() BalanceAssertions
}

type BalanceAssertions interface {
	IsEqualTo(expected bsv.Satoshis)
	IsGreaterThanOrEqualTo(expected bsv.Satoshis)
	IsZero()
}

type engineAssertions struct {
	engWithCfg *EngineWithConfig
	t          testing.TB
}

func Then(t testing.TB, engWithCfg *EngineWithConfig) EngineAssertions {
	return &engineAssertions{
		engWithCfg: engWithCfg,
		t:          t,
	}
}

type userAssertions struct {
	engWithCfg *EngineWithConfig
	t          testing.TB
	user       fixtures.User
	require    *require.Assertions
}

func (e *engineAssertions) User(user fixtures.User) UserAssertions {
	return &userAssertions{
		engWithCfg: e.engWithCfg,
		t:          e.t,
		user:       user,
		require:    require.New(e.t),
	}
}

func (u *userAssertions) Balance() BalanceAssertions {
	return u
}

func (u *userAssertions) balance() bsv.Satoshis {
	u.t.Helper()
	actual, err := u.engWithCfg.Engine.UsersService().GetBalance(context.Background(), u.user.ID())
	u.require.NoError(err)
	return actual
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

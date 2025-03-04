package testabilities

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet/actions/testabilities"
)

type IntegrationTestAssertion interface {
	Alice() SpvWalletActorsStateAssertions
	Bob() SpvWalletActorsStateAssertions
	Charlie() SpvWalletActorsStateAssertions
}

// SpvWalletActorsStateAssertions about spv-wallet users
type SpvWalletActorsStateAssertions interface {
	Balance() testabilities.BalanceAssertions
	Operations() testabilities.OperationsAssertions
}

type assertions struct {
	t       testing.TB
	fixture *fixture
	testabilities.SPVWalletApplicationAssertions
}

func (a *assertions) Alice() SpvWalletActorsStateAssertions {
	return &actorAssertions{
		userAssertions: a.User(a.fixture.alice.User),
	}
}

func (a *assertions) Bob() SpvWalletActorsStateAssertions {
	return &actorAssertions{
		userAssertions: a.User(a.fixture.bob.User),
	}
}

func (a *assertions) Charlie() SpvWalletActorsStateAssertions {
	return &actorAssertions{
		userAssertions: a.User(a.fixture.charlie.User),
	}
}

func newAssertions(t testing.TB, fixture *fixture, appAssertions testabilities.SPVWalletApplicationAssertions) IntegrationTestAssertion {
	return &assertions{
		t:                              t,
		fixture:                        fixture,
		SPVWalletApplicationAssertions: appAssertions,
	}
}

type actorAssertions struct {
	userAssertions testabilities.SPVWalletAppUserAssertions
}

// Balance returns balance assertions for the actor
func (a actorAssertions) Balance() testabilities.BalanceAssertions {
	return a.userAssertions.Balance()
}

// Operations returns operation assertions for the actor
func (a actorAssertions) Operations() testabilities.OperationsAssertions {
	return a.userAssertions.Operations()
}

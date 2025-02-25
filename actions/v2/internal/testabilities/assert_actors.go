package testabilities

import "github.com/bitcoin-sv/spv-wallet/actions/testabilities"

// Then starts an assertion chain for the specified actor
func (tc *ActorTestContext) Then(actor *User) ActorAssertions {
	return ActorAssertions{
		userAssertions: tc.baseAssertions.User(actor.User),
	}
}

// Balance returns balance assertions for the actor
func (a ActorAssertions) Balance() testabilities.BalanceAssertions {
	return a.userAssertions.Balance()
}

// Operations returns operation assertions for the actor
func (a ActorAssertions) Operations() testabilities.OperationsAssertions {
	return a.userAssertions.Operations()
}

package testabilities

import "testing"

func New(t testing.TB) (given DraftTransactionFixture, then DraftTransactionAssertion) {
	return Given(t), Then(t)
}

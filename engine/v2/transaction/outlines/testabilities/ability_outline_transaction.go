package testabilities

import "testing"

func New(t testing.TB) (given TransactionOutlineFixture, then TransactionOutlineAssertion) {
	given = Given(t)
	then = Then(t, given)
	return
}

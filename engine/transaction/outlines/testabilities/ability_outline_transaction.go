package testabilities

import "testing"

func New(t testing.TB) (given TransactionOutlineFixture, then TransactionOutlineAssertion) {
	return Given(t), Then(t)
}

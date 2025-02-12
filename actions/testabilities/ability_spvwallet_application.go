package testabilities

import "testing"

func New(t testing.TB) (given SPVWalletApplicationFixture, then SPVWalletApplicationAssertions) {
	given = Given(t)
	then = Then(t, given)
	return
}

func NewOf(givenSource SPVWalletApplicationFixture, t testing.TB) (given SPVWalletApplicationFixture, then SPVWalletApplicationAssertions) {
	given = givenSource.NewTest(t)
	then = Then(t, given)
	return
}

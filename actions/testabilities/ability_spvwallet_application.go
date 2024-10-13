package testabilities

import "testing"

func New(t testing.TB) (given SPVWalletApplicationFixture, then SPVWalletApplicationAssertions) {
	return Given(t), Then(t)
}

func NewOf(givenSource SPVWalletApplicationFixture, t testing.TB) (given SPVWalletApplicationFixture, then SPVWalletApplicationAssertions) {
	return givenSource.NewTest(t), Then(t)
}

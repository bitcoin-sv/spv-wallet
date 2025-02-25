package testabilities

import "testing"
import "github.com/bitcoin-sv/spv-wallet/actions/testabilities"

func New(t testing.TB) (given IntegrationTestFixtures, when IntegrationTestAction, then IntegrationTestAssertion) {
	appFixture, appAssertions := testabilities.New(t)

	integrationFixture := newFixture(t, appFixture)
	when = newActions(t, integrationFixture)
	then = newAssertions(t, integrationFixture, appAssertions)

	return integrationFixture, when, then
}

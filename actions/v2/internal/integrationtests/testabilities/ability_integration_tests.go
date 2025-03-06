package testabilities

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet/actions/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/testabilities/testmode"
)

// New creates a new integration test fixture with PostgreSQL container by default
func New(t testing.TB) (given IntegrationTestFixtures, when IntegrationTestAction, then IntegrationTestAssertion) {
	setPostgresContainer(t)

	appFixture, appAssertions := testabilities.New(t)
	integrationFixture := newFixture(t, appFixture)
	when = newActions(t, integrationFixture)
	then = newAssertions(t, integrationFixture, appAssertions)

	return integrationFixture, when, then
}

// setPostgresContainer configures the test to use a PostgreSQL container
func setPostgresContainer(t testing.TB) {
	t.Setenv(testmode.EnvDBMode, testmode.PostgresContainerMode)
}

package testabilities

import (
	"github.com/bitcoin-sv/spv-wallet/engine/testabilities/testmode"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/actions/testabilities"
)

// DBOption represents a database configuration option for integration tests
type DBOption func(t testing.TB)

// WithoutCleanup tells the test not to clean up the database after running
func WithoutCleanup() DBOption {
	return func(t testing.TB) {
		testmode.WithoutCleanup(t)
	}
}

// New creates a new integration test fixture with PostgreSQL container by default
func New(t testing.TB) (given IntegrationTestFixtures, when IntegrationTestAction, then IntegrationTestAssertion) {
	setPostgresContainer(t)

	appFixture, appAssertions := testabilities.New(t)
	integrationFixture := newFixture(t, appFixture)
	when = newActions(t, integrationFixture)
	then = newAssertions(t, integrationFixture, appAssertions)

	return integrationFixture, when, then
}

// NewWithOptions creates a new integration test fixture with database options
func NewWithOptions(t testing.TB, options ...DBOption) (given IntegrationTestFixtures, when IntegrationTestAction, then IntegrationTestAssertion) {
	setPostgresContainer(t)

	for _, option := range options {
		option(t)
	}

	appFixture, appAssertions := testabilities.New(t)
	integrationFixture := newFixture(t, appFixture)
	when = newActions(t, integrationFixture)
	then = newAssertions(t, integrationFixture, appAssertions)

	return integrationFixture, when, then
}

// setPostgresContainer configures the test to use a PostgreSQL container
func setPostgresContainer(t testing.TB) {
	container := testmode.GetOrCreatePostgres(t)
	t.Setenv(testmode.EnvDBHost, container.Host)
	t.Setenv(testmode.EnvDBPort, container.Port)
	t.Setenv(testmode.EnvDBMode, "postgres")
}

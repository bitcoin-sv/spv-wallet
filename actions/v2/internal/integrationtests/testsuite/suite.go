package testsuite

import (
	"os"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/actions/v2/internal/integrationtests/testabilities"
	"github.com/bitcoin-sv/spv-wallet/actions/v2/internal/integrationtests/testcontainer"
	testengine "github.com/bitcoin-sv/spv-wallet/engine/testabilities"
)

const (
	// DbmsPostgres is the PostgreSQL database management system
	DbmsPostgres = "postgres"
	// DbmsSqlite is the SQLite database management system
	DbmsSqlite = "sqlite"

	// EnvSQLiteDBPath is the environment variable to specify the SQLite database file path
	EnvSQLiteDBPath = "TEST_SQLITE_DB_PATH"
)

// RunOnAllDBMS runs the given test function on all database management systems
func RunOnAllDBMS(t *testing.T, testFunc func(t *testing.T, dbms string)) {
	t.Run("PostgreSQL", func(t *testing.T) {
		container := testcontainer.NewPostgresContainer(t)
		t.Logf("Using PostgreSQL container at %s:%s", container.Host, container.Port)

		testFunc(t, DbmsPostgres)
	})

	t.Run("SQLite", func(t *testing.T) {
		cleanupSQLiteFile(t)
		testFunc(t, DbmsSqlite)
	})
}

// SetupDBMSTest is a helper function to set up a test with the specified database system
func SetupDBMSTest(t *testing.T, dbms string) (
	given testabilities.IntegrationTestFixtures,
	when testabilities.IntegrationTestAction,
	then testabilities.IntegrationTestAssertion,
	cleanup func(),
) {
	t.Helper()

	given, when, then = testabilities.New(t)

	if dbms == DbmsPostgres {
		container := testcontainer.NewPostgresContainer(t)

		container.CleanDatabase(t)

		cleanup = given.StartedSPVWalletV2(testengine.WithPostgresConfig(
			container.Host,
			container.Port,
			container.User,
			container.Password,
			container.DBName,
		))
	} else {
		var sqliteOpts []testengine.ConfigOpts

		if dbPath := os.Getenv(EnvSQLiteDBPath); dbPath != "" {
			sqliteOpts = append(sqliteOpts, testengine.WithSQLiteFilePath(dbPath))
		}

		cleanup = given.StartedSPVWalletV2(sqliteOpts...)
	}

	return given, when, then, cleanup
}

// cleanupSQLiteFile removes any SQLite file used for testing if configured
func cleanupSQLiteFile(t *testing.T) {
	t.Helper()

	if dbPath := os.Getenv(EnvSQLiteDBPath); dbPath != "" {
		if _, err := os.Stat(dbPath); err == nil {
			if err := os.Remove(dbPath); err != nil {
				t.Logf("Warning: Failed to remove SQLite file %s: %s", dbPath, err)
			}
		}
	}
}

/*
Package testmode provides functions to set special modes for tests,
allowing to use actual Postgres or SQLite file for testing, especially for development purposes.
Important: It should be used only in LOCAL tests.
Calls of SetPostgresMode and SetFileSQLiteMode should not be committed.
*/
package testmode

import (
	"context"
	"os"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	// EnvDBMode is the environment variable to set the test database mode
	EnvDBMode = "TEST_DB_MODE"
	// EnvDBName is the environment variable to set the test database name
	EnvDBName = "TEST_DB_NAME"
	// EnvDBHost is the environment variable to set the test database host
	EnvDBHost = "TEST_DB_HOST"
	// EnvDBPort is the environment variable to set the test database port
	EnvDBPort = "TEST_DB_PORT"

	// PostgresContainerMode is the mode to use a PostgreSQL testcontainer for testing
	PostgresContainerMode = "postgres-container"
	// PostgresMode is the mode to use actual Postgres for testing
	PostgresMode = "postgres"
	// SQLiteFileMode is the mode to use SQLite file for testing
	SQLiteFileMode = "file"

	// DefaultPostgresName is the default database name for PostgreSQL
	DefaultPostgresName = "postgres"
	// DefaultPostgresUser is the default database user for PostgreSQL
	DefaultPostgresUser = "postgres"
	// DefaultPostgresPass is the default database password for PostgreSQL
	DefaultPostgresPass = "postgres"
)

// TestContainer represents a running test container
type TestContainer struct {
	Container testcontainers.Container
	Host      string
	Port      string
}

// StartPostgresContainer starts a PostgreSQL container and returns connection details
func StartPostgresContainer(t testing.TB) *TestContainer {
	t.Helper()

	ctx := context.Background()
	dbName := DefaultPostgresName
	dbUser := DefaultPostgresUser
	dbPassword := DefaultPostgresPass

	postgresContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2),
		),
	)
	if err != nil {
		t.Fatalf("failed to start postgres container: %s", err)
	}

	t.Cleanup(func() {
		if err := testcontainers.TerminateContainer(postgresContainer); err != nil {
			t.Logf("failed to terminate container: %s", err)
		}
	})

	host, err := postgresContainer.Host(ctx)
	if err != nil {
		t.Fatalf("failed to get container host: %s", err)
	}

	port, err := postgresContainer.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatalf("failed to get mapped port: %s", err)
	}

	t.Logf("Started PostgreSQL container at %s:%s", host, port.Port())

	return &TestContainer{
		Container: postgresContainer,
		Host:      host,
		Port:      port.Port(),
	}
}

// DevelopmentOnly_SetPostgresMode sets the test mode to use actual Postgres
func DevelopmentOnly_SetPostgresMode(t testing.TB) {
	t.Helper()
	t.Setenv(EnvDBMode, PostgresMode)
}

// DevelopmentOnly_SetPostgresModeWithName sets the test mode to use actual Postgres and sets the database name
func DevelopmentOnly_SetPostgresModeWithName(t testing.TB, dbName string) {
	t.Helper()
	t.Setenv(EnvDBMode, PostgresMode)
	t.Setenv(EnvDBName, dbName)
}

// DevelopmentOnly_SetFileSQLiteMode sets the test mode to use SQLite file
func DevelopmentOnly_SetFileSQLiteMode(t testing.TB) {
	t.Helper()
	t.Setenv(EnvDBMode, SQLiteFileMode)
}

// CheckPostgresMode checks if the test mode is set to use actual Postgres and returns the database name
func CheckPostgresMode() (ok bool, dbName string) {
	if os.Getenv(EnvDBMode) != PostgresMode {
		return false, ""
	}
	dbName = os.Getenv(EnvDBName)
	if dbName == "" {
		dbName = DefaultPostgresName
	}
	return true, dbName
}

// CheckFileSQLiteMode checks if the test mode is set to use SQLite file
func CheckFileSQLiteMode() bool {
	return os.Getenv(EnvDBMode) == SQLiteFileMode
}

// CheckPostgresContainerMode checks if the test mode is set to use PostgreSQL container
func CheckPostgresContainerMode() bool {
	return os.Getenv(EnvDBMode) == PostgresContainerMode
}

/*
Package testmode provides functions to set special modes for tests,
allowing to use actual Postgres or SQLite file for testing, especially for development purposes.
Important: It should be used only in LOCAL tests.
Calls of SetPostgresMode and SetFileSQLiteMode should not be committed.
*/
package testmode

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
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

	postgresContainer, err := postgres.Run(ctx,
		"postgres:latest",
		postgres.WithDatabase(DefaultPostgresName),
		postgres.WithUsername(DefaultPostgresUser),
		postgres.WithPassword(DefaultPostgresPass),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		t.Fatalf("failed to start postgres container: %s", err)
	}

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

// CleanDatabaseSchema resets the database schema for a fresh test
func CleanDatabaseSchema(t testing.TB, container *TestContainer) {
	t.Helper()

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		container.Host, container.Port, DefaultPostgresUser, DefaultPostgresPass, DefaultPostgresName,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("Failed to connect to database: %s", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			t.Fatalf("Failed to close database connection: %s", err)
		}
	}(db)

	_, err = db.Exec(`DROP SCHEMA IF EXISTS public CASCADE; CREATE SCHEMA public;`)

	if err != nil {
		t.Fatalf("Failed to clean database: %s", err)
	}

	t.Log("Database schema cleaned successfully")
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

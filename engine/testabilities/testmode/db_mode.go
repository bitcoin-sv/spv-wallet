/*
Package testmode provides functions to set special modes for tests,
allowing to use actual Postgres or SQLite file for testing, especially for development purposes.
Important: It should be used only in LOCAL tests.
Calls of SetPostgresMode and SetFileSQLiteMode should not be committed.
*/
package testmode

import (
	"context"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"os"
	"testing"
	"time"
)

const (
	EnvDBMode      = "TEST_DB_MODE"
	EnvDBContainer = "TEST_DB_CONTAINER"
	EnvDBName      = "TEST_DB_NAME"
	EnvDBHost      = "TEST_DB_HOST"
	EnvDBPort      = "TEST_DB_PORT"

	defaultPostgresDBName = "postgres"
)

// PostgresModeBuilder provides a fluent interface for configuring Postgres test mode
type PostgresModeBuilder struct {
	t testing.TB
}

// DevelopmentOnly_SetPostgresMode sets the test mode to use actual Postgres and sets the database name.
func DevelopmentOnly_SetPostgresMode(t testing.TB) *PostgresModeBuilder {
	t.Setenv(EnvDBMode, "postgres")
	return &PostgresModeBuilder{t: t}
}

// WithTestcontainersMode configures PostgreSQL to run in a testcontainer
func (b *PostgresModeBuilder) WithTestcontainersMode() {
	ctx := context.Background()
	container, err := startPostgresContainer(ctx)
	if err != nil {
		b.t.Fatalf("Failed to start PostgreSQL container: %v", err)
	}
	b.t.Setenv(EnvDBHost, container.host)
	b.t.Setenv(EnvDBPort, container.port)
	b.t.Setenv(EnvDBContainer, "true")

	b.t.Cleanup(func() {
		ctx := context.Background()
		if err := container.container.Terminate(ctx); err != nil {
			b.t.Logf("Failed to stop PostgreSQL container: %v", err)
		}
	})
}

// DevelopmentOnly_SetPostgresModeWithName sets the test mode to use actual Postgres and sets the database name.
func DevelopmentOnly_SetPostgresModeWithName(t testing.TB, dbName string) {
	DevelopmentOnly_SetPostgresMode(t)
	t.Setenv(EnvDBName, dbName)
}

// DevelopmentOnly_SetFileSQLiteMode sets the test mode to use SQLite file
func DevelopmentOnly_SetFileSQLiteMode(t testing.TB) {
	t.Setenv(EnvDBMode, "file")
}

// CheckPostgresMode checks if the test mode is set to use actual Postgres and returns the database name.
func CheckPostgresMode() (ok bool, dbName string) {
	if os.Getenv(EnvDBMode) != "postgres" {
		return false, ""
	}
	dbName = os.Getenv(EnvDBName)
	if dbName == "" {
		dbName = defaultPostgresDBName
	}
	return true, dbName
}

// CheckFileSQLiteMode checks if the test mode is set to use SQLite file
func CheckFileSQLiteMode() bool {
	return os.Getenv(EnvDBMode) == "file"
}

// postgresContainer holds the details of a running PostgreSQL container
type postgresContainer struct {
	container testcontainers.Container
	host      string
	port      string
}

// startPostgresContainer starts a PostgreSQL container
func startPostgresContainer(ctx context.Context) (*postgresContainer, error) {
	req := testcontainers.ContainerRequest{
		Image:        "postgres:14",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "postgres",
			"POSTGRES_PASSWORD": "postgres",
			"POSTGRES_DB":       "postgres",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").WithStartupTimeout(60 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start container: %w", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get container host: %w", err)
	}

	mappedPort, err := container.MappedPort(ctx, nat.Port("5432/tcp"))
	if err != nil {
		return nil, fmt.Errorf("failed to get mapped port: %w", err)
	}

	return &postgresContainer{
		container: container,
		host:      host,
		port:      mappedPort.Port(),
	}, nil
}

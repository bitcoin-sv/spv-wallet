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
	"os"
	"sync"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	EnvDBMode      = "TEST_DB_MODE"
	EnvDBName      = "TEST_DB_NAME"
	EnvDBHost      = "TEST_DB_HOST"
	EnvDBPort      = "TEST_DB_PORT"
	EnvSkipCleanup = "TEST_SKIP_CLEANUP"

	defaultPostgresDBName = "postgres"
)

// Global container instance to be reused across tests
var (
	sharedContainer *TestContainer
	containerMutex  sync.Mutex
)

// PostgresModeBuilder provides a fluent interface for configuring Postgres test mode
type PostgresModeBuilder struct {
	t testing.TB
}

// WithTestcontainersMode configures PostgreSQL to run in a testcontainer
func (b *PostgresModeBuilder) WithTestcontainersMode() *PostgresModeBuilder {
	container := GetOrCreatePostgres(b.t)

	b.t.Setenv(EnvDBHost, container.Host)
	b.t.Setenv(EnvDBPort, container.Port)
	b.t.Setenv(EnvDBMode, "postgres")

	return b
}

// WithoutCleanup prevents the container from being cleaned up after tests
func (b *PostgresModeBuilder) WithoutCleanup() *PostgresModeBuilder {
	WithoutCleanup(b.t)
	return b
}

// DevelopmentOnly_SetPostgresMode sets the test mode to use actual Postgres
// This should be used only for development purposes
func DevelopmentOnly_SetPostgresMode(t testing.TB) *PostgresModeBuilder {
	t.Helper()
	t.Setenv(EnvDBMode, "postgres")
	return &PostgresModeBuilder{t: t}
}

// DevelopmentOnly_SetPostgresModeWithName sets the test mode to use actual Postgres with a specific database name
// This should be used only for development purposes
func DevelopmentOnly_SetPostgresModeWithName(t testing.TB, dbName string) *PostgresModeBuilder {
	t.Helper()
	t.Setenv(EnvDBMode, "postgres")
	t.Setenv(EnvDBName, dbName)
	return &PostgresModeBuilder{t: t}
}

// DevelopmentOnly_SetFileSQLiteMode sets the test mode to use SQLite file
// This should be used only for development purposes
func DevelopmentOnly_SetFileSQLiteMode(t testing.TB) {
	t.Helper()
	t.Setenv(EnvDBMode, "file")
}

// CheckPostgresMode checks if the test mode is set to use actual Postgres and returns the database name
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

// TestContainer represents a running test container
type TestContainer struct {
	Container testcontainers.Container
	Host      string
	Port      string
}

// GetOrCreatePostgres returns an existing PostgreSQL container or creates a new one if none exists
// This helps reuse containers between tests
func GetOrCreatePostgres(t testing.TB) *TestContainer {
	t.Helper()

	containerMutex.Lock()
	defer containerMutex.Unlock()

	if sharedContainer != nil {
		t.Log("Reusing existing PostgreSQL container")
		return sharedContainer
	}

	container := startPostgresContainer(t)
	sharedContainer = container

	if !ShouldSkipCleanup() {
		t.Cleanup(func() {
			containerMutex.Lock()
			defer containerMutex.Unlock()

			if sharedContainer != nil {
				t.Log("Cleaning up PostgreSQL container")
				ctx := context.Background()
				if err := sharedContainer.Container.Terminate(ctx); err != nil {
					t.Logf("Error terminating container: %v", err)
				}
				sharedContainer = nil
			}
		})
	}

	return container
}

// startPostgresContainer starts a PostgreSQL container
func startPostgresContainer(t testing.TB) *TestContainer {
	t.Helper()

	skipCleanup := ShouldSkipCleanup()
	if skipCleanup {
		t.Log("Container cleanup will be skipped (EnvSkipCleanup=true)")
	}

	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:14",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "postgres",
			"POSTGRES_PASSWORD": "postgres",
			"POSTGRES_DB":       "postgres",
		},
		WaitingFor: wait.ForAll(
			wait.ForLog("database system is ready to accept connections").WithStartupTimeout(30*time.Second),
			wait.ForSQL("5432/tcp", "postgres", func(host string, port nat.Port) string {
				return fmt.Sprintf("postgres://postgres:postgres@%s:%s/postgres?sslmode=disable", host, port.Port())
			}).WithStartupTimeout(30*time.Second),
		),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("Failed to start PostgreSQL container: %v", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		t.Fatalf("Failed to get container host: %v", err)
	}

	mappedPort, err := container.MappedPort(ctx, nat.Port("5432/tcp"))
	if err != nil {
		t.Fatalf("Failed to get mapped port: %v", err)
	}

	time.Sleep(1 * time.Second)

	t.Logf("Started PostgreSQL container at %s:%s", host, mappedPort.Port())

	tc := &TestContainer{
		Container: container,
		Host:      host,
		Port:      mappedPort.Port(),
	}

	// for debugging purposes: print a message about how to connect to this database
	if skipCleanup {
		t.Logf("POSTGRES KEPT ALIVE - Connect with: psql -h %s -p %s -U postgres", host, mappedPort.Port())
	}

	return tc
}

// ShouldSkipCleanup returns true if cleanup should be skipped
func ShouldSkipCleanup() bool {
	return os.Getenv(EnvSkipCleanup) == "true"
}

// WithoutCleanup sets the environment to skip cleanup
func WithoutCleanup(t testing.TB) {
	t.Helper()
	t.Setenv(EnvSkipCleanup, "true")
	t.Log("Database cleanup will be skipped for this test")
}

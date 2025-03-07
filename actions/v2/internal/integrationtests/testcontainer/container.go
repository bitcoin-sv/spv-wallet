package testcontainer

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"testing"
)

const (
	schemaQuery = `SELECT schema_name
		FROM information_schema.schemata
		WHERE schema_name NOT IN ('pg_catalog', 'information_schema', 'pg_toast', 'pg_temp_1', 'pg_toast_temp_1')
		      AND schema_name NOT LIKE 'pg_temp_%'
		      AND schema_name NOT LIKE 'pg_toast_temp_%'`
	createSchemaQuery = "CREATE SCHEMA IF NOT EXISTS public"
)

const (
	// DefaultUser is the default username for the PostgreSQL container
	DefaultUser = "postgres"
	// DefaultPassword is the default password for the PostgreSQL container
	DefaultPassword = "postgres"
	// DefaultDBName is the default database name for the PostgreSQL container
	DefaultDBName = "postgres"
)

// PostgresTestContainer encapsulates a PostgreSQL container for testing
type PostgresTestContainer struct {
	Container testcontainers.Container
	Host      string
	Port      string
	User      string
	Password  string
	DBName    string
}

var sharedContainer *PostgresTestContainer

// NewPostgresContainer creates and starts a new PostgreSQL container for testing
func NewPostgresContainer(t testing.TB) *PostgresTestContainer {
	t.Helper()

	if sharedContainer != nil {
		return sharedContainer
	}

	ctx := context.Background()

	container, err := postgres.Run(ctx,
		"postgres:latest",
		postgres.WithDatabase(DefaultDBName),
		postgres.WithUsername(DefaultUser),
		postgres.WithPassword(DefaultPassword),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		t.Fatalf("Failed to start postgres container: %s", err)
	}

	t.Cleanup(func() {
		if err := container.Terminate(ctx); err != nil {
			t.Logf("Failed to terminate container: %s", err)
		}

		sharedContainer = nil
	})

	host, err := container.Host(ctx)
	if err != nil {
		t.Fatalf("Failed to get container host: %s", err)
	}

	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatalf("Failed to get mapped port: %s", err)
	}

	sharedContainer = &PostgresTestContainer{
		Container: container,
		Host:      host,
		Port:      port.Port(),
		User:      DefaultUser,
		Password:  DefaultPassword,
		DBName:    DefaultDBName,
	}

	return sharedContainer
}

// CleanDatabase drops all tables in the database to ensure a clean state for each test
func (p *PostgresTestContainer) CleanDatabase(t testing.TB) {
	t.Helper()

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		p.Host, p.Port, p.User, p.Password, p.DBName,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("Failed to connect to database: %s", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			t.Fatalf("Failed to close database: %s", err)
		}
	}(db)

	if err := db.Ping(); err != nil {
		t.Fatalf("Failed to ping database: %s", err)
	}

	rows, err := db.Query(schemaQuery)
	if err != nil {
		t.Fatalf("Failed to get schemas: %s", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			t.Fatalf("Failed to close rows: %s", err)
		}
	}(rows)

	for rows.Next() {
		var schemaName string
		if err := rows.Scan(&schemaName); err != nil {
			t.Fatalf("Failed to scan schema name: %s", err)
		}

		_, err := db.Exec(fmt.Sprintf("DROP SCHEMA IF EXISTS %s CASCADE", schemaName))
		if err != nil {
			t.Fatalf("Failed to drop schema %s: %s", schemaName, err)
		}
	}

	_, err = db.Exec(createSchemaQuery)
	if err != nil {
		t.Fatalf("Failed to recreate public schema: %s", err)
	}
}

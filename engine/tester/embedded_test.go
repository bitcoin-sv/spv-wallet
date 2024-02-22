//go:build database_tests
// +build database_tests

package tester

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	// testDatabasePort1    = 23902
	testDatabaseHost     = "localhost"
	testDatabaseName     = "test"
	testDatabasePassword = "tester-pw"
	testDatabasePort2    = 23903
	testDatabaseUser     = "tester"
	//testMongoVersion     = "4.2.1"
	testMongoVersion = "6.0.4"
)

// TestCreateMongoServer will test the method CreateMongoServer()
func TestCreateMongoServer(t *testing.T) {
	t.Parallel()

	t.Run("valid server", func(t *testing.T) {
		server, err := CreateMongoServer(
			testMongoVersion,
		)
		require.NoError(t, err)
		require.NotNil(t, server)
		server.Stop()
	})
}

/*
@mrz: This has some strange issues re-running and fails inconsistently

// TestCreatePostgresServer will test the method CreatePostgresServer()
func TestCreatePostgresServer(t *testing.T) {
	// t.Parallel() (disabled for now)

	t.Run("valid server", func(t *testing.T) {
		server, err := CreatePostgresServer(
			23902,
		)
		require.NoError(t, err)
		require.NotNil(t, server)
		err = server.Stop()
		require.NoError(t, err)
	})
}
*/

// TestCreateMySQL will test the method CreateMySQL()
func TestCreateMySQL(t *testing.T) {
	t.Parallel()

	t.Run("valid server", func(t *testing.T) {
		server, err := CreateMySQL(
			testDatabaseHost, testDatabaseName, testDatabaseUser,
			testDatabasePassword, testDatabasePort2,
		)
		require.NoError(t, err)
		require.NotNil(t, server)
		err = server.Close()
		require.NoError(t, err)
	})
}

// TestCreateMySQLTestDatabase will test the method CreateMySQLTestDatabase()
func TestCreateMySQLTestDatabase(t *testing.T) {
	t.Parallel()

	t.Run("valid db", func(t *testing.T) {
		db := CreateMySQLTestDatabase(testDatabaseName)
		require.NotNil(t, db)
		assert.Equal(t, testDatabaseName, db.Name())
	})
}

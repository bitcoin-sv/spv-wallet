/*
Package testmode provides functions to set special modes for tests,
allowing to use actual Postgres or SQLite file for testing, especially for development purposes.
Important: It should be used only in LOCAL tests.
Calls of SetPostgresMode and SetFileSQLiteMode should not be committed.
*/
package testmode

import (
	"os"
	"testing"
)

const (
	modeEnvVar = "TEST_DB_MODE"
	nameEnvVar = "TEST_DB_NAME"
)

// SetPostgresMode sets the test mode to use actual Postgres and sets the database name.
func SetPostgresMode(t testing.TB, dbName string) {
	t.Setenv(modeEnvVar, "postgres")
	t.Setenv(nameEnvVar, dbName)
}

// SetFileSQLiteMode sets the test mode to use SQLite file
func SetFileSQLiteMode(t testing.TB) {
	t.Setenv(modeEnvVar, "file")
}

// CheckPostgresMode checks if the test mode is set to use actual Postgres and returns the database name.
func CheckPostgresMode() (ok bool, dbName string) {
	if os.Getenv(modeEnvVar) != "postgres" {
		return false, ""
	}
	dbName = os.Getenv(nameEnvVar)
	if dbName == "" {
		panic("TEST_DB_NAME must be set when TEST_DB_MODE is 'postgres'")
	}
	return true, dbName
}

// CheckFileSQLiteMode checks if the test mode is set to use SQLite file
func CheckFileSQLiteMode() bool {
	return os.Getenv(modeEnvVar) == "file"
}

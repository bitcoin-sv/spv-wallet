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

	defaultPostgresDBName = "postgres"
)

// DevelopmentOnly_SetPostgresMode sets the test mode to use actual Postgres and sets the database name.
func DevelopmentOnly_SetPostgresMode(t testing.TB) {
	t.Setenv(modeEnvVar, "postgres")
}

// DevelopmentOnly_SetPostgresModeWithName sets the test mode to use actual Postgres and sets the database name.
func DevelopmentOnly_SetPostgresModeWithName(t testing.TB, dbName string) {
	DevelopmentOnly_SetPostgresMode(t)
	t.Setenv(nameEnvVar, dbName)
}

// DevelopmentOnly_SetFileSQLiteMode sets the test mode to use SQLite file
func DevelopmentOnly_SetFileSQLiteMode(t testing.TB) {
	t.Setenv(modeEnvVar, "file")
}

// CheckPostgresMode checks if the test mode is set to use actual Postgres and returns the database name.
func CheckPostgresMode() (ok bool, dbName string) {
	if os.Getenv(modeEnvVar) != "postgres" {
		return false, ""
	}
	dbName = os.Getenv(nameEnvVar)
	if dbName == "" {
		dbName = defaultPostgresDBName
	}
	return true, dbName
}

// CheckFileSQLiteMode checks if the test mode is set to use SQLite file
func CheckFileSQLiteMode() bool {
	return os.Getenv(modeEnvVar) == "file"
}

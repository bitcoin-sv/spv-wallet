package datastore

// Engine is the different engines that are supported (database)
type Engine string

// Supported engines (databases)
const (
	Empty      Engine = "empty"
	MongoDB    Engine = "mongodb"
	MySQL      Engine = "mysql"
	PostgreSQL Engine = "postgresql"
	SQLite     Engine = "sqlite"
)

// index creation constants
const (
	Postgres = "postgres"
	JSON     = "JSON"
	JSONB    = "JSONB"
)

// SQLDatabases is the list of supported SQL databases (via GORM)
var SQLDatabases = []Engine{
	MySQL,
	PostgreSQL,
	SQLite,
}

// String is the string version of engine
func (e Engine) String() string {
	return string(e)
}

// IsEmpty will return true if the datastore is not set
func (e Engine) IsEmpty() bool {
	return e == Empty
}

// IsSQLEngine check whether the string already is in the slice
func IsSQLEngine(e Engine) bool {
	for _, b := range SQLDatabases {
		if b == e {
			return true
		}
	}
	return false
}

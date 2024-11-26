package tgorm

import (
	"context"
	"database/sql"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// DBType Supported databases by this package.
// Added here to prevent circular dependencies.
type DBType string

const (
	// PostgreSQL is the PostgreSQL database type.
	PostgreSQL DBType = "postgresql"
	// SQLite is the SQLite database type.
	SQLite DBType = "sqlite"
)

// GormDBForPrintingSQL is creating a gorm.DB instance that can be used to print or check SQL (without connecting to postgres for example).
// This is useful for checking generated SQL statements.
func GormDBForPrintingSQL(dbType DBType) *gorm.DB {
	dialector := chooseDialector(dbType)

	db, err := gorm.Open(dialector, gormConfig())
	if err != nil {
		panic(err)
	}

	return db
}

func chooseDialector(dbType DBType) gorm.Dialector {
	switch dbType {
	case PostgreSQL:
		return postgres.New(postgres.Config{
			Conn: &doNothingConnectionPool{},
		})
	case SQLite:
		return sqlite.Open(":memory:")
	default:
		panic("unsupported database type " + dbType)
	}
}

func gormConfig() *gorm.Config {
	return &gorm.Config{
		DryRun:         true,
		TranslateError: true,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "xapi_",
			SingularTable: false,
		},
	}
}

type doNothingConnectionPool struct {
}

func (m *doNothingConnectionPool) PrepareContext(_ context.Context, _ string) (*sql.Stmt, error) {
	return nil, nil
}

func (m *doNothingConnectionPool) ExecContext(_ context.Context, _ string, _ ...interface{}) (sql.Result, error) {
	return nil, nil
}

func (m *doNothingConnectionPool) QueryContext(_ context.Context, _ string, _ ...interface{}) (*sql.Rows, error) {
	return nil, nil
}

func (m *doNothingConnectionPool) QueryRowContext(_ context.Context, _ string, _ ...interface{}) *sql.Row {
	return nil
}

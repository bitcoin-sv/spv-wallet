package datastore

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
)

// SQL related default settings
// todo: make this configurable for the end-user?
const (
	defaultFieldStringSize    uint = 256             // default size for string fields
	dsnDefault                     = "file::memory:" // DSN for connection (file or memory, default is memory)
	defaultPreparedStatements      = false           // Flag for prepared statements for SQL
)

// openSQLDatabase will open a new SQL database
func openSQLDatabase(optionalLogger glogger.Interface, configs ...*SQLConfig) (db *gorm.DB, err error) {
	// Try to find a source
	var sourceConfig *SQLConfig
	if sourceConfig, configs = getSourceDatabase(configs); sourceConfig == nil {
		return nil, ErrNoSourceFound
	}

	// Not a valid driver?
	if sourceConfig.Driver != PostgreSQL.String() {
		return nil, ErrUnsupportedDriver
	}

	// Switch on driver
	sourceDialector := getDialector(sourceConfig)

	// Create a new source connection
	// todo: make this configurable? (PrepareStmt)
	if db, err = gorm.Open(
		sourceDialector, getGormConfig(
			sourceConfig.TablePrefix, defaultPreparedStatements,
			sourceConfig.Debug, optionalLogger,
		),
	); err != nil {
		return
	}

	// Start the resolver (default is source and replica are the same)
	resolverConfig := dbresolver.Config{
		Policy:   dbresolver.RandomPolicy{},
		Replicas: []gorm.Dialector{sourceDialector},
		Sources:  []gorm.Dialector{sourceDialector},
	}

	// Do we have additional
	if len(configs) > 0 {

		// Clear the existing replica
		resolverConfig.Replicas = nil

		// Loop configs
		for _, config := range configs {

			// Get the dialector
			dialector := getDialector(config)

			// Set based on replica
			if config.Replica {
				resolverConfig.Replicas = append(resolverConfig.Replicas, dialector)
			} else {
				resolverConfig.Sources = append(resolverConfig.Sources, dialector)
			}
		}

		// No replica?
		if len(resolverConfig.Replicas) == 0 {
			resolverConfig.Replicas = append(resolverConfig.Replicas, sourceDialector)
		}
	}

	// Create the register and set the configuration
	//
	// See: https://gorm.io/docs/dbresolver.html
	// var register *dbresolver.DBResolver
	register := new(dbresolver.DBResolver)
	register.Register(resolverConfig)
	if sourceConfig.MaxConnectionIdleTime.String() != emptyTimeDuration {
		register = register.SetConnMaxIdleTime(sourceConfig.MaxConnectionIdleTime)
	}
	if sourceConfig.MaxConnectionTime.String() != emptyTimeDuration {
		register = register.SetConnMaxLifetime(sourceConfig.MaxConnectionTime)
	}
	if sourceConfig.MaxOpenConnections > 0 {
		register = register.SetMaxOpenConns(sourceConfig.MaxOpenConnections)
	}
	if sourceConfig.MaxIdleConnections > 0 {
		register = register.SetMaxIdleConns(sourceConfig.MaxIdleConnections)
	}

	// Use the register
	if err = db.Use(register); err != nil {
		return
	}

	// Return the connection
	return
}

// openSQLiteDatabase will open a SQLite database connection
func openSQLiteDatabase(optionalLogger glogger.Interface, config *SQLiteConfig) (db *gorm.DB, err error) {
	// Check for an existing connection
	var dialector gorm.Dialector
	if config.ExistingConnection != nil {
		dialector = sqlite.Dialector{Conn: config.ExistingConnection}
	} else {
		dialector = sqlite.Open(getDNS(config.DatabasePath, config.Shared))
	}

	/*
		// todo: implement this functionality (name spaced in-memory tables)
		NOTE: https://www.sqlite.org/inmemorydb.html
		If two or more distinct but shareable in-memory databases are needed in a single process, then the mode=memory
		query parameter can be used with a URI filename to create a named in-memory database:
		rc = sqlite3_open("file:memdb1?mode=memory&cache=shared", &db);
	*/

	// Create a new connection
	if db, err = gorm.Open(
		dialector, getGormConfig(
			config.TablePrefix, defaultPreparedStatements,
			config.Debug, optionalLogger,
		),
	); err != nil {
		return
	}

	// Return the connection
	return
}

// getDNS will return the DNS string
func getDNS(databasePath string, shared bool) (dsn string) {
	// Use a file based path?
	if len(databasePath) > 0 {
		dsn = databasePath
	} else { // Default is in-memory
		dsn = dsnDefault
	}

	// Shared?
	if shared {
		dsn += "?cache=shared"
	}
	return
}

// getDialector will return a new gorm.Dialector based on driver
func getDialector(config *SQLConfig) gorm.Dialector {
	return postgreSQLDialector(config)
}

// postgreSQLDialector will return a gorm.Dialector
func postgreSQLDialector(config *SQLConfig) gorm.Dialector {
	// Create the default PostgreSQL configuration
	cfg := postgres.Config{
		// DriverName: "nrpgx",
		PreferSimpleProtocol: true, // turn to TRUE to disable implicit prepared statement usage
		WithoutReturning:     false,
	}

	// Do we have an existing connection
	if config.ExistingConnection != nil {
		cfg.Conn = config.ExistingConnection
	} else {
		cfg.DSN = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
			config.Host, config.User, config.Password, config.Name, config.Port, config.SslMode, config.TimeZone)
	}

	return postgres.New(cfg)
}

// getSourceDatabase will loop all configs and get the first source
//
// todo: this will grab ANY source (create a better way to seed the source database)
func getSourceDatabase(configs []*SQLConfig) (*SQLConfig, []*SQLConfig) {
	for index, config := range configs {
		if !config.Replica {
			if len(configs) > 1 {
				var processed []*SQLConfig
				for i, c := range configs {
					if i != index {
						processed = append(processed, c)
					}
				}
				return configs[index], processed
			}
			return configs[index], nil
		}
	}
	return nil, configs
}

// getGormSessionConfig returns the gorm session config
func getGormSessionConfig(preparedStatement, debug bool, optionalLogger glogger.Interface) *gorm.Session {
	config := &gorm.Session{
		AllowGlobalUpdate:        false,
		CreateBatchSize:          0,
		DisableNestedTransaction: false,
		DryRun:                   false,
		FullSaveAssociations:     false,
		Logger:                   optionalLogger,
		NewDB:                    false,
		NowFunc:                  nil,
		PrepareStmt:              preparedStatement,
		QueryFields:              false,
		SkipDefaultTransaction:   false,
		SkipHooks:                true,
	}

	// Optional logger vs basic
	if optionalLogger == nil {
		logLevel := glogger.Silent
		if debug {
			logLevel = glogger.Info
		}

		config.Logger = glogger.New(
			log.New(os.Stdout, "\r\n ", log.LstdFlags), // io writer
			glogger.Config{
				SlowThreshold:             5 * time.Second, // Slow SQL threshold
				LogLevel:                  logLevel,        // Log level
				IgnoreRecordNotFoundError: true,            // Ignore ErrRecordNotFound error for logger
				Colorful:                  false,           // Disable color
			},
		)
	}

	return config
}

// getGormConfig will return a valid gorm.Config
//
// See: https://gorm.io/docs/gorm_config.html
func getGormConfig(tablePrefix string, preparedStatement, debug bool, optionalLogger glogger.Interface) *gorm.Config {
	// Set the prefix
	if len(tablePrefix) > 0 {
		tablePrefix = tablePrefix + "_"
	}

	// Create the configuration
	config := &gorm.Config{
		AllowGlobalUpdate:                        false,
		ClauseBuilders:                           nil,
		ConnPool:                                 nil,
		CreateBatchSize:                          0,
		Dialector:                                nil,
		DisableAutomaticPing:                     false,
		DisableForeignKeyConstraintWhenMigrating: true,
		DisableNestedTransaction:                 false,
		DryRun:                                   false, // toggle for extreme debugging
		FullSaveAssociations:                     false,
		Logger:                                   optionalLogger,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   tablePrefix, // table name prefix, table for `User` would be `t_users`
			SingularTable: false,       // use singular table name, table for `User` would be `user` with this option enabled
		},
		NowFunc:                nil,
		Plugins:                nil,
		PrepareStmt:            preparedStatement, // default is: false
		QueryFields:            false,
		SkipDefaultTransaction: false,
	}

	// Optional logger vs basic
	if optionalLogger == nil {
		logLevel := glogger.Silent
		if debug {
			logLevel = glogger.Info
		}

		config.Logger = glogger.New(
			log.New(os.Stdout, "\r\n ", log.LstdFlags), // io writer
			glogger.Config{
				SlowThreshold:             5 * time.Second, // Slow SQL threshold
				LogLevel:                  logLevel,        // Log level
				IgnoreRecordNotFoundError: true,            // Ignore ErrRecordNotFound error for logger
				Colorful:                  false,           // Disable color
			},
		)
	}

	return config
}

// closeSQLDatabase will close an SQL connection safely
func closeSQLDatabase(gormDB *gorm.DB) error {
	if gormDB == nil {
		return nil
	}
	sqlDB, err := gormDB.DB()
	if err != nil {
		return spverrors.Wrapf(err, "failed to close the database connection")
	}
	err = sqlDB.Close()
	return spverrors.Wrapf(err, "failed to close the database connection")
}

// sqlDefaults will set the default values if missing
func (s *SQLConfig) sqlDefaults() *SQLConfig {
	// Set the default(s)
	if s.TxTimeout.String() == emptyTimeDuration {
		s.TxTimeout = defaultDatabaseTxTimeout
	}
	if s.MaxConnectionTime.String() == emptyTimeDuration {
		s.MaxConnectionTime = defaultDatabaseMaxTimeout
	}
	if s.MaxConnectionIdleTime.String() == emptyTimeDuration {
		s.MaxConnectionIdleTime = defaultDatabaseMaxIdleTime
	}
	if len(s.Port) == 0 {
		s.Port = defaultPostgreSQLPort
	}
	if len(s.Host) == 0 {
		s.Host = defaultPostgreSQLHost
	}
	if len(s.TimeZone) == 0 {
		s.TimeZone = defaultTimeZone
	}
	if len(s.SslMode) == 0 {
		s.SslMode = defaultPostgreSQLSslMode
	}
	return s
}

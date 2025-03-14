package datastore

import (
	"context"

	zLogger "github.com/mrz1836/go-logger"
	"gorm.io/gorm"
	gLogger "gorm.io/gorm/logger"
)

type (

	// Client is the datastore client (configuration)
	Client struct {
		options *clientOptions
	}

	// clientOptions holds all the configuration for the client
	clientOptions struct {
		db          *gorm.DB                    // Database connection for Read-Only requests (can be same as Write)
		debug       bool                        // Setting for global debugging
		engine      Engine                      // Datastore engine (PostgreSQL, SQLite)
		fields      *fieldConfig                // Configuration for custom fields
		logger      zLogger.GormLoggerInterface // Custom logger interface (standard interface)
		loggerDB    gLogger.Interface           // Custom logger interface (for GORM)
		sqlConfigs  []*SQLConfig                // Configuration for a PostgreSQL datastore
		sqLite      *SQLiteConfig               // Configuration for a SQLite datastore
		tablePrefix string                      // Model table prefix
	}

	// fieldConfig is the configuration for custom fields
	fieldConfig struct {
		arrayFields  []string // Fields that are an array (string, string, string)
		objectFields []string // Fields that are objects/JSON (metadata)
	}
)

// NewClient creates a new client for all Datastore functionality
//
// If no options are given, it will use the defaultClientOptions()
func NewClient(opts ...ClientOps) (ClientInterface, error) {
	// Create a new client with defaults
	client := &Client{options: defaultClientOptions()}

	// Overwrite defaults with any set by user
	for _, opt := range opts {
		opt(client.options)
	}

	// Set logger (if not set already)
	if client.options.logger == nil {
		client.options.logger = zLogger.NewGormLogger(client.IsDebug(), 5)
	}

	// Create GORM logger
	client.options.loggerDB = &DatabaseLogWrapper{client.options.logger}

	// EMPTY! Engine was NOT set and will use the default (file based)
	if client.Engine().IsEmpty() {

		// Use default SQLite
		// Create a SQLite engine config
		opt := WithSQLite(&SQLiteConfig{
			CommonConfig: CommonConfig{
				Debug:       client.options.debug,
				TablePrefix: defaultTablePrefix,
			},
			DatabasePath: defaultSQLiteFileName,
			Shared:       defaultSQLiteSharing,
		})
		opt(client.options)
	}

	// Use different datastore configurations
	var err error
	if client.Engine() == PostgreSQL {
		if client.options.db, err = openSQLDatabase(
			client.options.loggerDB, client.options.sqlConfigs...,
		); err != nil {
			return nil, err
		}
	} else { // SQLite
		if client.options.db, err = openSQLiteDatabase(
			client.options.loggerDB, client.options.sqLite,
		); err != nil {
			return nil, err
		}
	}

	// Return the client
	return client, nil
}

// Close will terminate (close) the datastore and any open connections
func (c *Client) Close() error {

	if err := closeSQLDatabase(c.options.db); err != nil {
		return err
	}
	c.options.db = nil

	c.options.engine = Empty
	return nil
}

// Debug will set the debug flag
func (c *Client) Debug(on bool) {
	c.options.debug = on
}

// DebugLog will display verbose logs
func (c *Client) DebugLog(ctx context.Context, text string) {
	if c.IsDebug() && c.options.logger != nil {
		c.options.logger.Info(ctx, text)
	}
}

// Engine will return the client's engine
func (c *Client) Engine() Engine {
	return c.options.engine
}

// GetTableName will return the full table name for the given model name
func (c *Client) GetTableName(modelName string) string {
	if c.options.tablePrefix != "" {
		return c.options.tablePrefix + "_" + modelName
	}
	return modelName
}

// GetDatabaseName will return the full database name for the given model name
func (c *Client) GetDatabaseName() string {
	if c.Engine() == PostgreSQL {
		return c.options.sqlConfigs[0].Name
	}

	return ""
}

// GetArrayFields will return the array fields
func (c *Client) GetArrayFields() []string {
	return c.options.fields.arrayFields
}

// GetObjectFields will return the object fields
func (c *Client) GetObjectFields() []string {
	return c.options.fields.objectFields
}

// IsDebug will return the debug flag (bool)
func (c *Client) IsDebug() bool {
	return c.options.debug
}

// DB returns ready to use gorm instance
func (c *Client) DB() *gorm.DB {
	return c.options.db
}

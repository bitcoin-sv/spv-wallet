package datastore

import (
	"context"
	"strings"

	zLogger "github.com/mrz1836/go-logger"
	"github.com/newrelic/go-agent/v3/newrelic"
	"go.mongodb.org/mongo-driver/mongo"
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
		autoMigrate     bool                        // Setting for Auto Migration of SQL tables
		db              *gorm.DB                    // Database connection for Read-Only requests (can be same as Write)
		debug           bool                        // Setting for global debugging
		engine          Engine                      // Datastore engine (PostgreSQL, SQLite)
		fields          *fieldConfig                // Configuration for custom fields
		logger          zLogger.GormLoggerInterface // Custom logger interface (standard interface)
		loggerDB        gLogger.Interface           // Custom logger interface (for GORM)
		migratedModels  []string                    // List of models (types) that have been migrated
		migrateModels   []interface{}               // Models for migrations
		mongoDB         *mongo.Database             // Database connection for a MongoDB datastore
		mongoDBConfig   *MongoDBConfig              // Configuration for a MongoDB datastore
		newRelicEnabled bool                        // If NewRelic is enabled (parent application)
		sqlConfigs      []*SQLConfig                // Configuration for a PostgreSQL datastore
		sqLite          *SQLiteConfig               // Configuration for a SQLite datastore
		tablePrefix     string                      // Model table prefix
	}

	// fieldConfig is the configuration for custom fields
	fieldConfig struct {
		arrayFields                   []string                                // Fields that are an array (string, string, string)
		customMongoConditionProcessor func(conditions map[string]interface{}) // Function for processing custom conditions (arrays, objects)
		customMongoIndexer            func() map[string][]mongo.IndexModel    // Function for returning custom mongo indexes
		objectFields                  []string                                // Fields that are objects/JSON (metadata)
	}
)

// NewClient creates a new client for all Datastore functionality
//
// If no options are given, it will use the defaultClientOptions()
// ctx may contain a NewRelic txn (or one will be created)
func NewClient(ctx context.Context, opts ...ClientOps) (ClientInterface, error) {
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

	// Use NewRelic if it's enabled (use existing txn if found on ctx)
	ctx = client.options.getTxnCtx(ctx)

	// If NewRelic is enabled
	txn := newrelic.FromContext(ctx)
	if txn != nil {
		segment := txn.StartSegment("load_datastore")
		segment.AddAttribute("engine", client.Engine().String())
		defer segment.End()
	}

	// Use different datastore configurations
	var err error
	if client.Engine() == PostgreSQL {
		if client.options.db, err = openSQLDatabase(
			client.options.loggerDB, client.options.sqlConfigs...,
		); err != nil {
			return nil, err
		}
	} else if client.Engine() == MongoDB {
		if client.options.mongoDB, err = openMongoDatabase(
			ctx, client.options.mongoDBConfig,
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

	// Auto migrate
	if client.options.autoMigrate && len(client.options.migrateModels) > 0 {
		if err = client.AutoMigrateDatabase(ctx, client.options.migrateModels...); err != nil {
			return nil, err
		}
	}

	// Return the client
	return client, nil
}

// Close will terminate (close) the datastore and any open connections
func (c *Client) Close(ctx context.Context) error {
	if txn := newrelic.FromContext(ctx); txn != nil {
		defer txn.StartSegment("close_datastore").End()
	}

	// Close Mongo
	if c.Engine() == MongoDB {
		if err := c.options.mongoDB.Client().Disconnect(ctx); err != nil {
			return err
		}
		c.options.mongoDB = nil
	} else { // All other SQL database(s)
		if err := closeSQLDatabase(c.options.db); err != nil {
			return err
		}
		c.options.db = nil
	}

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

// GetMongoConditionProcessor will return a custom mongo condition processor if set
func (c *Client) GetMongoConditionProcessor() func(conditions map[string]interface{}) {
	if c.options.fields.customMongoConditionProcessor != nil {
		return c.options.fields.customMongoConditionProcessor
	}
	return nil
}

// GetMongoIndexer will return a custom mongo condition indexer
func (c *Client) GetMongoIndexer() func() map[string][]mongo.IndexModel {
	if c.options.fields.customMongoIndexer != nil {
		return c.options.fields.customMongoIndexer
	}
	return nil
}

// HasMigratedModel will return if the model type has been migrated
func (c *Client) HasMigratedModel(modelType string) bool {
	for _, t := range c.options.migratedModels {
		if strings.EqualFold(t, modelType) {
			return true
		}
	}
	return false
}

// IsDebug will return the debug flag (bool)
func (c *Client) IsDebug() bool {
	return c.options.debug
}

// IsNewRelicEnabled will return if new relic is enabled
func (c *Client) IsNewRelicEnabled() bool {
	return c.options.newRelicEnabled
}

package datastore

import (
	"context"
	"database/sql"

	zLogger "github.com/mrz1836/go-logger"
	"github.com/newrelic/go-agent/v3/newrelic"
	"go.mongodb.org/mongo-driver/mongo"
)

// ClientOps allow functional options to be supplied
// that overwrite default client options.
type ClientOps func(c *clientOptions)

// defaultClientOptions will return an clientOptions struct with the default settings
//
// Useful for starting with the default and then modifying as needed
func defaultClientOptions() *clientOptions {
	// Set the default options
	return &clientOptions{
		autoMigrate: false,
		engine:      Empty,
		fields: &fieldConfig{
			arrayFields:  nil,
			objectFields: []string{metadataField},
		},
		newRelicEnabled: false,
		sqLite: &SQLiteConfig{
			CommonConfig: CommonConfig{
				Debug: false,
			},
		},
	}
}

// getTxnCtx will check for an existing transaction
func (c *clientOptions) getTxnCtx(ctx context.Context) context.Context {
	if c.newRelicEnabled {
		txn := newrelic.FromContext(ctx)
		if txn != nil {
			ctx = newrelic.NewContext(ctx, txn)
		}
	}
	return ctx
}

// WithAutoMigrate will enable auto migrate database mode (given models)
//
// Pointers of structs (IE: &models.Xpub{})
func WithAutoMigrate(migrateModels ...interface{}) ClientOps {
	return func(c *clientOptions) {
		if len(migrateModels) == 0 {
			return
		}
		for index, model := range migrateModels {
			if model != nil {
				c.autoMigrate = true
				// todo: make a function to ensure these are unique models (no duplicates)
				c.migrateModels = append(c.migrateModels, migrateModels[index])
			}
		}
	}
}

// WithDebugging will enable debugging mode
func WithDebugging() ClientOps {
	return func(c *clientOptions) {
		c.debug = true
	}
}

// WithNewRelic will enable the NewRelic wrapper
func WithNewRelic() ClientOps {
	return func(c *clientOptions) {
		c.newRelicEnabled = true
	}
}

// WithSQLite will set the datastore to use SQLite
func WithSQLite(config *SQLiteConfig) ClientOps {
	return func(c *clientOptions) {
		if config == nil {
			return
		}
		c.sqLite = config
		c.sqLite.MaxIdleConnections = maxIdleConnectionsSQLite // @mrz set this for issues connecting to SQLite
		c.engine = SQLite
		c.tablePrefix = config.TablePrefix
		if c.sqLite.Debug {
			c.debug = true
		}
	}
}

// WithSQL will load a datastore using either an SQL database config or existing connection
func WithSQL(engine Engine, configs []*SQLConfig) ClientOps {
	return func(c *clientOptions) {
		// Do not set if engine is wrong
		if engine != MySQL && engine != PostgreSQL {
			return
		}

		// Loop configurations
		for _, config := range configs {

			// Don't add empty configs
			if config == nil {
				continue
			}

			// Set the defaults if using config vs existing connection
			config.Driver = engine.String()
			if config.ExistingConnection == nil {
				c.sqlConfigs = append(c.sqlConfigs, config.sqlDefaults(engine))
			} else {
				c.sqlConfigs = append(c.sqlConfigs, config)
			}
			if config.Debug {
				c.debug = true
			}
			c.tablePrefix = config.TablePrefix
		}

		// Set the engine
		if len(c.sqlConfigs) > 0 {
			c.engine = engine
		}
	}
}

// WithSQLConnection will set the datastore to an existing connection for MySQL or PostgreSQL
func WithSQLConnection(engine Engine, sqlDB *sql.DB, tablePrefix string) ClientOps {
	return func(c *clientOptions) {
		// Do not set if engine is wrong
		if engine != MySQL && engine != PostgreSQL {
			return
		}

		// Do not set if db is nil
		if sqlDB == nil {
			return
		}

		// this was set for mock testing in MySQL
		// failed to initialize database, got error all expectations were already fulfilled,
		// call to Query 'SELECT VERSION()' with args [] was not expected
		skipInitializeWithVersion := false
		if engine == MySQL {
			skipInitializeWithVersion = true
		}

		c.sqlConfigs = []*SQLConfig{{
			CommonConfig: CommonConfig{
				Debug:       c.debug,
				TablePrefix: tablePrefix,
			},
			Driver:                    engine.String(),
			ExistingConnection:        sqlDB,
			SkipInitializeWithVersion: skipInitializeWithVersion,
		}}
		c.engine = engine
		c.tablePrefix = tablePrefix
	}
}

// WithMongo will set the datastore to use MongoDB
func WithMongo(config *MongoDBConfig) ClientOps {
	return func(c *clientOptions) {
		if config == nil {
			return
		}
		c.engine = MongoDB
		c.tablePrefix = config.TablePrefix
		c.mongoDBConfig = config
		if config.Debug {
			c.debug = true
		}
	}
}

// WithMongoConnection will set the datastore to use an existing Mongo database connection
func WithMongoConnection(database *mongo.Database, tablePrefix string) ClientOps {
	return func(c *clientOptions) {
		if database == nil {
			return
		}
		c.engine = MongoDB
		c.tablePrefix = tablePrefix
		c.mongoDBConfig = &MongoDBConfig{
			CommonConfig: CommonConfig{
				Debug:       c.debug,
				TablePrefix: tablePrefix,
			},
			ExistingConnection: database,
		}
	}
}

// WithLogger will set the custom logger interface
func WithLogger(customLogger zLogger.GormLoggerInterface) ClientOps {
	return func(c *clientOptions) {
		if customLogger != nil {
			c.logger = customLogger
		}
	}
}

// WithCustomFields will add custom fields to the datastore
func WithCustomFields(arrayFields []string, objectFields []string) ClientOps {
	return func(c *clientOptions) {
		if len(arrayFields) > 0 {
			for _, field := range arrayFields {
				if !StringInSlice(field, c.fields.arrayFields) {
					c.fields.arrayFields = append(c.fields.arrayFields, field)
				}
			}
		}
		if len(objectFields) > 0 {
			for _, field := range objectFields {
				if !StringInSlice(field, c.fields.objectFields) {
					c.fields.objectFields = append(c.fields.objectFields, field)
				}
			}
		}
	}
}

// WithCustomMongoConditionProcessor will add a custom mongo condition processor function
func WithCustomMongoConditionProcessor(f func(conditions *map[string]interface{})) ClientOps {
	return func(c *clientOptions) {
		if f != nil {
			c.fields.customMongoConditionProcessor = f
		}
	}
}

// WithCustomMongoIndexer will add a custom mongo index function (returns custom mongo indexes)
func WithCustomMongoIndexer(f func() map[string][]mongo.IndexModel) ClientOps {
	return func(c *clientOptions) {
		if f != nil {
			c.fields.customMongoIndexer = f
		}
	}
}

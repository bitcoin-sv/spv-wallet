package datastore

import (
	"context"
	"errors"
	"fmt"

	"github.com/newrelic/go-agent/v3/newrelic"
	"go.mongodb.org/mongo-driver/mongo"
	mongoOptions "go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// AutoMigrateDatabase will detect the engine and migrate as needed
func (c *Client) AutoMigrateDatabase(ctx context.Context, models ...interface{}) error {
	// Gracefully skip if not enabled
	if !c.options.autoMigrate {
		c.DebugLog(ctx, "auto migrate is disabled, skipping...")
		return nil
	}

	// Make sure we have a supported engine
	if c.Engine() != MySQL &&
		c.Engine() != PostgreSQL &&
		c.Engine() != SQLite &&
		c.Engine() != MongoDB {
		return ErrUnsupportedEngine
	}

	// Check the models against previously migrated models
	for _, modelInterface := range models {
		modelType := fmt.Sprintf("%T", modelInterface)
		if c.HasMigratedModel(modelType) {
			return errors.New("model " + modelType + " was already migrated")
		}
		c.options.migratedModels = append(c.options.migratedModels, modelType)
	}

	// Debug logs
	c.DebugLog(ctx, fmt.Sprintf(
		"database migration starting... engine: %s model_count: %d, models: %v",
		c.Engine().String(),
		len(models),
		c.options.migratedModels,
	))

	// Migrate database for Mongo
	if c.Engine() == MongoDB {
		return autoMigrateMongoDatabase(ctx, c.Engine(), c.options, models...)
	}

	// Migrate database for SQL (using GORM)
	return autoMigrateSQLDatabase(ctx, c.Engine(), c.options.db, c.IsDebug(), c.options.loggerDB, models...)
}

// IsAutoMigrate returns whether auto migration is on
func (c *Client) IsAutoMigrate() bool {
	return c.options.autoMigrate
}

// autoMigrateMongoDatabase will start a new database for Mongo
func autoMigrateMongoDatabase(ctx context.Context, _ Engine, options *clientOptions,
	_ ...interface{},
) error {
	var err error

	if options.fields.customMongoIndexer != nil {
		for collectionName, idx := range options.fields.customMongoIndexer() {
			for _, index := range idx {
				if err = createMongoIndex(
					ctx, options, collectionName, false, index,
				); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// createMongoIndex will create a mongo index
func createMongoIndex(ctx context.Context, options *clientOptions, modelName string, withPrefix bool,
	index mongo.IndexModel,
) error {
	collectionName := modelName
	if !withPrefix {
		collectionName = setPrefix(options.mongoDBConfig.TablePrefix, collectionName)
	}
	collection := options.mongoDB.Collection(collectionName)
	_, err := collection.Indexes().CreateOne(
		ctx, index, mongoOptions.CreateIndexes().SetMaxTime(defaultDatabaseCreateIndexTimeout),
	)

	return err
}

// autoMigrateSQLDatabase will attempt to create or update table schema
//
// See: https://gorm.io/docs/migration.html
func autoMigrateSQLDatabase(ctx context.Context, engine Engine, sqlWriteDB *gorm.DB,
	debug bool, optionalLogger logger.Interface, models ...interface{},
) error {
	// Create a segment
	txn := newrelic.FromContext(ctx)
	if txn != nil {
		defer txn.StartSegment("auto_migrate_sql_database").End()
	}

	// Create a session with config settings
	sessionDb := sqlWriteDB.Session(getGormSessionConfig(sqlWriteDB.PrepareStmt, debug, optionalLogger))

	// Run the auto migrate method
	if engine == MySQL {
		return sessionDb.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(models...)
	}

	// PostgreSQL and SQLite
	return sessionDb.AutoMigrate(models...)
}

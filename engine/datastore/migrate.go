package datastore

import (
	"context"
	"errors"
	"fmt"

	"github.com/newrelic/go-agent/v3/newrelic"
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
	if c.Engine() != PostgreSQL &&
		c.Engine() != SQLite {
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

	// Migrate database for SQL (using GORM)
	return autoMigrateSQLDatabase(ctx, c.options.db, c.IsDebug(), c.options.loggerDB, models...)
}

// IsAutoMigrate returns whether auto migration is on
func (c *Client) IsAutoMigrate() bool {
	return c.options.autoMigrate
}

// autoMigrateSQLDatabase will attempt to create or update table schema
//
// See: https://gorm.io/docs/migration.html
func autoMigrateSQLDatabase(ctx context.Context, sqlWriteDB *gorm.DB, debug bool, optionalLogger logger.Interface, models ...interface{}) error {
	// Create a segment
	txn := newrelic.FromContext(ctx)
	if txn != nil {
		defer txn.StartSegment("auto_migrate_sql_database").End()
	}

	// Create a session with config settings
	sessionDb := sqlWriteDB.Session(getGormSessionConfig(sqlWriteDB.PrepareStmt, debug, optionalLogger))

	// PostgreSQL and SQLite
	return sessionDb.AutoMigrate(models...)
}

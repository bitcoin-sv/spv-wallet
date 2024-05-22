package tester

import (
	"database/sql/driver"
	"os"
	"time"

	"github.com/bitcoin-sv/spv-wallet/engine/datastore"
	embeddedPostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/tryvium-travels/memongo"
)

// AnyTime will fill the need for any timestamp field
type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

// AnyGUID will fill the need for any GUID field
type AnyGUID struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyGUID) Match(v driver.Value) bool {
	str, ok := v.(string)
	return ok && len(str) > 0
}

// CreatePostgresServer will create a new Postgresql server
func CreatePostgresServer(port uint32) (*embeddedPostgres.EmbeddedPostgres, error) {
	// Create the new database
	postgres := embeddedPostgres.NewDatabase(embeddedPostgres.DefaultConfig().Port(port))
	if postgres == nil {
		return nil, ErrFailedLoadingPostgresql
	}

	// Start the database
	if err := postgres.Start(); err != nil {
		return nil, err
	}

	// Return the database
	return postgres, nil
}

// CreateMongoServer will create a new mongo server
func CreateMongoServer(version string) (*memongo.Server, error) {
	mongoServer, err := memongo.StartWithOptions(
		&memongo.Options{
			MongoVersion:     version,
			ShouldUseReplica: false,
			DownloadURL:      os.Getenv("SPV_WALLET_MONGODB_DOWNLOAD_URL"),
		},
	)
	if err != nil {
		return nil, err
	}

	return mongoServer, nil
}

// SQLiteTestConfig will return a test-version of SQLite
func SQLiteTestConfig(debug, shared bool) *datastore.SQLiteConfig {
	return &datastore.SQLiteConfig{
		CommonConfig: datastore.CommonConfig{
			Debug:              debug,
			MaxIdleConnections: 1,
			MaxOpenConnections: 1,
			TablePrefix:        RandomTablePrefix(),
		},
		DatabasePath: "",
		Shared:       shared,
	}
}

package internal

import (
	"context"

	"github.com/bitcoin-sv/spv-wallet/config"
	"github.com/bitcoin-sv/spv-wallet/engine/datastore"
	"github.com/bitcoin-sv/spv-wallet/engine/logging"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/database"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/database/repository"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/fee"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/outlines"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/outlines/utxo"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/utils/must"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// Storage is a struct that holds the database connection.
type Storage struct {
	logger zerolog.Logger
	config *config.AppConfig
	db     *gorm.DB
}

// NewStorage creates a new instance of the storage.
func NewStorage(cfg *config.AppConfig, logger zerolog.Logger) *Storage {
	logger = logger.With().Str("subservice", "database").Logger()

	storage := &Storage{
		logger: logger,
		config: cfg,
	}

	storage.initGormDB()

	return storage
}

// CreateRepositories creates a new instance of the repositories.
func (s *Storage) CreateRepositories() *repository.All {
	must.BeTrue(s.db != nil, "Trying to create repositories on closed database connection")
	return repository.NewRepositories(s.db)
}

// CreateUTXOSelector creates a new instance of the UTXO selector.
func (s *Storage) CreateUTXOSelector(feeService *fee.Service) outlines.UTXOSelector {
	must.BeTrue(s.db != nil, "Trying to create utxo selector on closed database connection")
	// TODO: pass feeService instead of simply fee
	feeUnit, err := feeService.GetFeeUnit(context.Background())
	must.HaveNoErrorf(err, "failed to setup fee unit")
	utxoSelector := utxo.NewSelector(s.db, feeUnit)
	return utxoSelector
}

// Start starts the database connection and migrates the database.
func (s *Storage) Start() error {
	must.BeTrue(s.db != nil, "Trying to start on closed database connection")
	err := s.migrateDatabase()
	if err != nil {
		return err
	}
	return nil
}

// Close closes the database connection.
func (s *Storage) Close() error {
	if s.db == nil {
		return nil
	}
	sqlDB, err := s.db.DB()
	if err != nil {
		return spverrors.Wrapf(err, "failed to close the database connection")
	}
	err = sqlDB.Close()
	return spverrors.Wrapf(err, "failed to close the database connection")
}

func (s *Storage) initGormDB() {
	opts := make([]datastore.ClientOps, 0)

	opts = s.configureLogger(opts)

	opts = s.configureSQL(opts)

	store, err := datastore.NewClient(opts...)
	must.HaveNoErrorf(err, "failed to prepare database connection")
	s.db = store.DB()
}

func (s *Storage) configureLogger(opts []datastore.ClientOps) []datastore.ClientOps {
	var datastoreLogger *logging.GormLoggerAdapter
	loggingLevel := s.logger.GetLevel()
	if loggingLevel == zerolog.InfoLevel {
		warnLvlLogger := s.logger.Level(zerolog.WarnLevel)
		datastoreLogger = logging.CreateGormLoggerAdapter(&warnLvlLogger, "datastore")
	} else {
		datastoreLogger = logging.CreateGormLoggerAdapter(&s.logger, "datastore")
	}
	opts = append(opts, datastore.WithLogger(&datastore.DatabaseLogWrapper{GormLoggerInterface: datastoreLogger}))

	if loggingLevel == zerolog.DebugLevel || loggingLevel == zerolog.TraceLevel {
		opts = append(opts, datastore.WithDebugging())
	}
	return opts
}

func (s *Storage) configureSQL(options []datastore.ClientOps) []datastore.ClientOps {
	// Select the datastore
	if s.config.Db.Datastore.Engine == datastore.SQLite {
		tablePrefix := s.config.Db.Datastore.TablePrefix
		if len(s.config.Db.SQLite.TablePrefix) > 0 {
			tablePrefix = s.config.Db.SQLite.TablePrefix
		}
		options = append(options, datastore.WithSQLite(&datastore.SQLiteConfig{
			CommonConfig: datastore.CommonConfig{
				Debug:                 s.config.Db.Datastore.Debug,
				MaxConnectionIdleTime: s.config.Db.SQLite.MaxConnectionIdleTime,
				MaxConnectionTime:     s.config.Db.SQLite.MaxConnectionTime,
				MaxIdleConnections:    s.config.Db.SQLite.MaxIdleConnections,
				MaxOpenConnections:    s.config.Db.SQLite.MaxOpenConnections,
				TablePrefix:           tablePrefix,
			},
			DatabasePath:       s.config.Db.SQLite.DatabasePath, // "" for in memory
			Shared:             s.config.Db.SQLite.Shared,
			ExistingConnection: s.config.Db.SQLite.ExistingConnection,
		}))
	} else if s.config.Db.Datastore.Engine == datastore.PostgreSQL {
		tablePrefix := s.config.Db.Datastore.TablePrefix
		if len(s.config.Db.SQL.TablePrefix) > 0 {
			tablePrefix = s.config.Db.SQL.TablePrefix
		}

		options = append(options, datastore.WithSQL(s.config.Db.Datastore.Engine, []*datastore.SQLConfig{
			{
				CommonConfig: datastore.CommonConfig{
					Debug:                 s.config.Db.Datastore.Debug,
					MaxConnectionIdleTime: s.config.Db.SQL.MaxConnectionIdleTime,
					MaxConnectionTime:     s.config.Db.SQL.MaxConnectionTime,
					MaxIdleConnections:    s.config.Db.SQL.MaxIdleConnections,
					MaxOpenConnections:    s.config.Db.SQL.MaxOpenConnections,
					TablePrefix:           tablePrefix,
				},
				Driver:    s.config.Db.Datastore.Engine.String(),
				Host:      s.config.Db.SQL.Host,
				Name:      s.config.Db.SQL.Name,
				Password:  s.config.Db.SQL.Password,
				Port:      s.config.Db.SQL.Port,
				TimeZone:  s.config.Db.SQL.TimeZone,
				TxTimeout: s.config.Db.SQL.TxTimeout,
				User:      s.config.Db.SQL.User,
				SslMode:   s.config.Db.SQL.SslMode,
			},
		}))

	} else {
		panic(spverrors.Newf("invalid configuration: unsupported datastore engine: %s", s.config.Db.Datastore.Engine.String()))
	}

	return options
}

func (s *Storage) migrateDatabase() error {
	models := database.Models()
	if err := s.db.AutoMigrate(models...); err != nil {
		return spverrors.Wrapf(err, "failed to auto-migrate database")
	}
	return nil
}

// DB returns the database connection.
// Deprecated: used as adapter for engine v1
func (s *Storage) DB() *gorm.DB {
	return s.db
}

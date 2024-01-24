package config

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/BuxOrg/bux"
	"github.com/BuxOrg/bux-server/logging"
	"github.com/BuxOrg/bux/cluster"
	"github.com/BuxOrg/bux/taskmanager"
	"github.com/BuxOrg/bux/utils"
	broadcastclient "github.com/bitcoin-sv/go-broadcast-client/broadcast/broadcast-client"
	"github.com/go-redis/redis/v8"
	"github.com/mrz1836/go-cachestore"
	"github.com/mrz1836/go-datastore"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rs/zerolog"
)

// AppServices is the loaded services via config
type (
	AppServices struct {
		Bux      bux.ClientInterface
		NewRelic *newrelic.Application
		Logger   *zerolog.Logger
	}
)

// LoadServices will load and return new set of services, updating the AppConfig
func (a *AppConfig) LoadServices(ctx context.Context) (*AppServices, error) {
	// Start services
	_services := new(AppServices)
	var err error

	// Load NewRelic first - used for Application debugging & tracking
	if err = a.loadNewRelic(_services); err != nil {
		return nil, fmt.Errorf("error with loadNewRelic: " + err.Error())
	}

	// Start the NewRelic Tx
	txn := _services.NewRelic.StartTransaction("services_load")
	ctx = newrelic.NewContext(ctx, txn)
	defer txn.End()

	logger, err := logging.CreateLogger(a.Logging.InstanceName, a.Logging.Format, a.Logging.Level, a.Logging.LogOrigin)
	if err != nil {
		return nil, err
	}

	_services.Logger = logger

	// Load BUX
	if err = _services.loadBux(ctx, a, false, logger); err != nil {
		return nil, err
	}

	// Return the services
	return _services, nil
}

// LoadTestServices will load the "minimum" for testing
func (a *AppConfig) LoadTestServices(ctx context.Context) (*AppServices, error) {
	// Start services
	_services := new(AppServices)

	// Load New Relic
	err := a.loadNewRelic(_services)
	if err != nil {
		return nil, err
	}

	// Start the NewRelic Tx
	txn := _services.NewRelic.StartTransaction("services_load_test")
	defer txn.End()

	// Load bux for testing
	if err = _services.loadBux(ctx, a, true, _services.Logger); err != nil {
		return nil, err
	}

	// Return the services
	return _services, nil
}

// loadNewRelic will load New Relic for monitoring
func (a *AppConfig) loadNewRelic(services *AppServices) (err error) {
	// Load new relic
	services.NewRelic, err = newrelic.NewApplication(
		// newrelic.ConfigInfoLogger(os.Stdout),
		// newrelic.ConfigDebugLogger(os.Stdout),
		func(config *newrelic.Config) {
			config.AppName = ApplicationName + "-" + Version
			config.CustomInsightsEvents.Enabled = a.NewRelic.Enabled
			config.DistributedTracer.Enabled = true
			config.Enabled = a.NewRelic.Enabled
			config.HostDisplayName = ApplicationName + "." + Version + "." + a.NewRelic.DomainName
			config.License = a.NewRelic.LicenseKey
			config.TransactionEvents.Enabled = a.NewRelic.Enabled
		},
	)

	// If enabled
	if a.NewRelic.Enabled {
		err = services.NewRelic.WaitForConnection(5 * time.Second)
	}

	return
}

// CloseAll will close all connections to all services
func (s *AppServices) CloseAll(ctx context.Context) {
	// Close Bux
	if s.Bux != nil {
		_ = s.Bux.Close(ctx)
		s.Bux = nil
	}

	// Close new relic
	if s.NewRelic != nil {
		s.NewRelic.Shutdown(DefaultNewRelicShutdown)
		s.NewRelic = nil
	}

	// All services closed!
	if s.Logger != nil {
		s.Logger.Debug().Msg("all services have been closed")
	}
}

// loadBux will load the bux client (including CacheStore and DataStore)
func (s *AppServices) loadBux(ctx context.Context, appConfig *AppConfig, testMode bool, logger *zerolog.Logger) (err error) {
	var options []bux.ClientOps

	if appConfig.NewRelic.Enabled {
		options = append(options, bux.WithNewRelic(s.NewRelic))
	}

	options = append(options, bux.WithUserAgent(appConfig.GetUserAgent()))

	if appConfig.DisableITC {
		options = append(options, bux.WithITCDisabled())
	}

	if appConfig.ImportBlockHeaders != "" {
		options = append(options, bux.WithImportBlockHeaders(appConfig.ImportBlockHeaders))
	}

	if logger != nil {
		buxLogger := logger.With().Str("service", "bux").Logger()
		options = append(options, bux.WithLogger(&buxLogger))
	}

	if appConfig.Debug {
		options = append(options, bux.WithDebugging())
	}

	options = loadCachestore(appConfig, options)

	if options, err = loadCluster(appConfig, options); err != nil {
		return err
	}

	// Set the datastore
	if options, err = loadDatastore(options, appConfig, testMode); err != nil {
		return err
	}

	options = loadPaymail(appConfig, options)

	options = loadTaskManager(appConfig, options)

	if appConfig.Notifications != nil && appConfig.Notifications.Enabled {
		options = append(options, bux.WithNotifications(appConfig.Notifications.WebhookEndpoint))
	}

	if appConfig.Nodes.Protocol == NodesProtocolMapi {
		options = loadMinercraftMapi(appConfig, options)
	} else if appConfig.Nodes.Protocol == NodesProtocolArc {
		options = loadBroadcastClientArc(appConfig, options, logger)
	}

	options = append(options, bux.WithFeeQuotes(appConfig.Nodes.UseFeeQuotes))

	if appConfig.Nodes.FeeUnit != nil {
		options = append(options, bux.WithFeeUnit(&utils.FeeUnit{
			Satoshis: appConfig.Nodes.FeeUnit.Satoshis,
			Bytes:    appConfig.Nodes.FeeUnit.Bytes,
		}))
	}

	// Create the new client
	s.Bux, err = bux.NewClient(ctx, options...)

	return
}

func loadCachestore(appConfig *AppConfig, options []bux.ClientOps) []bux.ClientOps {
	if appConfig.Cache.Engine == cachestore.Redis {
		options = append(options, bux.WithRedis(&cachestore.RedisConfig{
			DependencyMode:        appConfig.Cache.Redis.DependencyMode,
			MaxActiveConnections:  appConfig.Cache.Redis.MaxActiveConnections,
			MaxConnectionLifetime: appConfig.Cache.Redis.MaxConnectionLifetime,
			MaxIdleConnections:    appConfig.Cache.Redis.MaxIdleConnections,
			MaxIdleTimeout:        appConfig.Cache.Redis.MaxIdleTimeout,
			URL:                   appConfig.Cache.Redis.URL,
			UseTLS:                appConfig.Cache.Redis.UseTLS,
		}))
	} else if appConfig.Cache.Engine == cachestore.FreeCache {
		options = append(options, bux.WithFreeCache())
	}

	return options
}

func loadCluster(appConfig *AppConfig, options []bux.ClientOps) ([]bux.ClientOps, error) {
	if appConfig.Cache.Cluster != nil {
		if appConfig.Cache.Cluster.Coordinator == cluster.CoordinatorRedis {
			var redisOptions *redis.Options

			if appConfig.Cache.Cluster.Redis != nil {
				redisURL, err := url.Parse(appConfig.Cache.Cluster.Redis.URL)
				if err != nil {
					return options, fmt.Errorf("error parsing redis url: %w", err)
				}
				password, _ := redisURL.User.Password()
				redisOptions = &redis.Options{
					Addr:        fmt.Sprintf("%s:%s", redisURL.Hostname(), redisURL.Port()),
					Username:    redisURL.User.Username(),
					Password:    password,
					IdleTimeout: appConfig.Cache.Cluster.Redis.MaxIdleTimeout,
				}
				if appConfig.Cache.Cluster.Redis.UseTLS {
					redisOptions.TLSConfig = &tls.Config{
						MinVersion: tls.VersionTLS12,
					}
				}
			} else if appConfig.Cache.Redis != nil {
				redisOptions = &redis.Options{
					Addr:        appConfig.Cache.Redis.URL,
					IdleTimeout: appConfig.Cache.Redis.MaxIdleTimeout,
				}
				if appConfig.Cache.Redis.UseTLS {
					redisOptions.TLSConfig = &tls.Config{
						MinVersion: tls.VersionTLS12,
					}
				}
			} else {
				return options, errors.New("could not load redis cluster coordinator")
			}
			options = append(options, bux.WithClusterRedis(redisOptions))
		}
		if appConfig.Cache.Cluster.Prefix != "" {
			options = append(options, bux.WithClusterKeyPrefix(appConfig.Cache.Cluster.Prefix))
		}
	}

	return options, nil
}

func loadPaymail(appConfig *AppConfig, options []bux.ClientOps) []bux.ClientOps {
	pm := appConfig.Paymail
	options = append(options, bux.WithPaymailSupport(
		pm.Domains,
		pm.DefaultFromPaymail,
		pm.DefaultNote,
		pm.DomainValidationEnabled,
		pm.SenderValidationEnabled,
	))
	if pm.Beef.enabled() {
		options = append(options, bux.WithPaymailBeefSupport(pm.Beef.PulseHeaderValidationURL, pm.Beef.PulseAuthToken))
	}
	return options
}

// loadDatastore will load the correct datastore based on the engine
func loadDatastore(options []bux.ClientOps, appConfig *AppConfig, testMode bool) ([]bux.ClientOps, error) {
	// Set the datastore options
	if testMode {
		var err error
		// Set the unique table prefix
		if appConfig.Db.SQLite.TablePrefix, err = utils.RandomHex(8); err != nil {
			return options, err
		}

		// Defaults for safe thread testing
		appConfig.Db.SQLite.MaxIdleConnections = 1
		appConfig.Db.SQLite.MaxOpenConnections = 1
	}

	// Select the datastore
	if appConfig.Db.Datastore.Engine == datastore.SQLite {
		tablePrefix := appConfig.Db.Datastore.TablePrefix
		if len(appConfig.Db.SQLite.TablePrefix) > 0 {
			tablePrefix = appConfig.Db.SQLite.TablePrefix
		}
		options = append(options, bux.WithSQLite(&datastore.SQLiteConfig{
			CommonConfig: datastore.CommonConfig{
				Debug:                 appConfig.Db.Datastore.Debug,
				MaxConnectionIdleTime: appConfig.Db.SQLite.MaxConnectionIdleTime,
				MaxConnectionTime:     appConfig.Db.SQLite.MaxConnectionTime,
				MaxIdleConnections:    appConfig.Db.SQLite.MaxIdleConnections,
				MaxOpenConnections:    appConfig.Db.SQLite.MaxOpenConnections,
				TablePrefix:           tablePrefix,
			},
			DatabasePath: appConfig.Db.SQLite.DatabasePath, // "" for in memory
			Shared:       appConfig.Db.SQLite.Shared,
		}))
	} else if appConfig.Db.Datastore.Engine == datastore.MySQL || appConfig.Db.Datastore.Engine == datastore.PostgreSQL {
		tablePrefix := appConfig.Db.Datastore.TablePrefix
		if len(appConfig.Db.SQL.TablePrefix) > 0 {
			tablePrefix = appConfig.Db.SQL.TablePrefix
		}

		options = append(options, bux.WithSQL(appConfig.Db.Datastore.Engine, &datastore.SQLConfig{
			CommonConfig: datastore.CommonConfig{
				Debug:                 appConfig.Db.Datastore.Debug,
				MaxConnectionIdleTime: appConfig.Db.SQL.MaxConnectionIdleTime,
				MaxConnectionTime:     appConfig.Db.SQL.MaxConnectionTime,
				MaxIdleConnections:    appConfig.Db.SQL.MaxIdleConnections,
				MaxOpenConnections:    appConfig.Db.SQL.MaxOpenConnections,
				TablePrefix:           tablePrefix,
			},
			Driver:    appConfig.Db.Datastore.Engine.String(),
			Host:      appConfig.Db.SQL.Host,
			Name:      appConfig.Db.SQL.Name,
			Password:  appConfig.Db.SQL.Password,
			Port:      appConfig.Db.SQL.Port,
			TimeZone:  appConfig.Db.SQL.TimeZone,
			TxTimeout: appConfig.Db.SQL.TxTimeout,
			User:      appConfig.Db.SQL.User,
		}))

	} else if appConfig.Db.Datastore.Engine == datastore.MongoDB {

		debug := appConfig.Db.Datastore.Debug
		tablePrefix := appConfig.Db.Datastore.TablePrefix
		if len(appConfig.Db.Mongo.TablePrefix) > 0 {
			tablePrefix = appConfig.Db.Mongo.TablePrefix
		}
		appConfig.Db.Mongo.Debug = debug
		appConfig.Db.Mongo.TablePrefix = tablePrefix
		options = append(options, bux.WithMongoDB(appConfig.Db.Mongo))
	} else {
		return nil, errors.New("unsupported datastore engine: " + appConfig.Db.Datastore.Engine.String())
	}

	options = append(options, bux.WithAutoMigrate(bux.BaseModels...))

	return options, nil
}

func loadTaskManager(appConfig *AppConfig, options []bux.ClientOps) []bux.ClientOps {
	ops := []taskmanager.TasqOps{}
	if appConfig.TaskManager.Factory == taskmanager.FactoryRedis {
		ops = append(ops, taskmanager.WithRedis(appConfig.Cache.Redis.URL))
	}
	options = append(options, bux.WithTaskqConfig(
		taskmanager.DefaultTaskQConfig(TaskManagerQueueName, ops...),
	))
	return options
}

func loadBroadcastClientArc(appConfig *AppConfig, options []bux.ClientOps, logger *zerolog.Logger) []bux.ClientOps {
	builder := broadcastclient.Builder()
	var bcLogger zerolog.Logger
	if logger == nil {
		bcLogger = zerolog.Nop()
	} else {
		bcLogger = logger.With().Str("service", "broadcast-client").Logger()
	}
	for _, arcClient := range appConfig.Nodes.toBroadcastClientArc() {
		builder.WithArc(*arcClient, &bcLogger)
	}
	broadcastClient := builder.Build()
	options = append(
		options,
		bux.WithBroadcastClient(broadcastClient),
	)
	return options
}

func loadMinercraftMapi(appConfig *AppConfig, options []bux.ClientOps) []bux.ClientOps {
	options = append(
		options,
		bux.WithMAPI(),
		bux.WithMinercraftAPIs(appConfig.Nodes.toMinercraftMapi()),
	)
	return options
}

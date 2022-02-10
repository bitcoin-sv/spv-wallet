package config

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/BuxOrg/bux"
	"github.com/BuxOrg/bux/cachestore"
	"github.com/BuxOrg/bux/datastore"
	"github.com/BuxOrg/bux/taskmanager"
	"github.com/BuxOrg/bux/utils"
	"github.com/dgraph-io/ristretto"
	"github.com/go-redis/redis/v8"
	"github.com/go-resty/resty/v2"
	"github.com/mrz1836/go-logger"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/tonicpow/go-minercraft"
	"github.com/tonicpow/go-paymail/server"
)

// AppServices is the loaded services via config
type (
	AppServices struct {
		HTTPClient *http.Client
		MinerCraft minercraft.ClientInterface
		NewRelic   *newrelic.Application
		Resty      *resty.Client
		Bux        bux.ClientInterface
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

	// Create a new Resty Client
	_services.loadResty(txn)

	// Create an HTTP Client
	_services.loadHTTPClient(txn)

	// Start MinerCraft
	if err = _services.loadMinerCraft(txn); err != nil {
		return nil, fmt.Errorf("loadMinerCraft: %w", err)
	}

	// Load bux
	if err = _services.loadBux(ctx, a); err != nil {
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

	// Create a new Resty Client
	_services.loadResty(txn)

	// Create an HTTP Client
	_services.loadHTTPClient(txn)

	// Start MinerCraft
	if err = _services.loadMinerCraft(txn); err != nil {
		return nil, fmt.Errorf("loadMinerCraft: %w", err)
	}

	// Load bux for testing
	if err = _services.loadTestBux(ctx, a); err != nil {
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
			config.AppName = ApplicationName + "-" + a.Environment
			config.CustomInsightsEvents.Enabled = a.NewRelic.Enabled
			config.DistributedTracer.Enabled = true
			config.Enabled = a.NewRelic.Enabled
			config.HostDisplayName = ApplicationName + "." + a.Environment + "." + a.NewRelic.DomainName
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
	logger.Data(2, logger.DEBUG, "all services have been closed")
}

// loadResty will load a resty client
func (s *AppServices) loadResty(txn *newrelic.Transaction) {
	defer txn.StartSegment("load_resty").End()
	s.Resty = s.NewRestyClient()
}

// NewRestyClient will return a new default resty client
func (s *AppServices) NewRestyClient() (restyClient *resty.Client) {
	restyClient = resty.New()
	restyClient.SetTimeout(DefaultHTTPRequestReadTimeout)
	restyClient.SetRetryCount(2)
	return
}

// loadHTTPClient will load the HTTP client
func (s *AppServices) loadHTTPClient(txn *newrelic.Transaction) {
	defer txn.StartSegment("load_http_client").End()
	s.HTTPClient = s.NewHTTPClient()
}

// NewHTTPClient will return a new default HTTP client
func (s *AppServices) NewHTTPClient() *http.Client {
	return &http.Client{
		Timeout: DefaultHTTPRequestWriteTimeout,
	}
}

// loadMinerCraft loads the MinerCraft service
func (s *AppServices) loadMinerCraft(txn *newrelic.Transaction) (err error) {
	defer txn.StartSegment("load_miner_craft").End()
	s.MinerCraft, err = minercraft.NewClient(
		nil,
		s.HTTPClient,
		nil,
	)
	return
}

// loadBux will load the bux client (including CacheStore and DataStore)
func (s *AppServices) loadBux(ctx context.Context, appConfig *AppConfig) (err error) {
	var options []bux.ClientOps

	// Set new relic if enabled
	if appConfig.NewRelic.Enabled {
		options = append(options, bux.WithNewRelic(s.NewRelic))
	}

	// Customize the outgoing user agent
	options = append(options, bux.WithUserAgent(appConfig.GetUserAgent()))

	// Set if the feature is disabled
	if appConfig.DisableITC {
		options = append(options, bux.WithITCDisabled())
	}

	// todo: feature: override the config from JSON env (side-load your own /envs/custom-config.json

	// Debugging
	if appConfig.Debug {
		options = append(options, bux.WithDebugging())
	}

	// Load cache
	if appConfig.Cachestore.Engine == cachestore.Redis {
		options = append(options, bux.WithRedis(&cachestore.RedisConfig{
			DependencyMode:        appConfig.Redis.DependencyMode,
			MaxActiveConnections:  appConfig.Redis.MaxActiveConnections,
			MaxConnectionLifetime: appConfig.Redis.MaxConnectionLifetime,
			MaxIdleConnections:    appConfig.Redis.MaxIdleConnections,
			MaxIdleTimeout:        appConfig.Redis.MaxIdleTimeout,
			URL:                   appConfig.Redis.URL,
			UseTLS:                appConfig.Redis.UseTLS,
		}))
	} else if appConfig.Cachestore.Engine == cachestore.MCache {
		options = append(options, bux.WithMcache())
	} else if appConfig.Cachestore.Engine == cachestore.Ristretto {
		options = append(options, bux.WithRistretto(&ristretto.Config{
			NumCounters:        appConfig.Ristretto.NumCounters,
			MaxCost:            appConfig.Ristretto.MaxCost,
			BufferItems:        appConfig.Ristretto.BufferItems,
			Metrics:            appConfig.Ristretto.Metrics,
			IgnoreInternalCost: appConfig.Ristretto.IgnoreInternalCost,
		}))
	}

	// Set the datastore
	if options, err = loadDatastore(options, appConfig); err != nil {
		return err
	}

	// Set the Paymail server if enabled
	if appConfig.Paymail.Enabled {

		// Append the server config (run LoadPaymailServer())
		options = append(options, bux.WithPaymailServer(
			nil,
			appConfig.Paymail.DefaultFromPaymail,
			appConfig.Paymail.DefaultNote,
		))
	}

	// Load task manager (redis or taskq)
	// todo: this needs more improvement with redis options etc
	if appConfig.TaskManager.Engine == taskmanager.TaskQ {
		config := taskmanager.DefaultTaskQConfig(appConfig.TaskManager.QueueName)
		if appConfig.TaskManager.Factory == taskmanager.FactoryRedis {
			options = append(
				options,
				bux.WithTaskQUsingRedis(
					config,
					&redis.Options{
						Addr: strings.Replace(appConfig.Redis.URL, "redis://", "", -1),
					},
				))
		} else {
			options = append(options, bux.WithTaskQ(config, appConfig.TaskManager.Factory))
		}
	}

	// Create the new client
	s.Bux, err = bux.NewClient(ctx, options...)

	return
}

// SetPaymailServer will modify the bux client with the Paymail server configuration
func (s *AppServices) SetPaymailServer(appConfig *AppConfig, serviceProvider server.PaymailServiceProvider) (err error) {

	// Set the Paymail server configuration if enabled
	if !appConfig.Paymail.Enabled {
		return
	}

	// Add each domain
	opts := make([]server.ConfigOps, 0)
	for _, domain := range appConfig.Paymail.Domains {
		opts = append(opts, server.WithDomain(domain))
	}

	// If sender validation is enabled
	if appConfig.Paymail.SenderValidationEnabled {
		opts = append(opts, server.WithSenderValidation())
	}

	// Create the paymail server configuration
	var config *server.Configuration
	if config, err = server.NewConfig(
		serviceProvider,
		append(opts, server.WithGenericCapabilities())...,
	); err != nil {
		return
	}

	// Modify the server configuration
	s.Bux.ModifyPaymailConfig(
		config,
		appConfig.Paymail.DefaultFromPaymail,
		appConfig.Paymail.DefaultNote,
	)

	return
}

// loadTestBux will create a bux for testing purposes
func (s *AppServices) loadTestBux(ctx context.Context, appConfig *AppConfig) (err error) {
	var options []bux.ClientOps

	// New relic
	if appConfig.NewRelic.Enabled {
		options = append(options, bux.WithNewRelic(s.NewRelic))
	}

	// Set if the feature is disabled
	if appConfig.DisableITC {
		options = append(options, bux.WithITCDisabled())
	}

	// Custom user agent
	options = append(options, bux.WithUserAgent(appConfig.GetUserAgent()))

	// Use in-memory caching
	options = append(options, bux.WithRistretto(&ristretto.Config{
		NumCounters:        appConfig.Ristretto.NumCounters,
		MaxCost:            appConfig.Ristretto.MaxCost,
		BufferItems:        appConfig.Ristretto.BufferItems,
		Metrics:            appConfig.Ristretto.Metrics,
		IgnoreInternalCost: appConfig.Ristretto.IgnoreInternalCost,
	}))

	// Use in-memory TaskQ
	// todo: read from JSON in buxServer config
	options = append(options, bux.WithTaskQ(
		// todo: use a custom queue name from the test or the appConfig?
		taskmanager.DefaultTaskQConfig(appConfig.Datastore.TablePrefix+"_queue"),
		taskmanager.FactoryMemory,
	))

	// Turn on debugging
	if appConfig.Debug {
		options = append(options, bux.WithDebugging())
	}

	// Set the unique table prefix
	if appConfig.SQLite.TablePrefix, err = utils.RandomHex(8); err != nil {
		return err
	}

	// Defaults for safe thread testing
	appConfig.SQLite.MaxIdleConnections = 1
	appConfig.SQLite.MaxOpenConnections = 1

	// Set the datastore
	if options, err = loadDatastore(options, appConfig); err != nil {
		return err
	}

	// Create the client
	s.Bux, err = bux.NewClient(ctx, options...)

	return
}

// loadDatastore will load the correct datastore based on the engine
func loadDatastore(options []bux.ClientOps, appConfig *AppConfig) ([]bux.ClientOps, error) {

	// Select the datastore
	if appConfig.Datastore.Engine == datastore.SQLite {
		debug := appConfig.Datastore.Debug
		if appConfig.SQLite.Debug {
			debug = appConfig.SQLite.Debug
		}
		tablePrefix := appConfig.Datastore.TablePrefix
		if len(appConfig.SQLite.TablePrefix) > 0 {
			tablePrefix = appConfig.SQLite.TablePrefix
		}
		options = append(options, bux.WithSQLite(&datastore.SQLiteConfig{
			CommonConfig: datastore.CommonConfig{
				Debug:       debug,
				TablePrefix: tablePrefix,
			},
			DatabasePath: appConfig.SQLite.DatabasePath, // "" for in memory
			Shared:       appConfig.SQLite.Shared,
		}))
	} else if appConfig.Datastore.Engine == datastore.MySQL || appConfig.Datastore.Engine == datastore.PostgreSQL {

		debug := appConfig.Datastore.Debug
		if appConfig.SQL.Debug {
			debug = appConfig.SQL.Debug
		}
		tablePrefix := appConfig.Datastore.TablePrefix
		if len(appConfig.SQL.TablePrefix) > 0 {
			tablePrefix = appConfig.SQL.TablePrefix
		}

		options = append(options, bux.WithSQL(appConfig.Datastore.Engine, &datastore.SQLConfig{
			CommonConfig: datastore.CommonConfig{
				Debug:                 debug,
				MaxConnectionIdleTime: appConfig.SQL.MaxConnectionIdleTime,
				MaxConnectionTime:     appConfig.SQL.MaxConnectionTime,
				MaxIdleConnections:    appConfig.SQL.MaxIdleConnections,
				MaxOpenConnections:    appConfig.SQL.MaxOpenConnections,
				TablePrefix:           tablePrefix,
			},
			Driver:    appConfig.Datastore.Engine.String(),
			Host:      appConfig.SQL.Host,
			Name:      appConfig.SQL.Name,
			Password:  appConfig.SQL.Password,
			Port:      appConfig.SQL.Port,
			TimeZone:  appConfig.SQL.TimeZone,
			TxTimeout: appConfig.SQL.TxTimeout,
			User:      appConfig.SQL.User,
		}))

	} else if appConfig.Datastore.Engine == datastore.MongoDB {

		debug := appConfig.Datastore.Debug
		if appConfig.Mongo.Debug {
			debug = appConfig.Mongo.Debug
		}
		tablePrefix := appConfig.Datastore.TablePrefix
		if len(appConfig.Mongo.TablePrefix) > 0 {
			tablePrefix = appConfig.Mongo.TablePrefix
		}
		appConfig.Mongo.Debug = debug
		appConfig.Mongo.TablePrefix = tablePrefix
		options = append(options, bux.WithMongoDB(appConfig.Mongo))
	} else {
		return nil, errors.New("unsupported datastore engine: " + appConfig.Datastore.Engine.String())
	}

	// Add the auto migrate
	if appConfig.Datastore.AutoMigrate {
		options = append(options, bux.WithAutoMigrate(bux.BaseModels...))
	}

	return options, nil
}

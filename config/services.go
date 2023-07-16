package config

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/BuxOrg/bux"
	"github.com/BuxOrg/bux/chainstate"
	"github.com/BuxOrg/bux/cluster"
	"github.com/BuxOrg/bux/taskmanager"
	"github.com/BuxOrg/bux/utils"
	"github.com/go-redis/redis/v8"
	"github.com/mrz1836/go-cachestore"
	"github.com/mrz1836/go-datastore"
	"github.com/mrz1836/go-logger"
	"github.com/newrelic/go-agent/v3/newrelic"
)

// AppServices is the loaded services via config
type (
	AppServices struct {
		Bux      bux.ClientInterface
		NewRelic *newrelic.Application
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

	// Load BUX
	if err = _services.loadBux(ctx, a, false); err != nil {
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
	if err = _services.loadBux(ctx, a, true); err != nil {
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

// loadBux will load the bux client (including CacheStore and DataStore)
func (s *AppServices) loadBux(ctx context.Context, appConfig *AppConfig, testMode bool) (err error) {
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

	// Set if the feature is disabled
	if appConfig.ImportBlockHeaders != "" {
		options = append(options, bux.WithImportBlockHeaders(appConfig.ImportBlockHeaders))
	}

	// todo: customize the logger

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
	} else if appConfig.Cachestore.Engine == cachestore.FreeCache {
		options = append(options, bux.WithFreeCache())
	}

	if appConfig.ClusterConfig != nil {
		if appConfig.ClusterConfig.Coordinator == cluster.CoordinatorRedis {
			var redisURL *url.URL
			redisURL, err = url.Parse(appConfig.ClusterConfig.Redis.URL)
			if err != nil {
				return fmt.Errorf("error parsing redis url: %w", err)
			}

			var redisOptions *redis.Options
			if appConfig.ClusterConfig.Redis != nil {
				// parse redis url
				password, _ := redisURL.User.Password()
				redisOptions = &redis.Options{
					Addr:        fmt.Sprintf("%s:%s", redisURL.Hostname(), redisURL.Port()),
					Username:    redisURL.User.Username(),
					Password:    password,
					IdleTimeout: appConfig.ClusterConfig.Redis.MaxIdleTimeout,
				}
				if appConfig.ClusterConfig.Redis.UseTLS {
					redisOptions.TLSConfig = &tls.Config{
						MinVersion: tls.VersionTLS12,
					}
				}
			} else if appConfig.Redis != nil {
				redisOptions = &redis.Options{
					Addr:        appConfig.Redis.URL,
					IdleTimeout: appConfig.Redis.MaxIdleTimeout,
				}
				if appConfig.Redis.UseTLS {
					redisOptions.TLSConfig = &tls.Config{
						MinVersion: tls.VersionTLS12,
					}
				}
			} else {
				return errors.New("could not load redis cluster coordinator")
			}
			options = append(options, bux.WithClusterRedis(redisOptions))
		}
		if appConfig.ClusterConfig.Prefix != "" {
			options = append(options, bux.WithClusterKeyPrefix(appConfig.ClusterConfig.Prefix))
		}
	}

	// Set the datastore options
	if testMode {
		// Set the unique table prefix
		if appConfig.SQLite.TablePrefix, err = utils.RandomHex(8); err != nil {
			return err
		}

		// Defaults for safe thread testing
		appConfig.SQLite.MaxIdleConnections = 1
		appConfig.SQLite.MaxOpenConnections = 1
	}

	// Set the datastore
	if options, err = loadDatastore(options, appConfig); err != nil {
		return err
	}

	// Set the Paymail server if enabled
	if appConfig.Paymail.Enabled {

		// Append the server config
		options = append(options, bux.WithPaymailSupport(
			appConfig.Paymail.Domains,
			appConfig.Paymail.DefaultFromPaymail,
			appConfig.Paymail.DefaultNote,
			appConfig.Paymail.DomainValidationEnabled,
			appConfig.Paymail.SenderValidationEnabled,
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

	// Load the notifications
	if appConfig.Notifications != nil && appConfig.Notifications.Enabled {
		options = append(options, bux.WithNotifications(appConfig.Notifications.WebhookEndpoint))
	}

	if appConfig.Monitor != nil && appConfig.Monitor.Enabled {
		if appConfig.Monitor.BuxAgentURL == "" {
			err = errors.New("CentrifugeServer is required for monitoring to work")
			return
		}
		options = append(options, bux.WithMonitoring(ctx, &chainstate.MonitorOptions{
			Debug:                       appConfig.Monitor.Debug,
			BuxAgentURL:                 appConfig.Monitor.BuxAgentURL,
			MonitorDays:                 appConfig.Monitor.MonitorDays,
			AuthToken:                   appConfig.Monitor.AuthToken,
			FalsePositiveRate:           appConfig.Monitor.FalsePositiveRate,
			MaxNumberOfDestinations:     appConfig.Monitor.MaxNumberOfDestinations,
			SaveTransactionDestinations: appConfig.Monitor.SaveTransactionDestinations,
			LoadMonitoredDestinations:   appConfig.Monitor.LoadMonitoredDestinations,
			ProcessMempoolOnConnect:     appConfig.Monitor.ProcessMempoolOnConnect,
		}))
	}

	if appConfig.UseMapiFeeQuotes {
		options = append(options, bux.WithMapiFeeQuotes())
	}

	// Overrides mApiConfig
	if appConfig.MAPIConfig != nil {
		options = append(options, bux.WithOverridenMAPIConfig(appConfig.MAPIConfig))
	}

	// Create the new client
	s.Bux, err = bux.NewClient(ctx, options...)

	return
}

// loadDatastore will load the correct datastore based on the engine
func loadDatastore(options []bux.ClientOps, appConfig *AppConfig) ([]bux.ClientOps, error) {

	// Select the datastore
	if appConfig.Datastore.Engine == datastore.SQLite {
		debug := appConfig.Datastore.Debug
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

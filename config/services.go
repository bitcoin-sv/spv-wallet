package config

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strings"

	broadcastclient "github.com/bitcoin-sv/go-broadcast-client/broadcast/broadcast-client"
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/engine/cluster"
	"github.com/bitcoin-sv/spv-wallet/engine/datastore"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/taskmanager"
	"github.com/bitcoin-sv/spv-wallet/engine/utils"
	"github.com/bitcoin-sv/spv-wallet/logging"
	"github.com/bitcoin-sv/spv-wallet/metrics"
	"github.com/go-redis/redis/v8"
	"github.com/mrz1836/go-cachestore"
	"github.com/rs/zerolog"
)

// explicitHTTPURLRegex is a regex pattern to check the callback URL (host)
var explicitHTTPURLRegex = regexp.MustCompile(`^https?://`)

// AppServices is the loaded services via config
type (
	AppServices struct {
		SpvWalletEngine engine.ClientInterface
		Logger          *zerolog.Logger
	}
)

// LoadServices will load and return new set of services, updating the AppConfig
func (a *AppConfig) LoadServices(ctx context.Context) (*AppServices, error) {
	// Start services
	_services := new(AppServices)
	var err error

	logger, err := logging.CreateLogger(a.Logging.InstanceName, a.Logging.Format, a.Logging.Level, a.Logging.LogOrigin)
	if err != nil {
		err = spverrors.Wrapf(err, "error creating logger")
		return nil, err
	}

	_services.Logger = logger

	// Load SPV Wallet
	if err := _services.loadSPVWallet(ctx, a, false, logger); err != nil {
		return nil, err
	}

	// Return the services
	return _services, nil
}

// LoadTestServices will load the "minimum" for testing
func (a *AppConfig) LoadTestServices(ctx context.Context) (*AppServices, error) {
	// Start services
	_services := new(AppServices)

	nopLogger := zerolog.Nop()
	_services.Logger = &nopLogger

	// Load SPV Wallet for testing
	if err := _services.loadSPVWallet(ctx, a, true, _services.Logger); err != nil {
		return nil, err
	}

	// Return the services
	return _services, nil
}

// CloseAll will close all connections to all services
func (s *AppServices) CloseAll(ctx context.Context) {
	// Close SPV Wallet Engine
	if s.SpvWalletEngine != nil {
		_ = s.SpvWalletEngine.Close(ctx)
		s.SpvWalletEngine = nil
	}

	// All services closed!
	if s.Logger != nil {
		s.Logger.Debug().Msg("all services have been closed")
	}
}

// loadSPVWallet will load the SPV Wallet client (including CacheStore and DataStore)
func (s *AppServices) loadSPVWallet(ctx context.Context, appConfig *AppConfig, testMode bool, logger *zerolog.Logger) (err error) {
	var options []engine.ClientOps

	if appConfig.Metrics.Enabled {
		collector := metrics.EnableMetrics()
		options = append(options, engine.WithMetrics(collector))
	}

	options = append(options, engine.WithUserAgent(appConfig.GetUserAgent()))

	if logger != nil {
		serviceLogger := logger.With().Str("service", "spv-wallet").Logger()
		options = append(options, engine.WithLogger(&serviceLogger))
	}

	if appConfig.Debug {
		options = append(options, engine.WithDebugging())
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
		options = append(options, engine.WithNotifications())
	}

	options = loadBroadcastClientArc(appConfig, options, logger)

	options = append(options, engine.WithARC(appConfig.ARC.URL, appConfig.ARC.Token, appConfig.ARC.DeploymentID))

	options, err = configureCallback(options, appConfig)
	if err != nil {
		logger.Err(err).Msg("error while configuring callback")
	}

	options = append(options, engine.WithFeeQuotes(appConfig.ARC.UseFeeQuotes))

	if appConfig.ARC.FeeUnit != nil {
		options = append(options, engine.WithFeeUnit(&utils.FeeUnit{
			Satoshis: appConfig.ARC.FeeUnit.Satoshis,
			Bytes:    appConfig.ARC.FeeUnit.Bytes,
		}))
	}

	// Create the new client
	s.SpvWalletEngine, err = engine.NewClient(ctx, options...)

	return
}

func loadCachestore(appConfig *AppConfig, options []engine.ClientOps) []engine.ClientOps {
	if appConfig.Cache.Engine == cachestore.Redis {
		options = append(options, engine.WithRedis(&cachestore.RedisConfig{
			DependencyMode:        appConfig.Cache.Redis.DependencyMode,
			MaxActiveConnections:  appConfig.Cache.Redis.MaxActiveConnections,
			MaxConnectionLifetime: appConfig.Cache.Redis.MaxConnectionLifetime,
			MaxIdleConnections:    appConfig.Cache.Redis.MaxIdleConnections,
			MaxIdleTimeout:        appConfig.Cache.Redis.MaxIdleTimeout,
			URL:                   appConfig.Cache.Redis.URL,
			UseTLS:                appConfig.Cache.Redis.UseTLS,
		}))
	} else if appConfig.Cache.Engine == cachestore.FreeCache {
		options = append(options, engine.WithFreeCache())
	}

	return options
}

func loadCluster(appConfig *AppConfig, options []engine.ClientOps) ([]engine.ClientOps, error) {
	if appConfig.Cache.Cluster != nil {
		if appConfig.Cache.Cluster.Coordinator == cluster.CoordinatorRedis {
			var redisOptions *redis.Options

			if appConfig.Cache.Cluster.Redis != nil {
				redisURL, err := url.Parse(appConfig.Cache.Cluster.Redis.URL)
				if err != nil {
					return options, spverrors.Wrapf(err, "error parsing redis url")
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
				return options, spverrors.Newf("could not load redis cluster coordinator")
			}
			options = append(options, engine.WithClusterRedis(redisOptions))
		}
		if appConfig.Cache.Cluster.Prefix != "" {
			options = append(options, engine.WithClusterKeyPrefix(appConfig.Cache.Cluster.Prefix))
		}
	}

	return options, nil
}

func loadPaymail(appConfig *AppConfig, options []engine.ClientOps) []engine.ClientOps {
	pm := appConfig.Paymail
	options = append(options, engine.WithPaymailSupport(
		pm.Domains,
		pm.DefaultFromPaymail,
		pm.DomainValidationEnabled,
		pm.SenderValidationEnabled,
	))
	if pm.Beef.enabled() {
		options = append(options, engine.WithPaymailBeefSupport(pm.Beef.BlockHeadersServiceHeaderValidationURL, pm.Beef.BlockHeadersServiceAuthToken))
	}
	if appConfig.ExperimentalFeatures.PikeContactsEnabled {
		options = append(options, engine.WithPaymailPikeContactSupport())
	}
	if appConfig.ExperimentalFeatures.PikePaymentEnabled {
		options = append(options, engine.WithPaymailPikePaymentSupport())
	}

	return options
}

// loadDatastore will load the correct datastore based on the engine
func loadDatastore(options []engine.ClientOps, appConfig *AppConfig, testMode bool) ([]engine.ClientOps, error) {
	// Set the datastore options
	if testMode {
		var err error
		// Set the unique table prefix
		if appConfig.Db.SQLite.TablePrefix, err = utils.RandomHex(8); err != nil {
			err = spverrors.Wrapf(err, "error generating random hex")
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
		options = append(options, engine.WithSQLite(&datastore.SQLiteConfig{
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
	} else if appConfig.Db.Datastore.Engine == datastore.PostgreSQL {
		tablePrefix := appConfig.Db.Datastore.TablePrefix
		if len(appConfig.Db.SQL.TablePrefix) > 0 {
			tablePrefix = appConfig.Db.SQL.TablePrefix
		}

		options = append(options, engine.WithSQL(appConfig.Db.Datastore.Engine, &datastore.SQLConfig{
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
			SslMode:   appConfig.Db.SQL.SslMode,
		}))

	} else {
		return nil, spverrors.Newf("unsupported datastore engine: %s", appConfig.Db.Datastore.Engine.String())
	}

	options = append(options, engine.WithAutoMigrate(engine.BaseModels...))

	return options, nil
}

func loadTaskManager(appConfig *AppConfig, options []engine.ClientOps) []engine.ClientOps {
	ops := []taskmanager.TasqOps{}
	if appConfig.TaskManager.Factory == taskmanager.FactoryRedis {
		ops = append(ops, taskmanager.WithRedis(appConfig.Cache.Redis.URL))
	}
	options = append(options, engine.WithTaskqConfig(
		taskmanager.DefaultTaskQConfig(TaskManagerQueueName, ops...),
	))
	return options
}

func loadBroadcastClientArc(appConfig *AppConfig, options []engine.ClientOps, logger *zerolog.Logger) []engine.ClientOps {
	builder := broadcastclient.Builder()
	var bcLogger zerolog.Logger
	if logger == nil {
		bcLogger = zerolog.Nop()
	} else {
		bcLogger = logger.With().Str("service", "broadcast-client").Logger()
	}
	builder.WithArc(broadcastclient.ArcClientConfig{
		Token:        appConfig.ARC.Token,
		APIUrl:       appConfig.ARC.URL,
		DeploymentID: appConfig.ARC.DeploymentID,
	}, &bcLogger)
	broadcastClient := builder.Build()
	options = append(
		options,
		engine.WithBroadcastClient(broadcastClient),
	)
	return options
}

func configureCallback(options []engine.ClientOps, appConfig *AppConfig) ([]engine.ClientOps, error) {
	if appConfig.ARC.Callback.Enabled {
		if !isValidURL(appConfig.ARC.Callback.Host) {
			return nil, spverrors.Newf("invalid callback host: %s - must be a valid external url - not a localhost", appConfig.ARC.Callback.Host)
		}

		if appConfig.ARC.Callback.Token == "" {
			callbackToken, err := utils.HashAdler32(DefaultAdminXpub)
			if err != nil {
				return nil, spverrors.Wrapf(err, "error while generating callback token")
			}
			appConfig.ARC.Callback.Token = callbackToken
		}

		options = append(options, engine.WithCallback(appConfig.ARC.Callback.Host+BroadcastCallbackRoute, appConfig.ARC.Callback.Token))
	}
	return options, nil
}

func isLocal(hostname string) bool {
	if strings.Contains(hostname, "localhost") {
		return true
	}

	ip := net.ParseIP(hostname)
	if ip != nil {
		_, private10, _ := net.ParseCIDR("10.0.0.0/8")
		_, private172, _ := net.ParseCIDR("172.16.0.0/12")
		_, private192, _ := net.ParseCIDR("192.168.0.0/16")
		_, loopback, _ := net.ParseCIDR("127.0.0.0/8")
		_, linkLocal, _ := net.ParseCIDR("169.254.0.0/16")

		return private10.Contains(ip) || private172.Contains(ip) || private192.Contains(ip) || loopback.Contains(ip) || linkLocal.Contains(ip)
	}

	return false
}

func isValidURL(rawURL string) bool {
	if !explicitHTTPURLRegex.MatchString(rawURL) {
		return false
	}
	u, err := url.Parse(rawURL)
	if err != nil {
		return false
	}

	hostname := u.Hostname()

	return !isLocal(hostname)
}

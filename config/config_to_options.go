package config

import (
	"crypto/tls"
	"fmt"
	"net/url"
	"time"

	broadcastclient "github.com/bitcoin-sv/go-broadcast-client/broadcast/broadcast-client"
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/engine/cluster"
	"github.com/bitcoin-sv/spv-wallet/engine/datastore"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/taskmanager"
	"github.com/bitcoin-sv/spv-wallet/engine/utils"
	"github.com/bitcoin-sv/spv-wallet/metrics"
	"github.com/go-redis/redis/v8"
	"github.com/go-resty/resty/v2"
	"github.com/mrz1836/go-cachestore"
	"github.com/rs/zerolog"
)

// ToEngineOptions converts the AppConfig to a slice of engine.ClientOps that can be used to create a new engine.Client.
func (c *AppConfig) ToEngineOptions(logger zerolog.Logger) (options []engine.ClientOps, err error) {
	options = c.addUserAgentOpts(options)

	options = c.addHttpClientOpts(options)

	options = c.addMetricsOpts(options)

	options = c.addLoggerOpts(logger, options)

	options = c.addDebugOpts(options)

	options = c.addCacheStoreOpts(options)

	if options, err = c.addClusterOpts(options); err != nil {
		return nil, err
	}

	if options, err = c.addDataStoreOpts(options); err != nil {
		return nil, err
	}

	options = c.addPaymailOpts(options)

	options = c.addTaskManagerOpts(options)

	options = c.addNotificationOpts(options)

	options = c.addARCOpts(options)

	options = c.addBroadcastClientOpts(options, logger)

	if options, err = c.addCallbackOpts(options); err != nil {
		return nil, err
	}

	options = c.addFeeQuotes(options)

	return options, nil
}

func (c *AppConfig) addHttpClientOpts(options []engine.ClientOps) []engine.ClientOps {
	client := resty.New()
	client.SetTimeout(20 * time.Second)
	client.SetDebug(c.Logging.Level == zerolog.LevelTraceValue)
	client.SetHeader("User-Agent", c.GetUserAgent())
	return append(options, engine.WithHTTPClient(client))
}

func (c *AppConfig) addFeeQuotes(options []engine.ClientOps) []engine.ClientOps {
	options = append(options, engine.WithFeeQuotes(c.ARC.UseFeeQuotes))

	if c.ARC.FeeUnit != nil {
		options = append(options, engine.WithFeeUnit(&utils.FeeUnit{
			Satoshis: c.ARC.FeeUnit.Satoshis,
			Bytes:    c.ARC.FeeUnit.Bytes,
		}))
	}

	return options
}

func (c *AppConfig) addUserAgentOpts(options []engine.ClientOps) []engine.ClientOps {
	return append(options, engine.WithUserAgent(c.GetUserAgent()))
}

func (c *AppConfig) addLoggerOpts(logger zerolog.Logger, options []engine.ClientOps) []engine.ClientOps {
	serviceLogger := logger.With().Str("service", "spv-wallet").Logger()
	return append(options, engine.WithLogger(&serviceLogger))
}

func (c *AppConfig) addMetricsOpts(options []engine.ClientOps) []engine.ClientOps {
	if c.Metrics.Enabled {
		collector := metrics.EnableMetrics()
		options = append(options, engine.WithMetrics(collector))
	}
	return options
}

func (c *AppConfig) addDebugOpts(options []engine.ClientOps) []engine.ClientOps {
	if c.Logging.Level == zerolog.LevelDebugValue || c.Logging.Level == zerolog.LevelTraceValue {
		options = append(options, engine.WithDebugging())
	}
	return options
}

func (c *AppConfig) addCacheStoreOpts(options []engine.ClientOps) []engine.ClientOps {
	if c.Cache.Engine == cachestore.Redis {
		options = append(options, engine.WithRedis(&cachestore.RedisConfig{
			DependencyMode:        c.Cache.Redis.DependencyMode,
			MaxActiveConnections:  c.Cache.Redis.MaxActiveConnections,
			MaxConnectionLifetime: c.Cache.Redis.MaxConnectionLifetime,
			MaxIdleConnections:    c.Cache.Redis.MaxIdleConnections,
			MaxIdleTimeout:        c.Cache.Redis.MaxIdleTimeout,
			URL:                   c.Cache.Redis.URL,
			UseTLS:                c.Cache.Redis.UseTLS,
		}))
	} else if c.Cache.Engine == cachestore.FreeCache {
		options = append(options, engine.WithFreeCache())
	}

	return options
}

func (c *AppConfig) addClusterOpts(options []engine.ClientOps) ([]engine.ClientOps, error) {
	if c.Cache.Cluster == nil {
		return options, nil
	}
	if c.Cache.Cluster.Coordinator == cluster.CoordinatorRedis {
		var redisOptions *redis.Options

		if c.Cache.Cluster.Redis != nil {
			redisURL, err := url.Parse(c.Cache.Cluster.Redis.URL)
			if err != nil {
				return options, spverrors.Wrapf(err, "error parsing redis url")
			}
			password, _ := redisURL.User.Password()
			redisOptions = &redis.Options{
				Addr:        fmt.Sprintf("%s:%s", redisURL.Hostname(), redisURL.Port()),
				Username:    redisURL.User.Username(),
				Password:    password,
				IdleTimeout: c.Cache.Cluster.Redis.MaxIdleTimeout,
			}
			if c.Cache.Cluster.Redis.UseTLS {
				redisOptions.TLSConfig = &tls.Config{
					MinVersion: tls.VersionTLS12,
				}
			}
		} else if c.Cache.Redis != nil {
			redisOptions = &redis.Options{
				Addr:        c.Cache.Redis.URL,
				IdleTimeout: c.Cache.Redis.MaxIdleTimeout,
			}
			if c.Cache.Redis.UseTLS {
				redisOptions.TLSConfig = &tls.Config{
					MinVersion: tls.VersionTLS12,
				}
			}
		} else {
			return options, spverrors.Newf("could not load redis cluster coordinator")
		}
		options = append(options, engine.WithClusterRedis(redisOptions))
	}
	if c.Cache.Cluster.Prefix != "" {
		options = append(options, engine.WithClusterKeyPrefix(c.Cache.Cluster.Prefix))
	}

	return options, nil
}

func (c *AppConfig) addDataStoreOpts(options []engine.ClientOps) ([]engine.ClientOps, error) {
	// Select the datastore
	if c.Db.Datastore.Engine == datastore.SQLite {
		tablePrefix := c.Db.Datastore.TablePrefix
		if len(c.Db.SQLite.TablePrefix) > 0 {
			tablePrefix = c.Db.SQLite.TablePrefix
		}
		options = append(options, engine.WithSQLite(&datastore.SQLiteConfig{
			CommonConfig: datastore.CommonConfig{
				Debug:                 c.Db.Datastore.Debug,
				MaxConnectionIdleTime: c.Db.SQLite.MaxConnectionIdleTime,
				MaxConnectionTime:     c.Db.SQLite.MaxConnectionTime,
				MaxIdleConnections:    c.Db.SQLite.MaxIdleConnections,
				MaxOpenConnections:    c.Db.SQLite.MaxOpenConnections,
				TablePrefix:           tablePrefix,
			},
			DatabasePath: c.Db.SQLite.DatabasePath, // "" for in memory
			Shared:       c.Db.SQLite.Shared,
		}))
	} else if c.Db.Datastore.Engine == datastore.PostgreSQL {
		tablePrefix := c.Db.Datastore.TablePrefix
		if len(c.Db.SQL.TablePrefix) > 0 {
			tablePrefix = c.Db.SQL.TablePrefix
		}

		options = append(options, engine.WithSQL(c.Db.Datastore.Engine, &datastore.SQLConfig{
			CommonConfig: datastore.CommonConfig{
				Debug:                 c.Db.Datastore.Debug,
				MaxConnectionIdleTime: c.Db.SQL.MaxConnectionIdleTime,
				MaxConnectionTime:     c.Db.SQL.MaxConnectionTime,
				MaxIdleConnections:    c.Db.SQL.MaxIdleConnections,
				MaxOpenConnections:    c.Db.SQL.MaxOpenConnections,
				TablePrefix:           tablePrefix,
			},
			Driver:    c.Db.Datastore.Engine.String(),
			Host:      c.Db.SQL.Host,
			Name:      c.Db.SQL.Name,
			Password:  c.Db.SQL.Password,
			Port:      c.Db.SQL.Port,
			TimeZone:  c.Db.SQL.TimeZone,
			TxTimeout: c.Db.SQL.TxTimeout,
			User:      c.Db.SQL.User,
			SslMode:   c.Db.SQL.SslMode,
		}))

	} else {
		return nil, spverrors.Newf("unsupported datastore engine: %s", c.Db.Datastore.Engine.String())
	}

	options = append(options, engine.WithAutoMigrate(engine.BaseModels...))

	return options, nil
}

func (c *AppConfig) addPaymailOpts(options []engine.ClientOps) []engine.ClientOps {
	pm := c.Paymail
	options = append(options, engine.WithPaymailSupport(
		pm.Domains,
		pm.DefaultFromPaymail,
		pm.DomainValidationEnabled,
		pm.SenderValidationEnabled,
	))
	if pm.Beef.enabled() {
		options = append(options, engine.WithPaymailBeefSupport(pm.Beef.BlockHeadersServiceHeaderValidationURL, pm.Beef.BlockHeadersServiceAuthToken))
	}
	if c.ExperimentalFeatures.PikeContactsEnabled {
		options = append(options, engine.WithPaymailPikeContactSupport())
	}
	if c.ExperimentalFeatures.PikePaymentEnabled {
		options = append(options, engine.WithPaymailPikePaymentSupport())
	}

	return options
}

func (c *AppConfig) addTaskManagerOpts(options []engine.ClientOps) []engine.ClientOps {
	var ops []taskmanager.TasqOps
	if c.TaskManager.Factory == taskmanager.FactoryRedis {
		ops = append(ops, taskmanager.WithRedis(c.Cache.Redis.URL))
	}

	return append(options, engine.WithTaskqConfig(
		taskmanager.DefaultTaskQConfig(TaskManagerQueueName, ops...),
	))
}

func (c *AppConfig) addNotificationOpts(options []engine.ClientOps) []engine.ClientOps {
	if c.Notifications != nil && c.Notifications.Enabled {
		options = append(options, engine.WithNotifications())
	}
	return options
}

func (c *AppConfig) addARCOpts(options []engine.ClientOps) []engine.ClientOps {
	return append(options, engine.WithARC(c.ARC.URL, c.ARC.Token, c.ARC.DeploymentID))
}

func (c *AppConfig) addBroadcastClientOpts(options []engine.ClientOps, logger zerolog.Logger) []engine.ClientOps {
	bcLogger := logger.With().Str("service", "broadcast-client").Logger()

	broadcastClient := broadcastclient.Builder().
		WithArc(broadcastclient.ArcClientConfig{
			Token:        c.ARC.Token,
			APIUrl:       c.ARC.URL,
			DeploymentID: c.ARC.DeploymentID,
		}, &bcLogger).
		Build()

	return append(
		options,
		engine.WithBroadcastClient(broadcastClient),
	)
}

func (c *AppConfig) addCallbackOpts(options []engine.ClientOps) ([]engine.ClientOps, error) {
	if !c.ARC.Callback.Enabled {
		return options, nil
	}

	if c.ARC.Callback.Token == "" {
		callbackToken, err := utils.HashAdler32(DefaultAdminXpub)
		if err != nil {
			return nil, spverrors.Wrapf(err, "error while generating callback token")
		}
		c.ARC.Callback.Token = callbackToken
	}

	options = append(options, engine.WithCallback(c.ARC.Callback.Host+BroadcastCallbackRoute, c.ARC.Callback.Token))
	return options, nil
}

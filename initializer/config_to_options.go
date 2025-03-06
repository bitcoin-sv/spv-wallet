package initializer

import (
	"crypto/tls"
	"fmt"
	"net/url"
	"time"

	"github.com/bitcoin-sv/spv-wallet/config"
	"github.com/bitcoin-sv/spv-wallet/conv"
	"github.com/bitcoin-sv/spv-wallet/engine"
	chainmodels "github.com/bitcoin-sv/spv-wallet/engine/chain/models"
	"github.com/bitcoin-sv/spv-wallet/engine/cluster"
	"github.com/bitcoin-sv/spv-wallet/engine/datastore"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/taskmanager"
	"github.com/bitcoin-sv/spv-wallet/engine/utils"
	"github.com/bitcoin-sv/spv-wallet/metrics"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"github.com/go-redis/redis/v8"
	"github.com/go-resty/resty/v2"
	"github.com/mrz1836/go-cachestore"
	"github.com/rs/zerolog"
)

// ToEngineOptions converts the AppConfig to a slice of engine.ClientOps that can be used to create a new engine.Client.
func ToEngineOptions(c *config.AppConfig, logger zerolog.Logger) (options []engine.ClientOps, err error) {
	options = addAppConfigOpts(c, options)

	options = addUserAgentOpts(c, options)

	options = addHttpClientOpts(c, options)

	options = addMetricsOpts(c, options)

	options = addLoggerOpts(c, logger, options)

	options = addDebugOpts(c, options)

	options = addCacheStoreOpts(c, options)

	if options, err = addClusterOpts(c, options); err != nil {
		return nil, err
	}

	if options, err = addDataStoreOpts(c, options); err != nil {
		return nil, err
	}

	options = addPaymailOpts(c, options)

	options = addTaskManagerOpts(c, options)

	options = addNotificationOpts(c, options)

	if options, err = addARCOpts(c, options); err != nil {
		return nil, err
	}

	options = addBHSOpts(c, options)

	options = addCustomFeeUnit(c, options)

	return options, nil
}

func addAppConfigOpts(c *config.AppConfig, options []engine.ClientOps) []engine.ClientOps {
	return append(options, engine.WithAppConfig(c))
}

func addHttpClientOpts(c *config.AppConfig, options []engine.ClientOps) []engine.ClientOps {
	client := resty.New()
	client.SetTimeout(20 * time.Second)
	client.SetDebug(c.Logging.Level == zerolog.LevelTraceValue)
	client.SetHeader("User-Agent", c.GetUserAgent())
	return append(options, engine.WithHTTPClient(client))
}

func addCustomFeeUnit(c *config.AppConfig, options []engine.ClientOps) []engine.ClientOps {
	if c.CustomFeeUnit != nil {
		satoshis, err := conv.IntToUint64(c.CustomFeeUnit.Satoshis)
		if err != nil {
			panic(spverrors.Wrapf(err, "error converting custom fee unit satoshis"))
		}
		options = append(options, engine.WithCustomFeeUnit(bsv.FeeUnit{
			Satoshis: bsv.Satoshis(satoshis),
			Bytes:    c.CustomFeeUnit.Bytes,
		}))
	}

	return options
}

func addUserAgentOpts(c *config.AppConfig, options []engine.ClientOps) []engine.ClientOps {
	return append(options, engine.WithUserAgent(c.GetUserAgent()))
}

func addLoggerOpts(c *config.AppConfig, logger zerolog.Logger, options []engine.ClientOps) []engine.ClientOps {
	serviceLogger := logger.With().Str("service", "spv-wallet").Logger()
	return append(options, engine.WithLogger(&serviceLogger))
}

func addMetricsOpts(c *config.AppConfig, options []engine.ClientOps) []engine.ClientOps {
	if c.Metrics.Enabled {
		collector := metrics.EnableMetrics()
		options = append(options, engine.WithMetrics(collector))
	}
	return options
}

func addDebugOpts(c *config.AppConfig, options []engine.ClientOps) []engine.ClientOps {
	if c.Logging.Level == zerolog.LevelDebugValue || c.Logging.Level == zerolog.LevelTraceValue {
		options = append(options, engine.WithDebugging())
	}
	return options
}

func addCacheStoreOpts(c *config.AppConfig, options []engine.ClientOps) []engine.ClientOps {
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

func addClusterOpts(c *config.AppConfig, options []engine.ClientOps) ([]engine.ClientOps, error) {
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

func addDataStoreOpts(c *config.AppConfig, options []engine.ClientOps) ([]engine.ClientOps, error) {
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
			DatabasePath:       c.Db.SQLite.DatabasePath, // "" for in memory
			Shared:             c.Db.SQLite.Shared,
			ExistingConnection: c.Db.SQLite.ExistingConnection,
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

	return options, nil
}

func addPaymailOpts(c *config.AppConfig, options []engine.ClientOps) []engine.ClientOps {
	pm := c.Paymail
	options = append(options, engine.WithPaymailSupport(
		pm.Domains,
		pm.DefaultFromPaymail,
		pm.DomainValidationEnabled,
		pm.SenderValidationEnabled,
	))
	if pm.Beef.Enabled() {
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

func addTaskManagerOpts(c *config.AppConfig, options []engine.ClientOps) []engine.ClientOps {
	var ops []taskmanager.TasqOps
	if c.TaskManager.Factory == taskmanager.FactoryRedis {
		ops = append(ops, taskmanager.WithRedis(c.Cache.Redis.URL))
	}

	return append(options, engine.WithTaskqConfig(
		taskmanager.DefaultTaskQConfig(config.TaskManagerQueueName, ops...),
	))
}

func addNotificationOpts(c *config.AppConfig, options []engine.ClientOps) []engine.ClientOps {
	if c.Notifications != nil && c.Notifications.Enabled {
		options = append(options, engine.WithNotifications())
	}
	return options
}

func addARCOpts(c *config.AppConfig, options []engine.ClientOps) ([]engine.ClientOps, error) {
	arcCfg := chainmodels.ARCConfig{
		URL:          c.ARC.URL,
		Token:        c.ARC.Token,
		DeploymentID: c.ARC.DeploymentID,
		WaitFor:      c.ARC.WaitForStatus,
	}

	if c.ARCCallbackEnabled() {
		var err error
		if c.ARC.Callback.Token == "" {
			// This also sets the token to the config reference and, it is used in the callbacktoken_middleware
			if c.ARC.Callback.Token, err = utils.HashAdler32(config.DefaultAdminXpub); err != nil {
				return nil, spverrors.Wrapf(err, "error while generating callback token")
			}
		}
		callbackURL, err := c.ARC.Callback.ShouldGetURL()
		if err != nil {
			return nil, spverrors.Wrapf(err, "error while getting callback url")
		}
		arcCfg.Callback = &chainmodels.ARCCallbackConfig{
			URL:   callbackURL.String(),
			Token: c.ARC.Callback.Token,
		}
	}

	if c.ExperimentalFeatures != nil && c.ExperimentalFeatures.UseJunglebus {
		arcCfg.UseJunglebus = true
	}

	return append(options, engine.WithARC(arcCfg)), nil
}

func addBHSOpts(c *config.AppConfig, options []engine.ClientOps) []engine.ClientOps {
	return append(options, engine.WithBHS(c.BHS.URL, c.BHS.AuthToken))
}

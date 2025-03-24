package engine

import (
	"database/sql"
	"net/url"
	"time"

	"github.com/bitcoin-sv/go-paymail"
	"github.com/bitcoin-sv/go-paymail/server"
	"github.com/bitcoin-sv/spv-wallet/config"
	"github.com/bitcoin-sv/spv-wallet/engine/chain/models"
	"github.com/bitcoin-sv/spv-wallet/engine/cluster"
	"github.com/bitcoin-sv/spv-wallet/engine/datastore"
	"github.com/bitcoin-sv/spv-wallet/engine/logging"
	"github.com/bitcoin-sv/spv-wallet/engine/metrics"
	"github.com/bitcoin-sv/spv-wallet/engine/taskmanager"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"github.com/coocood/freecache"
	"github.com/go-redis/redis/v8"
	"github.com/go-resty/resty/v2"
	"github.com/mrz1836/go-cache"
	"github.com/mrz1836/go-cachestore"
	"github.com/rs/zerolog"
	"github.com/vmihailenco/taskq/v3"
)

// ClientOps allow functional options to be supplied that overwrite default client options.
type ClientOps func(c *clientOptions)

// defaultClientOptions will return an clientOptions struct with the default settings
//
// Useful for starting with the default and then modifying as needed
func defaultClientOptions() *clientOptions {
	defaultLogger := logging.GetDefaultLogger()

	dWarnLogger := defaultLogger.Level(zerolog.WarnLevel)
	datastoreLogger := logging.CreateGormLoggerAdapter(&dWarnLogger, "datastore")
	// Set the default options
	return &clientOptions{
		// By default check input utxos (unless disabled by the user)
		iuc: true,

		cluster: &clusterOptions{
			options: []cluster.ClientOps{},
		},

		// Blank cache config
		cacheStore: &cacheStoreOptions{
			ClientInterface: nil,
			options:         []cachestore.ClientOps{},
		},

		// Blank Datastore config
		dataStore: &dataStoreOptions{
			ClientInterface: nil,
			options:         []datastore.ClientOps{datastore.WithLogger(&datastore.DatabaseLogWrapper{GormLoggerInterface: datastoreLogger})},
		},

		// Default http client
		httpClient: resty.New(),

		// Blank Paymail config
		paymail: &paymailOptions{
			client: nil,
			serverConfig: &PaymailServerOptions{
				options: []server.ConfigOps{},
			},
		},

		// Blank TaskManager config
		taskManager: &taskManagerOptions{
			TaskEngine:        nil,
			cronCustomPeriods: map[string]time.Duration{},
		},

		// Default user agent
		userAgent: defaultUserAgent,
	}
}

// DefaultModelOptions will set any default model options (from Client options->model)
func (c *Client) DefaultModelOptions(opts ...ModelOps) []ModelOps {
	// Set the Client from the spvwalletengine.Client onto the model
	opts = append(opts, WithClient(c))

	// Set the encryption key (if found)
	opts = append(opts, WithEncryptionKey(c.options.encryptionKey))

	// Return the new options
	return opts
}

// -----------------------------------------------------------------
// GENERAL
// -----------------------------------------------------------------

// WithUserAgent will overwrite the default useragent
func WithUserAgent(userAgent string) ClientOps {
	return func(c *clientOptions) {
		if len(userAgent) > 0 {
			c.userAgent = userAgent
		}
	}
}

// WithDebugging will set debugging in any applicable configuration
func WithDebugging() ClientOps {
	return func(c *clientOptions) {
		c.debug = true

		// Enable debugging on other services
		c.cacheStore.options = append(c.cacheStore.options, cachestore.WithDebugging())
		c.dataStore.options = append(c.dataStore.options, datastore.WithDebugging())
	}
}

// WithEncryption will set the encryption key and encrypt values using this key
func WithEncryption(key string) ClientOps {
	return func(c *clientOptions) {
		if len(key) > 0 {
			c.encryptionKey = key
		}
	}
}

// WithIUCDisabled will disable checking the input utxos
func WithIUCDisabled() ClientOps {
	return func(c *clientOptions) {
		c.iuc = false
	}
}

// WithHTTPClient will set the custom http interface
func WithHTTPClient(httpClient *resty.Client) ClientOps {
	return func(c *clientOptions) {
		if httpClient != nil {
			c.httpClient = httpClient
		}
	}
}

// WithLogger will set the custom logger interface
func WithLogger(customLogger *zerolog.Logger) ClientOps {
	return func(c *clientOptions) {
		if customLogger != nil {
			c.logger = customLogger

			// Enable the logger on all SPV Wallet Engine services
			taskManagerLogger := customLogger.With().Str("subservice", "taskManager").Logger()
			c.taskManager.options = append(c.taskManager.options, taskmanager.WithLogger(&taskManagerLogger))

			// Enable the logger on all external services
			var datastoreLogger *logging.GormLoggerAdapter
			if customLogger.GetLevel() == zerolog.InfoLevel {
				warnLvlLogger := customLogger.Level(zerolog.WarnLevel)
				datastoreLogger = logging.CreateGormLoggerAdapter(&warnLvlLogger, "datastore")

			} else {
				datastoreLogger = logging.CreateGormLoggerAdapter(customLogger, "datastore")
			}
			c.dataStore.options = append(c.dataStore.options, datastore.WithLogger(&datastore.DatabaseLogWrapper{GormLoggerInterface: datastoreLogger}))

			cachestoreLogger := logging.CreateGormLoggerAdapter(customLogger, "cachestore")
			c.cacheStore.options = append(c.cacheStore.options, cachestore.WithLogger(cachestoreLogger))
		}
	}
}

// -----------------------------------------------------------------
// METRICS
// -----------------------------------------------------------------

// WithMetrics will set the metrics with a collector interface
func WithMetrics(collector metrics.Collector) ClientOps {
	return func(c *clientOptions) {
		if collector != nil {
			c.metrics = metrics.NewMetrics(collector)
		}
	}
}

// -----------------------------------------------------------------
// CACHESTORE
// -----------------------------------------------------------------

// WithCustomCachestore will set the cachestore
func WithCustomCachestore(cacheStore cachestore.ClientInterface) ClientOps {
	return func(c *clientOptions) {
		if cacheStore != nil {
			c.cacheStore.ClientInterface = cacheStore
		}
	}
}

// WithFreeCache will set the cache client for both Read & Write clients
func WithFreeCache() ClientOps {
	return func(c *clientOptions) {
		c.cacheStore.options = append(c.cacheStore.options, cachestore.WithFreeCache())
	}
}

// WithFreeCacheConnection will set the cache client to an active FreeCache connection
func WithFreeCacheConnection(client *freecache.Cache) ClientOps {
	return func(c *clientOptions) {
		if client != nil {
			c.cacheStore.options = append(
				c.cacheStore.options,
				cachestore.WithFreeCacheConnection(client),
			)
		}
	}
}

// WithRedis will set the redis cache client for both Read & Write clients
//
// This will load new redis connections using the given parameters
func WithRedis(config *cachestore.RedisConfig) ClientOps {
	return func(c *clientOptions) {
		if config != nil {
			c.cacheStore.options = append(c.cacheStore.options, cachestore.WithRedis(config))
		}
	}
}

// WithRedisConnection will set the cache client to an active redis connection
func WithRedisConnection(activeClient *cache.Client) ClientOps {
	return func(c *clientOptions) {
		if activeClient != nil {
			c.cacheStore.options = append(
				c.cacheStore.options,
				cachestore.WithRedisConnection(activeClient),
			)
		}
	}
}

// -----------------------------------------------------------------
// DATASTORE
// -----------------------------------------------------------------

// WithCustomDatastore will set the datastore
func WithCustomDatastore(dataStore datastore.ClientInterface) ClientOps {
	return func(c *clientOptions) {
		if dataStore != nil {
			c.dataStore.ClientInterface = dataStore
		}
	}
}

// WithSQLite will set the Datastore to use SQLite
func WithSQLite(config *datastore.SQLiteConfig) ClientOps {
	return func(c *clientOptions) {
		if config != nil {
			c.dataStore.options = append(c.dataStore.options, datastore.WithSQLite(config))
		}
	}
}

// WithSQL will set the datastore to use the SQL config
func WithSQL(engine datastore.Engine, config *datastore.SQLConfig) ClientOps {
	return func(c *clientOptions) {
		if config != nil && !engine.IsEmpty() {
			c.dataStore.options = append(
				c.dataStore.options,
				datastore.WithSQL(engine, []*datastore.SQLConfig{config}),
			)
		}
	}
}

// WithSQLConnection will set the Datastore to an existing connection for PostgreSQL
func WithSQLConnection(engine datastore.Engine, sqlDB *sql.DB, tablePrefix string) ClientOps {
	return func(c *clientOptions) {
		if sqlDB != nil && !engine.IsEmpty() {
			c.dataStore.options = append(
				c.dataStore.options,
				datastore.WithSQLConnection(engine, sqlDB, tablePrefix),
			)
		}
	}
}

// -----------------------------------------------------------------
// PAYMAIL
// -----------------------------------------------------------------

// WithPaymailClient will set a custom paymail client
func WithPaymailClient(client paymail.ClientInterface) ClientOps {
	return func(c *clientOptions) {
		if client != nil {
			c.paymail.client = client
		}
	}
}

// WithPaymailSupport will set the configuration for Paymail support (as a server)
func WithPaymailSupport(domains []string, defaultFromPaymail string, domainValidation, senderValidation bool) ClientOps {
	return func(c *clientOptions) {
		// Add generic capabilities
		c.paymail.serverConfig.options = append(c.paymail.serverConfig.options, server.WithP2PCapabilities())

		// Add each domain
		for _, domain := range domains {
			c.paymail.serverConfig.options = append(c.paymail.serverConfig.options, server.WithDomain(domain))
		}

		// Set the sender validation
		if senderValidation {
			c.paymail.serverConfig.options = append(c.paymail.serverConfig.options, server.WithSenderValidation())
		}

		// Domain validation
		if !domainValidation {
			c.paymail.serverConfig.options = append(c.paymail.serverConfig.options, server.WithDomainValidationDisabled())
		}

		// Add default values
		if len(defaultFromPaymail) > 0 {
			c.paymail.serverConfig.DefaultFromPaymail = defaultFromPaymail
		}
	}
}

// WithPaymailBeefSupport will enable Paymail BEEF format support (as a server) and create a Block Headers Service client for Merkle Roots verification.
func WithPaymailBeefSupport(blockHeadersServiceURL, blockHeadersServiceAuthToken string) ClientOps {
	return func(c *clientOptions) {
		_, err := url.ParseRequestURI(blockHeadersServiceURL)
		if err != nil {
			panic(err)
		}
		c.paymail.serverConfig.options = append(c.paymail.serverConfig.options, server.WithBeefCapabilities())
	}
}

// WithPaymailPikeContactSupport will enable Paymail Pike Contact support
func WithPaymailPikeContactSupport() ClientOps {
	return func(c *clientOptions) {
		c.paymail.serverConfig.options = append(c.paymail.serverConfig.options, server.WithPikeContactCapabilities())
	}
}

// WithPaymailPikePaymentSupport will enable Paymail Pike Payment support
func WithPaymailPikePaymentSupport() ClientOps {
	return func(c *clientOptions) {
		c.paymail.serverConfig.options = append(c.paymail.serverConfig.options, server.WithPikePaymentCapabilities())
	}
}

// -----------------------------------------------------------------
// TASK MANAGER
// -----------------------------------------------------------------

// WithTaskqConfig will set the task manager to use TaskQ & in-memory
func WithTaskqConfig(config *taskq.QueueOptions) ClientOps {
	return func(c *clientOptions) {
		if config != nil {
			c.taskManager.options = append(
				c.taskManager.options,
				taskmanager.WithTaskqConfig(config),
			)
		}
	}
}

// WithCronCustomPeriod will set the custom cron jobs period which will override the default
func WithCronCustomPeriod(cronJobName string, period time.Duration) ClientOps {
	return func(c *clientOptions) {
		if c.taskManager != nil {
			c.taskManager.cronCustomPeriods[cronJobName] = period
		}
	}
}

// -----------------------------------------------------------------
// CLUSTER
// -----------------------------------------------------------------

// WithClusterRedis will set the cluster coordinator to use redis
func WithClusterRedis(redisOptions *redis.Options) ClientOps {
	return func(c *clientOptions) {
		if redisOptions != nil {
			c.cluster.options = append(c.cluster.options, cluster.WithRedis(redisOptions))
		}
	}
}

// WithClusterKeyPrefix will set the cluster key prefix to use for all keys in the cluster coordinator
func WithClusterKeyPrefix(prefix string) ClientOps {
	return func(c *clientOptions) {
		if prefix != "" {
			c.cluster.options = append(c.cluster.options, cluster.WithKeyPrefix(prefix))
		}
	}
}

// WithClusterClient will set the cluster options on the client
func WithClusterClient(clusterClient cluster.ClientInterface) ClientOps {
	return func(c *clientOptions) {
		if clusterClient != nil {
			c.cluster.ClientInterface = clusterClient
		}
	}
}

// -----------------------------------------------------------------
// NOTIFICATIONS
// -----------------------------------------------------------------

// WithNotifications will set the notifications config
func WithNotifications() ClientOps {
	return func(c *clientOptions) {
		c.notifications = &notificationsOptions{
			enabled: true,
		}
	}
}

// -----------------------------------------------------------------
// CHAIN
// -----------------------------------------------------------------

// WithCustomFeeUnit will set the custom fee unit for transactions
func WithCustomFeeUnit(feeUnit bsv.FeeUnit) ClientOps {
	return func(c *clientOptions) {
		c.feeUnit = &feeUnit
	}
}

// WithARC sets all the ARC options needed for broadcasting, querying transactions etc.
func WithARC(arcCfg chainmodels.ARCConfig) ClientOps {
	return func(c *clientOptions) {
		c.arcConfig = arcCfg
	}
}

// WithBHS set BHS url params
func WithBHS(url, token string) ClientOps {
	return func(c *clientOptions) {
		c.bhsConfig = chainmodels.BHSConfig{
			URL:       url,
			AuthToken: token,
		}
	}
}

// WithAppConfig passes the config struct into engine
func WithAppConfig(config *config.AppConfig) ClientOps {
	return func(c *clientOptions) {
		c.config = config
	}
}

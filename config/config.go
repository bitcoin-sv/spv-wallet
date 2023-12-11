// Package config provides a configuration for the API
package config

import (
	"time"

	"github.com/BuxOrg/bux/cluster"
	"github.com/BuxOrg/bux/taskmanager"
	"github.com/mrz1836/go-cachestore"
	"github.com/mrz1836/go-datastore"
	"github.com/tonicpow/go-minercraft/v2"
)

// Config constants used for optimization and value testing
const (
	ApplicationName         = "BuxServer"
	CurrentMajorVersion     = "v1"
	DefaultNewRelicShutdown = 10 * time.Second
	HealthRequestPath       = "health"
	Version                 = "v0.5.16"
)

// ConfigFilePathKey is the viper key under which a config file path is stored
const ConfigFilePathKey = "config_file"

// AppConfig is the configuration values and associated env vars
type AppConfig struct {
	Authentication     *AuthenticationConfig `json:"auth" mapstructure:"auth"`
	Cache              *CacheConfig          `json:"cache" mapstructure:"cache"`
	Db                 *DbConfig             `json:"db" mapstructure:"db"`
	Debug              bool                  `json:"debug" mapstructure:"debug"`
	DebugProfiling     bool                  `json:"debug_profiling" mapstructure:"debug_profiling"`
	DisableITC         bool                  `json:"disable_itc" mapstructure:"disable_itc"`
	GraphQL            *GraphqlConfig        `json:"graphql" mapstructure:"graphql"`
	ImportBlockHeaders string                `json:"import_block_headers" mapstructure:"import_block_headers"`
	Logging            *LoggingConfig        `json:"logging" mapstructure:"logging"`
	Monitor            *MonitorOptions       `json:"monitor" mapstructure:"monitor"`
	NewRelic           *NewRelicConfig       `json:"new_relic" mapstructure:"new_relic"`
	Nodes              *NodesConfig          `json:"nodes" mapstructure:"nodes"`
	Notifications      *NotificationsConfig  `json:"notifications" mapstructure:"notifications"`
	Paymail            *PaymailConfig        `json:"paymail" mapstructure:"paymail"`
	RequestLogging     bool                  `json:"request_logging" mapstructure:"request_logging"`
	Server             *ServerConfig         `json:"server" mapstructure:"server"`
	TaskManager        *TaskManagerConfig    `json:"task_manager" mapstructure:"task_manager"`
}

// General config options keys for Viper
const (
	DebugKey              = "debug"
	DebugProfilingKey     = "debug_profiling"
	DisableITCKey         = "disable_itc"
	ImportBlockHeadersKey = "import_block_headers"
	RequestLoggingKey     = "request_logging"
)

// AuthenticationConfig is the configuration for Authentication
type AuthenticationConfig struct {
	AdminKey        string `json:"admin_key" mapstructure:"admin_key"`               // key that is used for administrative requests
	RequireSigning  bool   `json:"require_signing" mapstructure:"require_signing"`   // if the signing is required
	Scheme          string `json:"scheme" mapstructure:"scheme"`                     // authentication scheme to use (default is: xpub)
	SigningDisabled bool   `json:"signing_disabled" mapstructure:"signing_disabled"` // NOTE: Only for development (turns off signing)
}

// Authentication config option keys for Viper
const (
	AuthAdminKey           = "auth.admin_key"
	AuthRequireSigningKey  = "auth.require_signing"
	AuthSchemeKey          = "auth.scheme"
	AuthSigningDisabledKey = "auth.signing_disabled"
)

// CachestoreConfig is a configuration for cachestore
type CacheConfig struct {
	Engine  cachestore.Engine `json:"engine" mapstructure:"engine"` // Cache engine to use (redis, freecache)
	Cluster *ClusterConfig    `json:"cluster" mapstructure:"cluster"`
	Redis   *RedisConfig      `json:"redis" mapstructure:"redis"`
}

// ClusterConfig is a configuration for the Bux cluster
type ClusterConfig struct {
	Coordinator cluster.Coordinator `json:"coordinator" mapstructure:"coordinator"` // redis or memory (default)
	Prefix      string              `json:"prefix" mapstructure:"prefix"`           // prefix string to use for all cluster keys, "bux" by default
	Redis       *RedisConfig        `json:"redis" mapstrcuture:"redis"`             // will use cache config if redis is set and this is empty
}

// RedisConfig is a configuration for Redis cachestore or taskmanager
type RedisConfig struct {
	DependencyMode        bool          `json:"dependency_mode" mapstructure:"dependency_mode"`                 // Only in Redis with script enabled
	MaxActiveConnections  int           `json:"max_active_connections" mapstructure:"max_active_connections"`   // Max active connections
	MaxConnectionLifetime time.Duration `json:"max_connection_lifetime" mapstructure:"max_connection_lifetime"` // Max connection lifetime
	MaxIdleConnections    int           `json:"max_idle_connections" mapstructure:"max_idle_connections"`       // Max idle connections
	MaxIdleTimeout        time.Duration `json:"max_idle_timeout" mapstructure:"max_idle_timeout"`               // Max idle timeout
	URL                   string        `json:"url" mapstructure:"url"`                                         // Redis URL connection string
	UseTLS                bool          `json:"use_tls" mapstructure:"use_tls"`                                 // Flag for using TLS
}

// Cache config keys for Viper
const (
	CacheEngineKey                = "cache.engine"
	ClusterCoordinatorKey         = "cache.cluster.coordinator"
	ClusterPrefixKey              = "cache.cluster.prefix"
	ClusterRedisURLKey            = "cache.cluster.redis.url"
	ClusterRedisMaxIdleTimeoutKey = "cache.cluster.redis.max_idle_timeout"
	ClusterRedisUseTLSKey         = "cache.cluster.redis.use_tls"
	RedisDependencyModeKey        = "cache.redis.dependency_mode"
	RedisMaxActiveConnectionsKey  = "cache.redis.max_active_connections"
	RedisMaxConnectionLifetimeKey = "cache.redis.max_connection_lifetime"
	RedisMaxIdleConnectionsKey    = "cache.redis.max_idle_connections"
	RedisMaxIdleTimeoutKey        = "cache.redis.max_idle_timeout"
	RedisURLKey                   = "cache.redis.url"
	RedisUseTLSKey                = "cache.redis.use_tls"
)

// DbConfig consists of datastore config and specific dbs configs
type DbConfig struct {
	Datastore *DatastoreConfig         `json:"datastore" mapstructure:"datastore"`
	Mongo     *datastore.MongoDBConfig `json:"mongodb" mapstructure:"mongodb"`
	SQL       *datastore.SQLConfig     `json:"sql" mapstructure:"sql"`
	SQLite    *datastore.SQLiteConfig  `json:"sqlite" mapstructure:"sqlite"`
}

// DatastoreConfig is a configuration for the datastore
type DatastoreConfig struct {
	Debug       bool             `json:"debug" mapstructure:"debug"`               // true for sql statements
	Engine      datastore.Engine `json:"engine" mapstructure:"engine"`             // mysql, sqlite
	TablePrefix string           `json:"table_prefix" mapstructure:"table_prefix"` // pre_users (pre)
}

// Common datastore config keys
const (
	DatastoreDebugKey       = "db.datastore.debug"
	DatastoreEngineKey      = "db.datastore.engine"
	DatastoreTablePrefixKey = "db.datastore.table_prefix"
)

// MongoDB config keys
const (
	MongoDatabaseNameKey = "db.mongodb.db_name"
	MongoTransactionsKey = "db.mongodb.transactions"
	MongoURIKey          = "db.mongodb.uri"
)

// SQL (MySQL, PostgreSQL) config keys
const (
	SQLDriverKey                    = "db.sql.driver"
	SQLHostKey                      = "db.sql.host"
	SQLNameKey                      = "db.sql.name"
	SQLPasswordKey                  = "db.sql.password"
	SQLPortKey                      = "db.sql.port"
	SQLReplicaKey                   = "db.sql.replica"
	SQLSkipInitializeWithVersionKey = "db.sql.skip_initialize_with_version"
	SQLTimeZoneKey                  = "db.sql.time_zone"
	SQLTxTimeoutKey                 = "db.sql.tx_timeout"
	SQLUserKey                      = "db.sql.user"
)

// SQLite config keys
const (
	SQLiteDatabasePathKey = "db.sqlite.database_path"
	SQLiteSharedKey       = "db.sqlite.shared"
)

// GraphqlConfig is the configuration for the GraphQL server
type GraphqlConfig struct {
	Enabled        bool   `json:"enabled" mapstructure:"enabled"`                 // true/false
}

// GraphQL config keys for Viper
const (
	GraphqlEnabledKey        = "graphql.enabled"
)

// MonitorOptions is the configuration for blockchain monitoring
type MonitorOptions struct {
	AuthToken                   string  `json:"auth_token" mapstructure:"auth_token"`                                       // Token to connect to the server with
	BuxAgentURL                 string  `json:"bux_agent_url" mapstructure:"bux_agent_url"`                                 // The BUX agent server url address
	Debug                       bool    `json:"debug" mapstructure:"debug"`                                                 // true/false
	Enabled                     bool    `json:"enabled" mapstructure:"enabled"`                                             // true/false
	FalsePositiveRate           float64 `json:"false_positive_rate" mapstructure:"false_positive_rate"`                     // how many false positives do we except (default: 0.01)
	LoadMonitoredDestinations   bool    `json:"load_monitored_destinations" mapstructure:"load_monitored_destinations"`     // Whether to load monitored destinations`
	MaxNumberOfDestinations     int     `json:"max_number_of_destinations" mapstructure:"max_number_of_destinations"`       // how many destinations can the filter hold (default: 100,000)
	MonitorDays                 int     `json:"monitor_days" mapstructure:"monitor_days"`                                   // how many days in the past should we monitor an address (default: 7)
	ProcessorType               string  `json:"processor_type" mapstructure:"processor_type"`                               // Type of processor to start monitor with. Default: bloom
	SaveTransactionDestinations bool    `json:"save_transaction_destinations" mapstructure:"save_transaction_destinations"` // Whether to save destinations on monitored transactions
}

// Monitor config keys for Viper
const (
	MonitorAuthTokenKey                   = "monitor.auth_token" // #nosec G101
	MonitorBuxAgentURLKey                 = "monitor.bux_agent_url"
	MonitorDebugKey                       = "monitor.debug"
	MonitorEnabledKey                     = "monitor.enabled"
	MonitorFalsePositiveRateKey           = "monitor.false_positive_rate"
	MonitorLoadMonitoredDestinationsKey   = "monitor.load_monitored_destinations"
	MonitorMaxNumberOfDestinationsKey     = "monitor.max_number_of_destinations"
	MonitorMonitorDaysKey                 = "monitor.monitor_days"
	MonitorProcessorTypeKey               = "monitor.processor_type"
	MonitorSaveTransactionDestinationsKey = "monitor.save_transaction_destinations"
)

// NewRelicConfig is the configuration for New Relic
type NewRelicConfig struct {
	DomainName string `json:"domain_name" mapstructure:"domain_name"` // used for hostname display
	Enabled    bool   `json:"enabled" mapstructure:"enabled"`         // true/false
	LicenseKey string `json:"license_key" mapstructure:"license_key"` // 2342-3423523-62
}

// NewRelic config keys for Viper
const (
	NewRelicDomainNameKey = "new_relic.domain_name"
	NewRelicEnabledKey    = "new_relic.enabled"
	NewRelicLicenseKeyKey = "new_relic.license_key"
)

// NodesConfig consists of blockchain nodes (such as Minercraft and Arc) configuration
type NodesConfig struct {
	UseMapiFeeQuotes     bool                    `json:"use_mapi_fee_quotes" mapstructure:"use_mapi_fee_quotes"`
	MinercraftAPI        string                  `json:"minercraft_api" mapstructure:"minercraft_api"`
	MinercraftCustomAPIs []*minercraft.MinerAPIs `json:"minercraft_custom_apis" mapstructure:"minercraft_custom_apis"`
	BroadcastClientAPIs  []string                `json:"broadcast_client_apis" mapstructure:"broadcast_client_apis"`
}

// Nodes config keys for viper
const (
	NodesUseMapiFeeQuotesKey    = "nodes.use_mapi_fee_quotes"
	NodesMinercraftAPIKey       = "nodes.minercraft_api"
	NodesBroadcastClientAPIsKey = "nodes.broadcast_client_apis"
)

// NotificationsConfig is the configuration for notifications
type NotificationsConfig struct {
	Enabled         bool   `json:"enabled" mapstructure:"enabled"`
	WebhookEndpoint string `json:"webhook_endpoint" mapstructure:"webhook_endpoint"`
}

// Notification config keys for Viper
const (
	NotificationsEnabledKey         = "notifications.enabled"
	NotificationsWebhookEndpointKey = "notifications.webhook_endpoint"
)

// LoggingConfig is a configuration for logging
type LoggingConfig struct {
	Level        string `json:"level" mapstructure:"level"`
	Format       string `json:"format" mapstructure:"format"`
	InstanceName string `json:"instance_name" mapstructure:"instance_name"`
	LogOrigin    bool   `json:"log_origin" mapstructure:"log_origin"`
}

// PaymailConfig is the configuration for the built-in Paymail server
type PaymailConfig struct {
	Beef                    *BeefConfig `json:"beef" mapstructure:"beef"`                                           // Background Evaluation Extended Format (BEEF)
	DefaultFromPaymail      string      `json:"default_from_paymail" mapstructure:"default_from_paymail"`           // IE: from@domain.com
	DefaultNote             string      `json:"default_note" mapstructure:"default_note"`                           // IE: message needed for address resolution
	Domains                 []string    `json:"domains" mapstructure:"domains"`                                     // List of allowed domains
	DomainValidationEnabled bool        `json:"domain_validation_enabled" mapstructure:"domain_validation_enabled"` // Turn off if hosted domain is not paymail related
	Enabled                 bool        `json:"enabled" mapstructure:"enabled"`                                     // Flag for enabling the Paymail Server Service
	SenderValidationEnabled bool        `json:"sender_validation_enabled" mapstructure:"sender_validation_enabled"` // Turn on extra security
}

// BeefConfig consists of components required to use beef, e.g. Pulse for merkle roots validation
type BeefConfig struct {
	UseBeef                  bool   `json:"use_beef" mapstructure:"use_beef"`
	PulseHeaderValidationURL string `json:"pulse_url" mapstructure:"pulse_url"`
	PulseAuthToken           string `json:"pulse_auth_token" mapstructure:"pulse_auth_token"`
}

// Paymail config keys for Viper
const (
	UseBeefKey                        = "paymail.beef.use_beef"
	PulseHeaderValidationURLKey       = "paymail.beef.pulse_url"
	PulseAuthTokenKey                 = "paymail.beef.pulse_auth_token" // #nosec G101
	PaymailDefaultFromPaymailKey      = "paymail.default_from_paymail"
	PaymailDefaultNoteKey             = "paymail.default_note"
	PaymailDomainsKey                 = "paymail.domains"
	PaymailDomainValidationEnabledKey = "paymail.domain_validation_enabled"
	PaymailEnabledKey                 = "paymail.enabled"
	PaymailSenderValidationEnabledKey = "paymail.sender_validation_enabled"
)

// TaskManagerConfig is a configuration for the taskmanager
type TaskManagerConfig struct {
	// QueueOptions *taskq.QueueOptions
	Factory   taskmanager.Factory `json:"factory" mapstructure:"factory"`       // Factory (memory, redis)
}

// TaskManager config keys for Viper
const (
	TaskManagerFactoryKey   = "task_manager.factory"
)

// ServerConfig is a configuration for the HTTP Server
type ServerConfig struct {
	IdleTimeout  time.Duration `json:"idle_timeout" mapstructure:"idle_timeout"`   // 60s
	ReadTimeout  time.Duration `json:"read_timeout" mapstructure:"read_timeout"`   // 15s
	WriteTimeout time.Duration `json:"write_timeout" mapstructure:"write_timeout"` // 15s
	Port         string        `json:"port" mapstructure:"port"`                   // 3003
}

// Server config keys for Viper
const (
	ServerIdleTimeoutKey  = "server.idle_timeout"
	ServerReadTimeoutKey  = "server.read_timeout"
	ServerWriteTimeoutKey = "server.write_timeout"
	ServerPortKey         = "server.port"
)

// GetUserAgent will return the outgoing user agent
func (a *AppConfig) GetUserAgent() string {
	return "BUX-Server " + Version
}

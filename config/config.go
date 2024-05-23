// Package config provides a configuration for the API
package config

import (
	"time"

	"github.com/bitcoin-sv/spv-wallet/engine/cluster"
	"github.com/bitcoin-sv/spv-wallet/engine/datastore"
	"github.com/bitcoin-sv/spv-wallet/engine/taskmanager"
	"github.com/mrz1836/go-cachestore"
)

// Config constants used for spv-wallet
const (
	ApplicationName         = "SPVWallet"
	APIVersion              = "v1"
	DefaultNewRelicShutdown = 10 * time.Second
	HealthRequestPath       = "health"
	Version                 = "v0.12.0"
	ConfigFilePathKey       = "config_file"
	DefaultConfigFilePath   = "config.yaml"
	EnvPrefix               = "SPVWALLET"
	BroadcastCallbackRoute  = "/transaction/broadcast/callback"
)

// AppConfig is the configuration values and associated env vars
type AppConfig struct {
	// NewRelic is New Relic related settings.
	NewRelic *NewRelicConfig `json:"new_relic" mapstructure:"new_relic"`
	// TaskManager is a configuration for Task Manager in SPV Wallet.
	TaskManager *TaskManagerConfig `json:"task_manager" mapstructure:"task_manager"`
	// Authentication is the configuration for keys authentication in SPV Wallet.
	Authentication *AuthenticationConfig `json:"auth" mapstructure:"auth"`
	// Server is a general configuration for spv-wallet.
	Server *ServerConfig `json:"server_config" mapstructure:"server_config"`
	// Nodes is a config for BSV nodes, mAPI and Arc.
	Nodes *NodesConfig `json:"nodes" mapstructure:"nodes"`
	// Metrics is a configuration for metrics in SPV Wallet.
	Metrics *MetricsConfig `json:"metrics" mapstructure:"metrics"`
	// ExperimentalFeatures is a configuration that allows to enable features that are considered experimental/non-production.
	ExperimentalFeatures *ExperimentalConfig `json:"experimental_features" mapstructure:"experimental_features"`
	// Notifications is a config for Notification service.
	Notifications *NotificationsConfig `json:"notifications" mapstructure:"notifications"`
	// Db is the configuration for database related settings.
	Db *DbConfig `json:"db" mapstructure:"db"`
	// Cache is the configuration for cache, memory or redis, and cluster cache settings.
	Cache *CacheConfig `json:"cache" mapstructure:"cache"`
	// Logging is the configuration for zerolog used in SPV Wallet.
	Logging *LoggingConfig `json:"logging" mapstructure:"logging"`
	// Paymail is a config for Paymail and BEEF.
	Paymail *PaymailConfig `json:"paymail" mapstructure:"paymail"`
	// ImportBlockHeaders is a URL from where the headers can be downloaded.
	ImportBlockHeaders string `json:"import_block_headers" mapstructure:"import_block_headers"`
	// Debug is a flag for enabling additional information from SPV Wallet.
	Debug bool `json:"debug" mapstructure:"debug"`
	// DebugProfiling is a flag for enabling additinal debug profiling.
	DebugProfiling bool `json:"debug_profiling" mapstructure:"debug_profiling"`
	// DisableITC is a flag for disabling Incoming Transaction Checking.
	DisableITC bool `json:"disable_itc" mapstructure:"disable_itc"`
	// RequestLogging is flag for enabling logging in go-api-router.
	RequestLogging bool `json:"request_logging" mapstructure:"request_logging"`
}

// AuthenticationConfig is the configuration for Authentication
type AuthenticationConfig struct {
	// AdminKey is used for administrative requests
	AdminKey string `json:"admin_key" mapstructure:"admin_key"`
	// Scheme it the authentication scheme to use (default is: xpub)
	Scheme string `json:"scheme" mapstructure:"scheme"`
	// RequireSigning is the flag that decides if the signing is required
	RequireSigning bool `json:"require_signing" mapstructure:"require_signing"`
	// SigningDisabled turns off signing. NOTE: Only for development
	SigningDisabled bool `json:"signing_disabled" mapstructure:"signing_disabled"`
}

// CacheConfig is a configuration for cachestore
type CacheConfig struct {
	// Cluster is the cluster-specific configuration for SPV Wallet.
	Cluster *ClusterConfig `json:"cluster" mapstructure:"cluster"`
	// Redis is a general config for redis if the engine is set to it.
	Redis *RedisConfig `json:"redis" mapstructure:"redis"`
	// Engine is the cache engine to use (redis, freecache).
	Engine cachestore.Engine `json:"engine" mapstructure:"engine"`
}

// CallbackConfig is the configuration for callbacks
type CallbackConfig struct {
	// CallbackHost is the URL for broadcast callback registration.
	CallbackHost string `json:"callback_host" mapstructure:"callback_host"`
	// CallbackToken is the token for broadcast callback registration.
	CallbackToken string `json:"callback_token" mapstructure:"callback_token"`
}

// ClusterConfig is a configuration for the SPV Wallet cluster
type ClusterConfig struct {
	// Redis is cluster-specific redis config, will use cache config if this is unset.
	Redis *RedisConfig `json:"redis" mapstrcuture:"redis"`
	// Coordinator is a cluster coordinator (redis or memory).
	Coordinator cluster.Coordinator `json:"coordinator" mapstructure:"coordinator"`
	// Prefix is the string to use for all cluster keys.
	Prefix string `json:"prefix" mapstructure:"prefix"`
}

// RedisConfig is a configuration for Redis cachestore or taskmanager
type RedisConfig struct {
	// URL is Redis url connection string.
	URL string `json:"url" mapstructure:"url"`
	// MaxActiveConnections is maximum number of active redis connections.
	MaxActiveConnections int `json:"max_active_connections" mapstructure:"max_active_connections"`
	// MaxIdleConnections is the maximum number of idle connections.
	MaxIdleConnections int `json:"max_idle_connections" mapstructure:"max_idle_connections"`
	// MaxConnectionLifetime is the maximum duration of the connection.
	MaxConnectionLifetime time.Duration `json:"max_connection_lifetime" mapstructure:"max_connection_lifetime"`
	// MaxIdleTimeout is the maximum duration of idle redis connection before timeout.
	MaxIdleTimeout time.Duration `json:"max_idle_timeout" mapstructure:"max_idle_timeout"`
	// DependencyMode works only in Redis with script enabled.
	DependencyMode bool `json:"dependency_mode" mapstructure:"dependency_mode"`
	// UseTLS is a flag which decides whether to use TLS
	UseTLS bool `json:"use_tls" mapstructure:"use_tls"`
}

// DbConfig consists of datastore config and specific dbs configs
type DbConfig struct {
	// Datastore general config.
	Datastore *DatastoreConfig `json:"datastore" mapstructure:"datastore"`
	// Mongo is a config for MongoDb. Works only if datastore engine is set to mongodb.
	Mongo *datastore.MongoDBConfig `json:"mongodb" mapstructure:"mongodb"`
	// SQL is a config for PostgreSQL. Works only if datastore engine is set to postgresql.
	SQL *datastore.SQLConfig `json:"sql" mapstructure:"sql"`
	// SQLite is a config for SQLite. Works only if datastore engine is set to sqlite.
	SQLite *datastore.SQLiteConfig `json:"sqlite" mapstructure:"sqlite"`
}

// DatastoreConfig is a configuration for the datastore
type DatastoreConfig struct {
	// TablePrefix is the prefix for all table names in the database.
	TablePrefix string `json:"table_prefix" mapstructure:"table_prefix"`
	// Engine is the database to be used, sqlite, postgresql.
	Engine datastore.Engine `json:"engine" mapstructure:"engine"`
	// Debug is a flag that decides whether additional output (such as sql statements) should be produced from datastore.
	Debug bool `json:"debug" mapstructure:"debug"`
}

// NewRelicConfig is the configuration for New Relic
type NewRelicConfig struct {
	// DomainName is used for hostname display.
	DomainName string `json:"domain_name" mapstructure:"domain_name"`
	// LicenseKey is the New Relic license key.
	LicenseKey string `json:"license_key" mapstructure:"license_key"`
	// Enabled is the flag that enables New Relic service.
	Enabled bool `json:"enabled" mapstructure:"enabled"`
}

// NodesConfig consists of blockchain nodes (such as Minercraft and Arc) configuration
type NodesConfig struct {
	Callback     *CallbackConfig `json:"callback" mapstructure:"callback"`
	FeeUnit      *FeeUnitConfig  `json:"fee_unit" mapstructure:"fee_unit"`
	DeploymentID string          `json:"deployment_id" mapstructure:"deployment_id"`
	Protocol     NodesProtocol   `json:"protocol" mapstructure:"protocol"`
	Apis         []*MinerAPI     `json:"apis" mapstructure:"apis"`
	UseFeeQuotes bool            `json:"use_fee_quotes" mapstructure:"use_fee_quotes"`
}

// FeeUnitConfig reflects the utils.FeeUnit struct with proper annotations for json and mapstructure
type FeeUnitConfig struct {
	Satoshis int `json:"satoshis" mapstructure:"satoshis"`
	Bytes    int `json:"bytes" mapstructure:"bytes"`
}

// MinerAPI holds connection info for a single miner endpoint
type MinerAPI struct {
	Token   string `json:"token" mapstructure:"token"`
	ArcURL  string `json:"arc_url" mapstructure:"arc_url"`
	MapiURL string `json:"mapi_url" mapstructure:"mapi_url"`

	// MinerID is not used with ARC potocol
	MinerID string `json:"minerid" mapstructure:"minerid"`
}

// NotificationsConfig is the configuration for notifications
type NotificationsConfig struct {
	// WebhookEndpoint is the endpoint for webhook registration.
	WebhookEndpoint string `json:"webhook_endpoint" mapstructure:"webhook_endpoint"`
	// Enabled is the flag that enables notifications service.
	Enabled bool `json:"enabled" mapstructure:"enabled"`
}

// LoggingConfig is a configuration for logging
type LoggingConfig struct {
	// Level is the importance and amount of information printed: debug, info, warn, error, fatal, panic, etc.
	Level string `json:"level" mapstructure:"level"`
	// Format is the format of logs, json (for gathering eg. into elastic) or console (for stdout).
	Format string `json:"format" mapstructure:"format"`
	// InstanceName is the name of the zerolog instance.
	InstanceName string `json:"instance_name" mapstructure:"instance_name"`
	// LogOrigin is the flag for whether the origin of logs should be printed.
	LogOrigin bool `json:"log_origin" mapstructure:"log_origin"`
}

// PaymailConfig is the configuration for the built-in Paymail server
type PaymailConfig struct {
	// Beef is for Background Evaluation Extended Format (BEEF) config.
	Beef *BeefConfig `json:"beef" mapstructure:"beef"`
	// DefaultFromPaymail IE: from@domain.com.
	DefaultFromPaymail string `json:"default_from_paymail" mapstructure:"default_from_paymail"`
	// Domains is a list of allowed domains.
	Domains []string `json:"domains" mapstructure:"domains"`
	// DomainValidationEnabled should be turned off if hosted domain is not paymail related.
	DomainValidationEnabled bool `json:"domain_validation_enabled" mapstructure:"domain_validation_enabled"`
	// SenderValidationEnabled should be turned on for extra security.
	SenderValidationEnabled bool `json:"sender_validation_enabled" mapstructure:"sender_validation_enabled"`
}

// BeefConfig consists of components required to use beef, e.g. Block Headers Service for merkle roots validation
type BeefConfig struct {
	// BlockHeaderServiceHeaderValidationURL is the URL for merkle roots validation in Block Headers Service.
	BlockHeaderServiceHeaderValidationURL string `json:"block_header_service_url" mapstructure:"block_header_service_url"`
	// BlockHeaderServiceAuthToken is the authentication token for validating merkle roots in Block Headers Service.
	BlockHeaderServiceAuthToken string `json:"block_header_service_auth_token" mapstructure:"block_header_service_auth_token"`
	// UseBeef is a flag for enabling BEEF transactions format.
	UseBeef bool `json:"use_beef" mapstructure:"use_beef"`
}

func (b *BeefConfig) enabled() bool {
	return b != nil && b.UseBeef
}

// TaskManagerConfig is a configuration for the taskmanager
type TaskManagerConfig struct {
	// Factory is the Task Manager factory, memory or redis.
	Factory taskmanager.Factory `json:"factory" mapstructure:"factory"`
}

// ServerConfig is a configuration for the HTTP Server
type ServerConfig struct {
	// IdleTimeout is the maximum duration before server timeout.
	IdleTimeout time.Duration `json:"idle_timeout" mapstructure:"idle_timeout"`
	// ReadTimeout is the maximum duration for server read timeout.
	ReadTimeout time.Duration `json:"read_timeout" mapstructure:"read_timeout"`
	// WriteTimeout is the maximum duration for server write timeout.
	WriteTimeout time.Duration `json:"write_timeout" mapstructure:"write_timeout"`
	// Port is the port that the server should use.
	Port int `json:"port" mapstructure:"port"`
}

// MetricsConfig represents a metrics config.
type MetricsConfig struct {
	// Enabled is a flag for enabling metrics.
	Enabled bool `json:"enabled" mapstructure:"enabled"`
}

// ExperimentalConfig represents a feature flag config.
type ExperimentalConfig struct {
	PikeEnabled bool `json:"pike_enabled" mapstructure:"pike_enabled"`
}

// GetUserAgent will return the outgoing user agent
func (a *AppConfig) GetUserAgent() string {
	return "SPV Wallet " + Version
}

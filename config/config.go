// Package config provides a configuration for the API
package config

import (
	"time"

	"github.com/BuxOrg/bux/cluster"
	"github.com/BuxOrg/bux/taskmanager"
	"github.com/mrz1836/go-cachestore"
	"github.com/mrz1836/go-datastore"
)

// Config constants used for bux-server
const (
	ApplicationName         = "BuxServer"
	APIVersion              = "v1"
	DefaultNewRelicShutdown = 10 * time.Second
	HealthRequestPath       = "health"
	Version                 = "v0.5.16"
	ConfigFilePathKey       = "config_file"
	DefaultConfigFilePath   = "config.yaml"
	ConfigEnvPrefix         = "BUX_"
)

// AppConfig is the configuration values and associated env vars
type AppConfig struct {
	// Authentication is the configuration for keys authentication in bux.
	Authentication *AuthenticationConfig `json:"auth" mapstructure:"auth"`
	// Cache is the configuration for cache, memory or redis, and cluster cache settings.
	Cache *CacheConfig `json:"cache" mapstructure:"cache"`
	// Db is the configuration for database related settings.
	Db *DbConfig `json:"db" mapstructure:"db"`
	// Debug is a flag for enabling additional information from bux.
	Debug bool `json:"debug" mapstructure:"debug"`
	// DebugProfiling is a flag for enabling additinal debug profiling.
	DebugProfiling bool `json:"debug_profiling" mapstructure:"debug_profiling"`
	// DisableITC is a flag for disabling Incoming Transaction Checking.
	DisableITC bool `json:"disable_itc" mapstructure:"disable_itc"`
	// GraphQL is GraphQL related settings.
	GraphQL *GraphqlConfig `json:"graphql" mapstructure:"graphql"`
	// ImportBlockHeaders is a URL from where the headers can be downloaded.
	ImportBlockHeaders string `json:"import_block_headers" mapstructure:"import_block_headers"`
	// Logging is the configuration for zerolog used in bux.
	Logging *LoggingConfig `json:"logging" mapstructure:"logging"`
	// NewRelic is New Relic related settings.
	NewRelic *NewRelicConfig `json:"new_relic" mapstructure:"new_relic"`
	// Nodes is a config for BSV nodes, mAPI and Arc.
	Nodes *NodesConfig `json:"nodes" mapstructure:"nodes"`
	// Notifications is a config for Notification service.
	Notifications *NotificationsConfig `json:"notifications" mapstructure:"notifications"`
	// Paymail is a config for Paymail and BEEF.
	Paymail *PaymailConfig `json:"paymail" mapstructure:"paymail"`
	// RequestLogging is flag for enabling logging in go-api-router.
	RequestLogging bool `json:"request_logging" mapstructure:"request_logging"`
	// Server is a general configuration for bux-server.
	Server *ServerConfig `json:"server_config" mapstructure:"server_config"`
	// TaskManager is a configuration for Task Manager in bux.
	TaskManager *TaskManagerConfig `json:"task_manager" mapstructure:"task_manager"`
}

// AuthenticationConfig is the configuration for Authentication
type AuthenticationConfig struct {
	// AdminKey is used for administrative requests
	AdminKey string `json:"admin_key" mapstructure:"admin_key"`
	// RequireSigning is the flag that decides if the signing is required
	RequireSigning bool `json:"require_signing" mapstructure:"require_signing"`
	// Scheme it the authentication scheme to use (default is: xpub)
	Scheme string `json:"scheme" mapstructure:"scheme"`
	// SigningDisabled turns off signing. NOTE: Only for development
	SigningDisabled bool `json:"signing_disabled" mapstructure:"signing_disabled"`
}

// CacheConfig is a configuration for cachestore
type CacheConfig struct {
	// Engine is the cache engine to use (redis, freecache).
	Engine cachestore.Engine `json:"engine" mapstructure:"engine"`
	// Cluster is the cluster-specific configuration for bux.
	Cluster *ClusterConfig `json:"cluster" mapstructure:"cluster"`
	// Redis is a general config for redis if the engine is set to it.
	Redis *RedisConfig `json:"redis" mapstructure:"redis"`
}

// ClusterConfig is a configuration for the Bux cluster
type ClusterConfig struct {
	// Coordinator is a cluster coordinator (redis or memory).
	Coordinator cluster.Coordinator `json:"coordinator" mapstructure:"coordinator"`
	// Prefix is the string to use for all cluster keys.
	Prefix string `json:"prefix" mapstructure:"prefix"`
	// Redis is cluster-specific redis config, will use cache config if this is unset.
	Redis *RedisConfig `json:"redis" mapstrcuture:"redis"`
}

// RedisConfig is a configuration for Redis cachestore or taskmanager
type RedisConfig struct {
	// DependencyMode works only in Redis with script enabled.
	DependencyMode bool `json:"dependency_mode" mapstructure:"dependency_mode"`
	// MaxActiveConnections is maximum number of active redis connections.
	MaxActiveConnections int `json:"max_active_connections" mapstructure:"max_active_connections"`
	// MaxConnectionLifetime is the maximum duration of the connection.
	MaxConnectionLifetime time.Duration `json:"max_connection_lifetime" mapstructure:"max_connection_lifetime"`
	// MaxIdleConnections is the maximum number of idle connections.
	MaxIdleConnections int `json:"max_idle_connections" mapstructure:"max_idle_connections"`
	// MaxIdleTimeout is the maximum duration of idle redis connection before timeout.
	MaxIdleTimeout time.Duration `json:"max_idle_timeout" mapstructure:"max_idle_timeout"`
	// URL is Redis url connection string.
	URL string `json:"url" mapstructure:"url"`
	// UseTLS is a flag which decides whether to use TLS
	UseTLS bool `json:"use_tls" mapstructure:"use_tls"`
}

// DbConfig consists of datastore config and specific dbs configs
type DbConfig struct {
	// Datastore is a general go-datastore config.
	Datastore *DatastoreConfig `json:"datastore" mapstructure:"datastore"`
	// Mongo is a config for MongoDb. Works only if datastore engine is set to mongodb.
	Mongo *datastore.MongoDBConfig `json:"mongodb" mapstructure:"mongodb"`
	// SQL is a config for PostgreSQL or MySQL. Works only if datastore engine is set to postgresql or mysql.
	SQL *datastore.SQLConfig `json:"sql" mapstructure:"sql"`
	// SQLite is a config for SQLite. Works only if datastore engine is set to sqlite.
	SQLite *datastore.SQLiteConfig `json:"sqlite" mapstructure:"sqlite"`
}

// DatastoreConfig is a configuration for the datastore
type DatastoreConfig struct {
	// Debug is a flag that decides whether additional output (such as sql statements) should be produced from datastore.
	Debug bool `json:"debug" mapstructure:"debug"`
	// Engine is the database to be used, mysql, sqlite, postgresql.
	Engine datastore.Engine `json:"engine" mapstructure:"engine"`
	// TablePrefix is the prefix for all table names in the database.
	TablePrefix string `json:"table_prefix" mapstructure:"table_prefix"`
}

// GraphqlConfig is the configuration for the GraphQL server
type GraphqlConfig struct {
	// Enabled is a flag that says whether graphql should be enabled.
	Enabled bool `json:"enabled" mapstructure:"enabled"`
}

// NewRelicConfig is the configuration for New Relic
type NewRelicConfig struct {
	// DomainName is used for hostname display.
	DomainName string `json:"domain_name" mapstructure:"domain_name"`
	// Enabled is the flag that enables New Relic service.
	Enabled bool `json:"enabled" mapstructure:"enabled"`
	// LicenseKey is the New Relic license key.
	LicenseKey string `json:"license_key" mapstructure:"license_key"`
}

// NodesConfig consists of blockchain nodes (such as Minercraft and Arc) configuration
type NodesConfig struct {
	Protocol NodesProtocol `json:"protocol" mapstructure:"protocol"`
	Apis     []*MinerAPI   `json:"apis" mapstructure:"apis"`
	Mapi     *MapiConfig   `json:"mapi" mapstructure:"mapi"`
}

// MinerAPI holds connection info for a single miner endpoint
type MinerAPI struct {
	Token   string `json:"token" mapstructure:"token"`
	ArcURL  string `json:"arc_url" mapstructure:"arc_url"`
	MapiURL string `json:"mapi_url" mapstructure:"mapi_url"`

	// MinerID is not used with ARC potocol
	MinerID string `json:"minerid" mapstructure:"minerid"`
}

// MapiConfig holds mApi-specific configuration
type MapiConfig struct {
	UseFeeQuotes bool `json:"use_fee_quotes" mapstructure:"use_fee_quotes"`
}

// NotificationsConfig is the configuration for notifications
type NotificationsConfig struct {
	// Enabled is the flag that enables notifications service.
	Enabled bool `json:"enabled" mapstructure:"enabled"`
	// WebhookEndpoint is the endpoint for webhook registration.
	WebhookEndpoint string `json:"webhook_endpoint" mapstructure:"webhook_endpoint"`
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
	// DefaultNote IE: message needed for address resolution.
	DefaultNote string `json:"default_note" mapstructure:"default_note"`
	// Domains is a list of allowed domains.
	Domains []string `json:"domains" mapstructure:"domains"`
	// DomainValidationEnabled should be turned off if hosted domain is not paymail related.
	DomainValidationEnabled bool `json:"domain_validation_enabled" mapstructure:"domain_validation_enabled"`
	// Enabled is a flag for enabling the Paymail Server Service.
	Enabled bool `json:"enabled" mapstructure:"enabled"`
	// SenderValidationEnabled should be turned on for extra security.
	SenderValidationEnabled bool `json:"sender_validation_enabled" mapstructure:"sender_validation_enabled"`
}

// BeefConfig consists of components required to use beef, e.g. Pulse for merkle roots validation
type BeefConfig struct {
	// UseBeef is a flag for enabling BEEF transactions format.
	UseBeef bool `json:"use_beef" mapstructure:"use_beef"`
	// PulseHeaderValidationURL is the URL for headers validation in Pulse.
	PulseHeaderValidationURL string `json:"pulse_url" mapstructure:"pulse_url"`
	// PulseAuthToken is the authentication token for validating headers in Pulse.
	PulseAuthToken string `json:"pulse_auth_token" mapstructure:"pulse_auth_token"`
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

// GetUserAgent will return the outgoing user agent
func (a *AppConfig) GetUserAgent() string {
	return "BUX-Server " + Version
}

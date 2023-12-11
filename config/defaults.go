package config

import (
	"github.com/spf13/viper"
)

// DefaultConfigFilePath is in the root folder by default
const DefaultConfigFilePath = "config.yml"

// General defaults
const (
	DebugDefault              = true
	DebugProfilingDefault     = true
	DisableITCDefault         = true
	ImportBlockHeadersDefault = ""
	RequestLoggingDefault     = true
)

// Authentication defaults
const (
	AuthAdminKeyDefault        = "xpub661MyMwAqRbcFrBJbKwBGCB7d3fr2SaAuXGM95BA62X41m6eW2ehRQGW4xLi9wkEXUGnQZYxVVj4PxXnyrLk7jdqvBAs1Qq9gf6ykMvjR7J"
	AuthRequireSigningDefault  = false
	AuthSchemeDefault          = "xpub"
	AuthSigningDisabledDefault = true
)

// Cache defaults
const (
	CacheEngineDefault                = "freecache"
	ClusterCoordinatorDefault         = "memory"
	ClusterPrefixDefault              = "bux_cluser_"
	ClusterRedisURLDefault            = "localhost:6379"
	ClusterRedisMaxIdleTimeoutDefault = "10s"
	ClusterRedisUseTLSDefault         = false
	RedisDependencyModeDefault        = true
	RedisMaxActiveConnectionsDefault  = 0
	RedisMaxConnectionLifetimeDefault = "60s"
	RedisMaxIdleConnectionsDefault    = 10
	RedisMaxIdleTimeoutDefault        = "10s"
	RedisURLDefault                   = "redis://localhost:6379"
	RedisUseTLSDefault                = false
)

// Datastore defaults
const (
	DatastoreAutoMigrateDefault = true
	DatastoreDebugDefault       = false
	DatastoreEngineDefault      = "sqlite"
	DatastoreTablePrefixDefault = "xapi"

	// Mongo defaults
	MongoDatabaseNameDefault = "xapi"
	MongoTransactionsDefault = false
	MongoURIDefault          = "mongodb://localhost:27017/xapi"

	// SQL (MySQL, PostgreSQL) config keys
	SQLDriverDefault                    = "postgresql"
	SQLHostDefault                      = "localhost"
	SQLNameDefault                      = "xapi"
	SQLPasswordDefault                  = ""
	SQLPortDefault                      = "5432"
	SQLReplicaDefault                   = false
	SQLSkipInitializeWithVersionDefault = true
	SQLTimeZoneDefault                  = "UTC"
	SQLTxTimeoutDefault                 = "10s"
	SQLUserDefault                      = "postgres"

	// SQLite config keys
	SQLiteDatabasePathDefault = "./bux.db"
	SQLiteSharedDefault       = true
)

// Graphql defaults
const (
	GraphqlEnabledDefault        = false
	GraphqlPlaygroundPathDefault = "/graphql"
	GraphqlServerPathDefault     = "/graphiql"
)

// Monitor defaults
const (
	MonitorAuthTokenDefault                   = ""
	MonitorBuxAgentURLDefault                 = "ws://localhost:8000/websocket"
	MonitorDebugDefault                       = false
	MonitorEnabledDefault                     = false
	MonitorFalsePositiveRateDefault           = 0.01
	MonitorLoadMonitoredDestinationsDefault   = false
	MonitorMaxNumberOfDestinationsDefault     = 100000
	MonitorMonitorDaysDefault                 = 7
	MonitorProcessorTypeDefault               = "bloom"
	MonitorSaveTransactionDestinationsDefault = true
)

// NewRelic defaults
const (
	NewRelicDomainNameDefault = "domain.com"
	NewRelicEnabledDefault    = false
	NewRelicLicenseKeyDefault = "BOGUS-LICENSE-KEY-1234567890987654321234"
)

// Nodes defaults
const (
	NodesUseMapiFeeQuotesDefault = true
	NodesMinercraftAPIDefault    = "mAPI"
)

// NodesBroadcastClientAPIsDefault is the default template of broadcast-client api tokens
var NodesBroadcastClientAPIsDefault = []string{}

// Notification defaults
const (
	NotificationsEnabledDefault         = false
	NotificationsWebhookEndpointDefault = ""
)

// Paymail defaults
const (
	UseBeefDefault                        = false
	PulseHeaderValidationURLDefault       = "http://localhost:8080/api/v1/chain/merkleroot/verify"
	PulseAuthTokenDefault                 = "mQZQ6WmxURxWz5ch"
	PaymailEnabledDefault                 = true
	PaymailDefaultFromPaymailDefault      = "from@domain.com"
	PaymailDefaultNoteDefault             = "bux Address Resolution"
	PaymailDomainValidationEnabledDefault = true
	PaymailSenderValidationEnabledDefault = true
)

// PaymailDomainsDefault by default allows for localhost domain
var PaymailDomainsDefault = []string{"localhost"}

// TaskManager defaults
const (
	TaskManagerEngineDefault    = "taskq"
	TaskManagerFactoryDefault   = "memory"
	TaskManagerQueueNameDefault = "development_queue"
)

// Server defaults
const (
	ServerIdleTimeoutDefault  = "60s"
	ServerReadTimeoutDefault  = "15s"
	ServerWriteTimeoutDefault = "15s"
	ServerPortDefault         = "3003"
)

func setDefaults() {
	viper.SetDefault(ConfigFilePathKey, DefaultConfigFilePath)

	setGeneralDefaults()
	setAuthDefaults()
	setBeefDefaults()
	setCachestoreDefaults()
	setClusterDefaults()
	setDbDefaults()
	setGraphqlDefaults()
	setMonitorDefaults()
	setNewRelicDefaults()
	setNodesDefaults()
	setNotificationsDefaults()
	setPaymailDefaults()
	setRedisDefaults()
	setTaskManagerDefaults()
	setServerDefaults()
}

func setGeneralDefaults() {
	viper.SetDefault(DebugKey, DebugDefault)
	viper.SetDefault(DebugProfilingKey, DebugProfilingDefault)
	viper.SetDefault(DisableITCKey, DisableITCDefault)
	viper.SetDefault(ImportBlockHeadersKey, ImportBlockHeadersDefault)
	viper.SetDefault(RequestLoggingKey, RequestLoggingDefault)
}

func setAuthDefaults() {
	viper.SetDefault(AuthAdminKey, AuthAdminKeyDefault)
	viper.SetDefault(AuthRequireSigningKey, AuthRequireSigningDefault)
	viper.SetDefault(AuthSchemeKey, AuthSchemeDefault)
	viper.SetDefault(AuthSigningDisabledKey, AuthSigningDisabledDefault)
}

func setBeefDefaults() {
	viper.SetDefault(UseBeefKey, UseBeefDefault)
	viper.SetDefault(PulseHeaderValidationURLKey, PulseHeaderValidationURLDefault)
	viper.SetDefault(PulseAuthTokenKey, PulseAuthTokenDefault)
}

func setCachestoreDefaults() {
	viper.SetDefault(CacheEngineKey, CacheEngineDefault)
}

func setClusterDefaults() {
	viper.SetDefault(ClusterCoordinatorKey, ClusterCoordinatorDefault)
	viper.SetDefault(ClusterPrefixKey, ClusterPrefixDefault)
	viper.SetDefault(ClusterRedisURLKey, ClusterRedisURLDefault)
	viper.SetDefault(ClusterRedisMaxIdleTimeoutKey, ClusterRedisMaxIdleTimeoutDefault)
	viper.SetDefault(ClusterRedisUseTLSKey, ClusterRedisUseTLSDefault)
}

func setDbDefaults() {
	viper.SetDefault(DatastoreAutoMigrateKey, DatastoreAutoMigrateDefault)
	viper.SetDefault(DatastoreDebugKey, DatastoreDebugDefault)
	viper.SetDefault(DatastoreEngineKey, DatastoreEngineDefault)
	viper.SetDefault(DatastoreTablePrefixKey, DatastoreTablePrefixDefault)

	viper.SetDefault(MongoDatabaseNameKey, MongoDatabaseNameDefault)
	viper.SetDefault(MongoTransactionsKey, MongoTransactionsDefault)
	viper.SetDefault(MongoURIKey, MongoURIDefault)

	viper.SetDefault(SQLDriverKey, SQLDriverDefault)
	viper.SetDefault(SQLHostKey, SQLHostDefault)
	viper.SetDefault(SQLNameKey, SQLNameDefault)
	viper.SetDefault(SQLPasswordKey, SQLPasswordDefault)
	viper.SetDefault(SQLPortKey, SQLPortDefault)
	viper.SetDefault(SQLReplicaKey, SQLReplicaDefault)
	viper.SetDefault(SQLSkipInitializeWithVersionKey, SQLSkipInitializeWithVersionDefault)
	viper.SetDefault(SQLTimeZoneKey, SQLTimeZoneDefault)
	viper.SetDefault(SQLTxTimeoutKey, SQLTxTimeoutDefault)
	viper.SetDefault(SQLUserKey, SQLUserDefault)

	viper.SetDefault(SQLiteDatabasePathKey, SQLiteDatabasePathDefault)
	viper.SetDefault(SQLiteSharedKey, SQLiteSharedDefault)
}

func setGraphqlDefaults() {
	viper.SetDefault(GraphqlEnabledKey, GraphqlEnabledDefault)
}

func setMonitorDefaults() {
	viper.SetDefault(MonitorAuthTokenKey, MonitorAuthTokenDefault)
	viper.SetDefault(MonitorBuxAgentURLKey, MonitorBuxAgentURLDefault)
	viper.SetDefault(MonitorDebugKey, MonitorDebugDefault)
	viper.SetDefault(MonitorEnabledKey, MonitorEnabledDefault)
	viper.SetDefault(MonitorFalsePositiveRateKey, MonitorFalsePositiveRateDefault)
	viper.SetDefault(MonitorLoadMonitoredDestinationsKey, MonitorLoadMonitoredDestinationsDefault)
	viper.SetDefault(MonitorMaxNumberOfDestinationsKey, MonitorMaxNumberOfDestinationsDefault)
	viper.SetDefault(MonitorMonitorDaysKey, MonitorMonitorDaysDefault)
	viper.SetDefault(MonitorProcessorTypeKey, MonitorProcessorTypeDefault)
	viper.SetDefault(MonitorSaveTransactionDestinationsKey, MonitorSaveTransactionDestinationsDefault)
}

func setNewRelicDefaults() {
	viper.SetDefault(NewRelicDomainNameKey, NewRelicDomainNameDefault)
	viper.SetDefault(NewRelicEnabledKey, NewRelicEnabledDefault)
	viper.SetDefault(NewRelicLicenseKeyKey, NewRelicLicenseKeyDefault)
}

func setNodesDefaults() {
	viper.SetDefault(NodesUseMapiFeeQuotesKey, NodesUseMapiFeeQuotesDefault)
	viper.SetDefault(NodesMinercraftAPIKey, NodesMinercraftAPIDefault)
	viper.SetDefault(NodesBroadcastClientAPIsKey, NodesBroadcastClientAPIsDefault)
}

func setNotificationsDefaults() {
	viper.SetDefault(NotificationsEnabledKey, NotificationsEnabledDefault)
	viper.SetDefault(NotificationsWebhookEndpointKey, NotificationsWebhookEndpointDefault)
}

func setPaymailDefaults() {
	viper.SetDefault(PaymailEnabledKey, PaymailEnabledDefault)
	viper.SetDefault(PaymailDefaultFromPaymailKey, PaymailDefaultFromPaymailDefault)
	viper.SetDefault(PaymailDefaultNoteKey, PaymailDefaultNoteDefault)
	viper.SetDefault(PaymailDomainValidationEnabledKey, PaymailDomainValidationEnabledDefault)
	viper.SetDefault(PaymailSenderValidationEnabledKey, PaymailSenderValidationEnabledDefault)
	viper.SetDefault(PaymailDomainsKey, PaymailDomainsDefault)
}

func setRedisDefaults() {
	viper.SetDefault(RedisDependencyModeKey, RedisDependencyModeDefault)
	viper.SetDefault(RedisMaxActiveConnectionsKey, RedisMaxActiveConnectionsDefault)
	viper.SetDefault(RedisMaxConnectionLifetimeKey, RedisMaxConnectionLifetimeDefault)
	viper.SetDefault(RedisMaxIdleConnectionsKey, RedisMaxIdleConnectionsDefault)
	viper.SetDefault(RedisMaxIdleTimeoutKey, RedisMaxIdleTimeoutDefault)
	viper.SetDefault(RedisURLKey, RedisURLDefault)
	viper.SetDefault(RedisUseTLSKey, RedisUseTLSDefault)
}

func setTaskManagerDefaults() {
	viper.SetDefault(TaskManagerFactoryKey, TaskManagerFactoryDefault)
}

func setServerDefaults() {
	viper.SetDefault(ServerIdleTimeoutKey, ServerIdleTimeoutDefault)
	viper.SetDefault(ServerReadTimeoutKey, ServerReadTimeoutDefault)
	viper.SetDefault(ServerWriteTimeoutKey, ServerWriteTimeoutDefault)
	viper.SetDefault(ServerPortKey, ServerPortDefault)
}

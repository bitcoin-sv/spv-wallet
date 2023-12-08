package config

import (
	"github.com/spf13/viper"
)

// General defaults
const (
	DebugDefault              = true
	DebugProfilingDefault     = false
	DisableITCDefault         = true
	ImportBlockHeadersDefault = ""
	RequestLoggingDefault     = true
)

// Authentication defaults
const (
	AuthAdminKeyDefault        = "xpub661MyMwAqRbcFaYeQLxmExXvTCjw9jjBRpifkoGggkAitXNNjva4TStLJuYjjEmU4AzXRPGwoECjXo3Rgqg8zQqW6UPVfkKtsrogGBw8xz7"
	AuthRequireSigningDefault  = false
	AuthSchemeDefault          = "xpub"
	AuthSigningDisabledDefault = true
)

// Beef defaults
const (
	UseBeefDefault                  = true
	PulseHeaderValidationURLDefault = "http://localhost:8000/api/v1/chain/merkleroot/verify"
	PulseAuthTokenDefault           = "asd"
)

// Cachestore defaults
const (
	CacheEngineDefault = "freecache"
)

// Cluster defaults
const (
	ClusterCoordinatorDefault         = "redis"
	ClusterPrefixDefault              = "bux_cluser_"
	ClusterRedisUrlDefault            = "localhost:6379"
	ClusterRedisMaxIdleTimeoutDefault = "10s"
	ClusterRedisUseTlsDefault         = false
)

// Datastore defaults
const (
	DatastoreAutoMigrateDefault = true
	DatastoreDebugDefault       = false
	DatastoreEngineDefault      = "sqlite"
	DatastoreTablePrefixDefault = "xapi"
)

// MongoDB config keys
const (
	MongoDatabaseNameDefault = "xapi"
	MongoTransactionsDefault = false
	MongoUriDefault          = "mongodb://localhost:27017/xapi"
)

// SQL (MySQL, PostgreSQL) config keys
const (
	SqlDriverDefault                    = "postgresql"
	SqlHostDefault                      = "localhost"
	SqlNameDefault                      = "xapi"
	SqlPasswordDefault                  = ""
	SqlPortDefault                      = "5432"
	SqlReplicaDefault                   = false
	SqlSkipInitializeWithVersionDefault = true
	SqlTimeZoneDefault                  = "UTC"
	SqlTxTimeoutDefault                 = "10s"
	SqlUserDefault                      = "postgres"
)

// SQLite config keys
const (
	SqliteDatabasePathDefault = "./test-json.db"
	SqliteSharedDefault       = true
)

// Graphql defaults
const (
	GraphqlEnabledDefault        = true
	GraphqlPlaygroundPathDefault = "/graphql"
	GraphqlServerPathDefault     = "/graphql"
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

var NodesBroadcastClientAPIsDefault = []string{"url|token"}

// Notification defaults
const (
	NotificationsEnabledDefault         = false
	NotificationsWebhookEndpointDefault = ""
)

// Paymail defaults
const (
	PaymailEnabledDefault                 = true
	PaymailDefaultFromPaymailDefault      = "from@domain.com"
	PaymailDefaultNoteDefault             = "bux Address Resolution"
	PaymailDomainValidationEnabledDefault = false
	PaymailSenderValidationEnabledDefault = true
)

var PaymailDomainsDefault = []string{"localhost"}

// Redis defaults
const (
	RedisDependencyModeDefault        = true
	RedisMaxActiveConnectionsDefault  = 0
	RedisMaxConnectionLifetimeDefault = "60s"
	RedisMaxIdleConnectionsDefault    = 10
	RedisMaxIdleTimeoutDefault        = "10s"
	RedisURLDefault                   = "redis://localhost:6379"
	RedisUseTLSDefault                = false
)

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

func setDefaults(configFilePath string) {
	viper.SetDefault(ConfigFilePathKey, configFilePath)

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
	viper.SetDefault(ClusterRedisUrlKey, ClusterRedisUrlDefault)
	viper.SetDefault(ClusterRedisMaxIdleTimeoutKey, ClusterRedisMaxIdleTimeoutDefault)
	viper.SetDefault(ClusterRedisUseTlsKey, ClusterRedisUseTlsDefault)
}

func setDbDefaults() {
	viper.SetDefault(DatastoreAutoMigrateKey, DatastoreAutoMigrateDefault)
	viper.SetDefault(DatastoreDebugKey, DatastoreDebugDefault)
	viper.SetDefault(DatastoreEngineKey, DatastoreEngineDefault)
	viper.SetDefault(DatastoreTablePrefixKey, DatastoreTablePrefixDefault)

	viper.SetDefault(MongoDatabaseNameKey, MongoDatabaseNameDefault)
	viper.SetDefault(MongoTransactionsKey, MongoTransactionsDefault)
	viper.SetDefault(MongoUriKey, MongoUriDefault)

	viper.SetDefault(SqlDriverKey, SqlDriverDefault)
	viper.SetDefault(SqlHostKey, SqlHostDefault)
	viper.SetDefault(SqlNameKey, SqlNameDefault)
	viper.SetDefault(SqlPasswordKey, SqlPasswordDefault)
	viper.SetDefault(SqlPortKey, SqlPortDefault)
	viper.SetDefault(SqlReplicaKey, SqlReplicaDefault)
	viper.SetDefault(SqlSkipInitializeWithVersionKey, SqlSkipInitializeWithVersionDefault)
	viper.SetDefault(SqlTimeZoneKey, SqlTimeZoneDefault)
	viper.SetDefault(SqlTxTimeoutKey, SqlTxTimeoutDefault)
	viper.SetDefault(SqlUserKey, SqlUserDefault)

	viper.SetDefault(SqliteDatabasePathKey, SqliteDatabasePathDefault)
	viper.SetDefault(SqliteSharedKey, SqliteSharedDefault)
}

func setGraphqlDefaults() {
	viper.SetDefault(GraphqlEnabledKey, GraphqlEnabledDefault)
	viper.SetDefault(GraphqlPlaygroundPathKey, GraphqlPlaygroundPathDefault)
	viper.SetDefault(GraphqlServerPathKey, GraphqlServerPathDefault)
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
	viper.SetDefault(TaskManagerEngineKey, TaskManagerEngineDefault)
	viper.SetDefault(TaskManagerFactoryKey, TaskManagerFactoryDefault)
	viper.SetDefault(TaskManagerQueueNameKey, TaskManagerQueueNameDefault)
}

func setServerDefaults() {
	viper.SetDefault(ServerIdleTimeoutKey, ServerIdleTimeoutDefault)
	viper.SetDefault(ServerReadTimeoutKey, ServerReadTimeoutDefault)
	viper.SetDefault(ServerWriteTimeoutKey, ServerWriteTimeoutDefault)
	viper.SetDefault(ServerPortKey, ServerPortDefault)
}

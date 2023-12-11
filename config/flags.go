package config

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type buxFlags struct {
	pflag.FlagSet
}

type cliFlags struct {
	showVersion bool `mapstructure:"version"`
	showHelp    bool `mapstructure:"help"`
	dumpConfig  bool `mapstructure:"dump_config"`
}

func loadFlags() error {
	if !anyFlagsPassed() {
		return nil
	}

	cli := cliFlags{}
	bux := buxFlags{}

	bux.initFlags(&cli)

	err := bux.Parse(os.Args[1:])
	if err != nil {
		fmt.Printf("Flags can't be parsed: %v\n", err)
		os.Exit(1)
	}

	err = viper.BindPFlags(&bux.FlagSet)
	if err != nil {
		return err
	}

	bux.parseCliFlags(&cli)

	return nil
}

func anyFlagsPassed() bool {
	return len(os.Args) > 1
}

func (fs *buxFlags) initFlags(cliFlags *cliFlags) {
	fs.StringP(ConfigFilePathKey, "C", "", "custom config file path")

	fs.initCliFlags(cliFlags)

	fs.initGeneralFlags()
	fs.initAuthFlags()
	fs.initBeefFlags()
	fs.initCachestoreFlags()
	fs.initClusterFlags()
	fs.initDbFlags()
	fs.initGraphqlFlags()
	fs.initMonitorFlags()
	fs.initNewRelicFlags()
	fs.initNodesFlags()
	fs.initNotificationFlags()
	fs.initPaymailFlags()
	fs.initRedisFlags()
	fs.initTaskManagerFlags()
	fs.initServerFlags()
}

func (fs *buxFlags) initCliFlags(cliFlags *cliFlags) {
	if cliFlags != nil {
		fs.BoolVarP(&cliFlags.showHelp, "help", "h", false, "show help")
		fs.BoolVarP(&cliFlags.showVersion, "version", "v", false, "show version")
		fs.BoolVarP(&cliFlags.dumpConfig, "dump_config", "d", false, "dump config to file, specified by config_file flag")
	}
}

func (fs *buxFlags) parseCliFlags(cli *cliFlags) {
	if cli.showHelp {
		fs.PrintDefaults()
		os.Exit(0)
	}

	if cli.showVersion {
		fmt.Println("bux-sever", "version", Version)
		os.Exit(0)
	}

	if cli.dumpConfig {
		configPath := viper.GetString(ConfigFilePathKey)
		if configPath == "" {
			configPath = DefaultConfigFilePath
		}
		err := viper.SafeWriteConfigAs(configPath)
		if err != nil {
			fmt.Printf("error while dumping config: %v", err.Error())
		}
		os.Exit(0)
	}
}

func (fs *buxFlags) initGeneralFlags() {
	fs.Bool(DebugKey, DebugDefault, "enable debug logging")
	fs.Bool(DebugProfilingKey, DebugProfilingDefault, "enable debug profiling")
	fs.Bool(DisableITCKey, DisableITCDefault, "disable ITC - Incoming Transaction Checking")
	fs.String(ImportBlockHeadersKey, ImportBlockHeadersDefault, "path or URL to blockheaders file")
	fs.Bool(RequestLoggingKey, RequestLoggingDefault, "request logging from api routers (rest and graphql)")
}

func (fs *buxFlags) initAuthFlags() {
	fs.String(AuthAdminKey, AuthAdminKeyDefault, "key that is used for administrative requests")
	fs.Bool(AuthRequireSigningKey, AuthRequireSigningDefault, "require signing")
	fs.String(AuthSchemeKey, AuthSchemeDefault, "authentication scheme to use")
	fs.Bool(AuthSigningDisabledKey, AuthSigningDisabledDefault, "NOTE: Only for development, turns off signing")
}

func (fs *buxFlags) initBeefFlags() {
	fs.Bool(UseBeefKey, UseBeefDefault, "enables BEEF transaction format, requires Pulse settings")
	fs.String(PulseHeaderValidationURLKey, PulseHeaderValidationURLDefault, "pulse url for validating merkle roots")
	fs.String(PulseAuthTokenKey, PulseAuthTokenDefault, "authentication token for pulse")
}

func (fs *buxFlags) initCachestoreFlags() {
	fs.String(CacheEngineKey, CacheEngineDefault, "cache engine: redis, freecache or empty")
}

func (fs *buxFlags) initClusterFlags() {
	fs.String(ClusterCoordinatorKey, ClusterCoordinatorDefault, "redis or memory")
	fs.String(ClusterPrefixKey, ClusterPrefixDefault, "prefix string to use for all cluster keys")
	fs.String(ClusterRedisURLKey, ClusterRedisURLDefault, "Redis URL for cluster coordinator, if redis is chosen")
	fs.String(ClusterRedisMaxIdleTimeoutKey, ClusterRedisMaxIdleTimeoutDefault, "max idle timeout for redis for cluster, if redis is chosen")
	fs.Bool(ClusterRedisUseTLSKey, ClusterRedisUseTLSDefault, "should redis cluster coordinator use tls, if redis is chosen")
}

func (fs *buxFlags) initDbFlags() {
	fs.Bool(DatastoreAutoMigrateKey, DatastoreAutoMigrateDefault, "loads a blank database")
	fs.Bool(DatastoreDebugKey, DatastoreDebugDefault, "show sql statements")
	fs.String(DatastoreEngineKey, DatastoreEngineDefault, "mysql, sqlite, postgresql, mongodb, empty")
	fs.String(DatastoreTablePrefixKey, DatastoreTablePrefixDefault, "prefix for all tables in db")

	fs.String(MongoDatabaseNameKey, MongoDatabaseNameDefault, "database name for MongoDB")
	fs.Bool(MongoTransactionsKey, MongoTransactionsDefault, "has transactions")
	fs.String(MongoURIKey, MongoURIDefault, "connection uri to MongoDB")

	fs.String(SQLDriverKey, SQLDriverDefault, "mysql, postgresql")
	fs.String(SQLHostKey, SQLHostDefault, "db host")
	fs.String(SQLUserKey, SQLUserDefault, "db user")
	fs.String(SQLNameKey, SQLNameDefault, "db name")
	fs.String(SQLPasswordKey, SQLPasswordDefault, "db password")
	fs.String(SQLPortKey, SQLPortDefault, "db port")
	fs.Bool(SQLReplicaKey, SQLReplicaDefault, "true if it's a replica (Read-Only)")
	fs.Bool(SQLSkipInitializeWithVersionKey, SQLSkipInitializeWithVersionDefault, "skip using MySQL in test mode")
	fs.String(SQLTimeZoneKey, SQLTimeZoneDefault, "time zone for db")
	fs.String(SQLTxTimeoutKey, SQLTxTimeoutDefault, "timeout for transactions")

	fs.String(SQLiteDatabasePathKey, SQLiteDatabasePathDefault, "db path for sqlite")
	fs.Bool(SQLiteSharedKey, SQLiteSharedDefault, "adds a shared param to the connection string")
}

func (fs *buxFlags) initGraphqlFlags() {
	fs.Bool(GraphqlEnabledKey, GraphqlEnabledDefault, "enable graphql")
}

func (fs *buxFlags) initMonitorFlags() {
	fs.String(MonitorAuthTokenKey, MonitorAuthTokenDefault, "token to connect to the server with")
	fs.String(MonitorBuxAgentURLKey, MonitorBuxAgentURLDefault, "the bux agent server url address")
	fs.Bool(MonitorDebugKey, MonitorDebugDefault, "enable debug")
	fs.Bool(MonitorEnabledKey, MonitorEnabledDefault, "enable monitor")
	fs.Float64(MonitorFalsePositiveRateKey, MonitorFalsePositiveRateDefault, "percentage of false positives to expect")
	fs.Bool(MonitorLoadMonitoredDestinationsKey, MonitorLoadMonitoredDestinationsDefault, "load monitored destinations")
	fs.Int(MonitorMaxNumberOfDestinationsKey, MonitorMaxNumberOfDestinationsDefault, "number of destinations that the filter can hold")
	fs.Int(MonitorMonitorDaysKey, MonitorMonitorDaysDefault, "number of days in the past that an address should be monitored for")
	fs.String(MonitorProcessorTypeKey, MonitorProcessorTypeDefault, "type of processor to start monitor with")
	fs.Bool(MonitorSaveTransactionDestinationsKey, MonitorSaveTransactionDestinationsDefault, "save destinations on monitored transactions")
}

func (fs *buxFlags) initNewRelicFlags() {
	fs.String(NewRelicDomainNameKey, NewRelicDomainNameDefault, "used for hostname display")
	fs.Bool(NewRelicEnabledKey, NewRelicEnabledDefault, "enable NewRelic")
	fs.String(NewRelicLicenseKeyKey, NewRelicLicenseKeyDefault, "license key")
}

func (fs *buxFlags) initNodesFlags() {
	fs.Bool(NodesUseMapiFeeQuotesKey, NodesUseMapiFeeQuotesDefault, "use mAPI fee quotes")
	fs.String(NodesMinercraftAPIKey, NodesMinercraftAPIDefault, "type of api to use by minercraft, arc of mapi")
	fs.StringSlice(NodesBroadcastClientAPIsKey, NodesBroadcastClientAPIsDefault, "go-broadcastClient api keys in fromat 'api_url|token'")
}

func (fs *buxFlags) initNotificationFlags() {
	fs.Bool(NotificationsEnabledKey, NotificationsEnabledDefault, "enable notifications")
	fs.String(NotificationsWebhookEndpointKey, NotificationsWebhookEndpointDefault, "webhook endpoint for notifications")
}

func (fs *buxFlags) initPaymailFlags() {
	fs.String(PaymailDefaultFromPaymailKey, PaymailDefaultFromPaymailDefault, "default 'from:@domain.com' paymail")
	fs.String(PaymailDefaultNoteKey, PaymailDefaultNoteDefault, "default paymail note, IE: message needed for address resolution")
	fs.StringSlice(PaymailDomainsKey, PaymailDomainsDefault, "list of allowed paymail domains")
	fs.Bool(PaymailDomainValidationEnabledKey, PaymailDomainValidationEnabledDefault, "enable paymail domain validation, turn off if hosted domain is not paymail related")
	fs.Bool(PaymailEnabledKey, PaymailEnabledDefault, "enable paymail")
	fs.Bool(PaymailSenderValidationEnabledKey, PaymailSenderValidationEnabledDefault, "enable paymail sender validation - extra security")
}

func (fs *buxFlags) initRedisFlags() {
	fs.Bool(RedisDependencyModeKey, RedisDependencyModeDefault, "only in Redis with script enabled")
	fs.Int(RedisMaxActiveConnectionsKey, RedisMaxActiveConnectionsDefault, "max active redis connections")
	fs.String(RedisMaxConnectionLifetimeKey, RedisMaxConnectionLifetimeDefault, "max redis connection lifetime")
	fs.Int(RedisMaxIdleConnectionsKey, RedisMaxIdleConnectionsDefault, "max idle redis connections")
	fs.String(RedisMaxIdleTimeoutKey, RedisMaxIdleTimeoutDefault, "max idle redis timeout")
	fs.String(RedisURLKey, RedisURLDefault, "redis url connections string")
	fs.Bool(RedisUseTLSKey, RedisUseTLSDefault, "enable redis TLS")
}

func (fs *buxFlags) initTaskManagerFlags() {
	fs.String(TaskManagerFactoryKey, TaskManagerFactoryDefault, "memory, redis, empty")
}

func (fs *buxFlags) initServerFlags() {
	fs.String(ServerIdleTimeoutKey, ServerIdleTimeoutDefault, "server idle timeout")
	fs.String(ServerReadTimeoutKey, ServerReadTimeoutDefault, "server read timeout")
	fs.String(ServerWriteTimeoutKey, ServerWriteTimeoutDefault, "server write timout")
	fs.String(ServerPortKey, ServerPortDefault, "server port")
}

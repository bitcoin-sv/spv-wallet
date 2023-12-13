package config

import (
	"time"

	"github.com/mrz1836/go-datastore"
	"github.com/tonicpow/go-minercraft/v2"
)

// DefaultAppConfig are the default values for AppConfig
var DefaultAppConfig = &AppConfig{
	Authentication:     authConfigDefault,
	Cache:              cacheDefault,
	Db:                 dbDefaut,
	Debug:              true,
	DebugProfiling:     true,
	DisableITC:         true,
	GraphQL:            graphqlDefault,
	ImportBlockHeaders: "",
	Monitor:            monitorDefault,
	NewRelic:           newRelicDefault,
	Nodes:              nodesDefault,
	Notifications:      notificationDefault,
	Paymail:            paymailDefault,
	RequestLogging:     true,
	Server:             serverDefault,
	TaskManager:        taskManagerDefault,
}

var authConfigDefault = &AuthenticationConfig{
	AdminKey:        "xpub661MyMwAqRbcFrBJbKwBGCB7d3fr2SaAuXGM95BA62X41m6eW2ehRQGW4xLi9wkEXUGnQZYxVVj4PxXnyrLk7jdqvBAs1Qq9gf6ykMvjR7J",
	RequireSigning:  false,
	Scheme:          "xpub",
	SigningDisabled: true,
}

var cacheDefault = &CacheConfig{
	Engine: "freecache",
	Cluster: &ClusterConfig{
		Coordinator: "memory",
		Prefix:      "bux_cluster_",
		Redis:       nil,
	},
	Redis: &RedisConfig{
		DependencyMode:        true,
		MaxActiveConnections:  0,
		MaxConnectionLifetime: 60 * time.Second,
		MaxIdleConnections:    10,
		MaxIdleTimeout:        10 * time.Second,
		URL:                   "redis://localhost:6379",
		UseTLS:                false,
	},
}

var dbDefaut = &DbConfig{
	Datastore: &DatastoreConfig{
		Debug:       false,
		Engine:      "sqlite",
		TablePrefix: "xapi",
	},
	Mongo: &datastore.MongoDBConfig{
		DatabaseName:       "xapi",
		ExistingConnection: nil,
		Transactions:       false,
		URI:                "mongodb://localhost:27017/xapi",
	},
	SQL: &datastore.SQLConfig{
		Driver:                    "postgresql",
		ExistingConnection:        nil,
		Host:                      "localhost",
		Name:                      "xapi",
		Password:                  "",
		Port:                      "5432",
		Replica:                   false,
		SkipInitializeWithVersion: true,
		TimeZone:                  "UTC",
		TxTimeout:                 10 * time.Second,
		User:                      "postgres",
	},
	SQLite: &datastore.SQLiteConfig{
		DatabasePath:       "./bux.db",
		ExistingConnection: nil,
		Shared:             true,
	},
}

var graphqlDefault = &GraphqlConfig{
	Enabled: false,
}

var monitorDefault = &MonitorOptions{
	AuthToken:                   "",
	BuxAgentURL:                 "ws://localhost:8000/websocket",
	Debug:                       false,
	Enabled:                     false,
	FalsePositiveRate:           0.01,
	LoadMonitoredDestinations:   false,
	MaxNumberOfDestinations:     100000,
	MonitorDays:                 7,
	ProcessorType:               "bloom",
	SaveTransactionDestinations: true,
}

var newRelicDefault = &NewRelicConfig{
	DomainName: "domain.com",
	Enabled:    false,
	LicenseKey: "BOGUS-LICENSE-KEY-1234567890987654321234",
}

var nodesDefault = &NodesConfig{
	UseMapiFeeQuotes:     true,
	MinercraftAPI:        "mAPI",
	MinercraftCustomAPIs: []*minercraft.MinerAPIs{},
	BroadcastClientAPIs:  []string{},
}

var notificationDefault = &NotificationsConfig{
	Enabled:         false,
	WebhookEndpoint: "",
}

var paymailDefault = &PaymailConfig{
	Beef: &BeefConfig{
		UseBeef:                  false,
		PulseHeaderValidationURL: "http://localhost:8080/api/v1/chain/merkleroot/verify",
		PulseAuthToken:           "mQZQ6WmxURxWz5ch", // #nosec G101
	},
	DefaultFromPaymail:      "from@domain.com",
	DefaultNote:             "bux Address Resolution",
	Domains:                 []string{"localhost"},
	DomainValidationEnabled: true,
	Enabled:                 true,
	SenderValidationEnabled: true,
}

var taskManagerDefault = &TaskManagerConfig{
	Factory: "memory",
}

var serverDefault = &ServerConfig{
	IdleTimeout:  60 * time.Second,
	ReadTimeout:  15 * time.Second,
	WriteTimeout: 15 * time.Second,
	Port:         "3003",
}

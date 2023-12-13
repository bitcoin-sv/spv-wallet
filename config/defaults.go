package config

import (
	"time"

	"github.com/mrz1836/go-datastore"
	"github.com/tonicpow/go-minercraft/v2"
)

// DefaultAppConfig is the default config for AppConfig
var DefaultAppConfig = &AppConfig{
	Authentication:     AuthConfigDefault,
	Cache:              CacheDefault,
	Db:                 DbDefaut,
	Debug:              true,
	DebugProfiling:     true,
	DisableITC:         true,
	GraphQL:            GraphqlDefault,
	ImportBlockHeaders: "",
	Monitor:            MonitorDefault,
	NewRelic:           NewRelicDefault,
	Nodes:              NodesDefault,
	Notifications:      NotificationDefault,
	Paymail:            PaymailDefault,
	RequestLogging:     true,
	Server:             ServerDefault,
	TaskManager:        TaskManagerDefault,
}

// AuthConfigDefault is the default config for AuthenticationConfig
var AuthConfigDefault = &AuthenticationConfig{
	AdminKey:        "xpub661MyMwAqRbcFrBJbKwBGCB7d3fr2SaAuXGM95BA62X41m6eW2ehRQGW4xLi9wkEXUGnQZYxVVj4PxXnyrLk7jdqvBAs1Qq9gf6ykMvjR7J",
	RequireSigning:  false,
	Scheme:          "xpub",
	SigningDisabled: true,
}

// CacheDefault is the default config for CacheConfig
var CacheDefault = &CacheConfig{
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

// DbDefaut is the default config for DbConfig
var DbDefaut = &DbConfig{
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

// GraphqlDefault is the default settings for GraphqlConfig
var GraphqlDefault = &GraphqlConfig{
	Enabled: false,
}

// MonitorDefault is the default settings for MonitorOptions
var MonitorDefault = &MonitorOptions{
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

// NewRelicDefault is the default settings for NewRelicConfig
var NewRelicDefault = &NewRelicConfig{
	DomainName: "domain.com",
	Enabled:    false,
	LicenseKey: "BOGUS-LICENSE-KEY-1234567890987654321234",
}

// NodesDefault is the default settings for NodesConfig
var NodesDefault = &NodesConfig{
	UseMapiFeeQuotes:     true,
	MinercraftAPI:        "mAPI",
	MinercraftCustomAPIs: []*minercraft.MinerAPIs{},
	BroadcastClientAPIs:  []string{},
}

// NotificationDefault is the default settings for NotificationConfig
var NotificationDefault = &NotificationsConfig{
	Enabled:         false,
	WebhookEndpoint: "",
}

// PaymailDefault is the default settings for PaymailConfig
var PaymailDefault = &PaymailConfig{
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

// TaskManagerDefault is the default settings for TaskManagerConfig
var TaskManagerDefault = &TaskManagerConfig{
	Factory: "memory",
}

// ServerDefault is the default settings for ServerConfig
var ServerDefault = &ServerConfig{
	IdleTimeout:  60 * time.Second,
	ReadTimeout:  15 * time.Second,
	WriteTimeout: 15 * time.Second,
	Port:         "3003",
}

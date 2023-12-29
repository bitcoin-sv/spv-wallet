package config

import (
	"time"

	"github.com/mrz1836/go-datastore"
	"github.com/tonicpow/go-minercraft/v2"
)

func getDefaultAppConfig() *AppConfig {
	return &AppConfig{
		Authentication:     getAuthConfigDefaults(),
		Cache:              getCacheDefaults(),
		Db:                 getDbDefaults(),
		Debug:              true,
		DebugProfiling:     true,
		DisableITC:         true,
		GraphQL:            getGraphqlDefaults(),
		ImportBlockHeaders: "",
		Logging:            getLoggingDefaults(),
		NewRelic:           getNewRelicDefaults(),
		Nodes:              getNodesDefaults(),
		Notifications:      getNotificationDefaults(),
		Paymail:            getPaymailDefaults(),
		RequestLogging:     true,
		Server:             getServerDefaults(),
		TaskManager:        getTaskManagerDefault(),
	}
}

func getAuthConfigDefaults() *AuthenticationConfig {
	return &AuthenticationConfig{
		AdminKey:        "xpub661MyMwAqRbcFrBJbKwBGCB7d3fr2SaAuXGM95BA62X41m6eW2ehRQGW4xLi9wkEXUGnQZYxVVj4PxXnyrLk7jdqvBAs1Qq9gf6ykMvjR7J",
		RequireSigning:  false,
		Scheme:          "xpub",
		SigningDisabled: true,
	}
}

func getCacheDefaults() *CacheConfig {
	return &CacheConfig{
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
}

func getDbDefaults() *DbConfig {
	return &DbConfig{
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
}

func getGraphqlDefaults() *GraphqlConfig {
	return &GraphqlConfig{
		Enabled: false,
	}
}

func getLoggingDefaults() *LoggingConfig {
	return &LoggingConfig{
		Level:        "debug",
		Format:       "console",
		InstanceName: "bux-server",
		LogOrigin:    false,
	}
}

func getNewRelicDefaults() *NewRelicConfig {
	return &NewRelicConfig{
		DomainName: "domain.com",
		Enabled:    false,
		LicenseKey: "BOGUS-LICENSE-KEY-1234567890987654321234",
	}
}

func getNodesDefaults() *NodesConfig {
	return &NodesConfig{
		UseMapiFeeQuotes:     true,
		MinercraftAPI:        "mAPI",
		MinercraftCustomAPIs: []*minercraft.MinerAPIs{},
		BroadcastClientAPIs:  []*BroadcastClientAPI{},
	}
}

func getNotificationDefaults() *NotificationsConfig {
	return &NotificationsConfig{
		Enabled:         false,
		WebhookEndpoint: "",
	}
}

func getPaymailDefaults() *PaymailConfig {
	return &PaymailConfig{
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
}

func getTaskManagerDefault() *TaskManagerConfig {
	return &TaskManagerConfig{
		Factory: "memory",
	}
}

func getServerDefaults() *ServerConfig {
	return &ServerConfig{
		IdleTimeout:  60 * time.Second,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		Port:         3003,
	}
}

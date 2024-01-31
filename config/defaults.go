package config

import (
	"time"

	"github.com/mrz1836/go-datastore"
)

const DefaultAdminXpub = "xpub661MyMwAqRbcFgfmdkPgE2m5UjHXu9dj124DbaGLSjaqVESTWfCD4VuNmEbVPkbYLCkykwVZvmA8Pbf8884TQr1FgdG2nPoHR8aB36YdDQh"

func getDefaultAppConfig() *AppConfig {
	return &AppConfig{
		Authentication:     getAuthConfigDefaults(),
		Cache:              getCacheDefaults(),
		Callback:           getCallbackDefaults(),
		Db:                 getDbDefaults(),
		Debug:              true,
		DebugProfiling:     true,
		DisableITC:         true,
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
		AdminKey:        DefaultAdminXpub,
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

func getCallbackDefaults() *CallbackConfig {
	return &CallbackConfig{
		CallbackHost:  "http://localhost:3003",
		CallbackToken: "",
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
			SslMode:                   "disable",
		},
		SQLite: &datastore.SQLiteConfig{
			DatabasePath:       "./bux.db",
			ExistingConnection: nil,
			Shared:             true,
		},
	}
}

func getLoggingDefaults() *LoggingConfig {
	return &LoggingConfig{
		Level:        "info",
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
		Protocol: NodesProtocolArc,
		Apis: []*MinerAPI{
			{
				ArcURL: "https://arc.gorillapool.io",
				// GorillaPool does not support querying (Merkle proofs)
				Token:   "",
				MinerID: "03ad780153c47df915b3d2e23af727c68facaca4facd5f155bf5018b979b9aeb83",
			},
		},
		UseFeeQuotes: true,
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
			UseBeef:                  true,
			PulseHeaderValidationURL: "http://localhost:8080/api/v1/chain/merkleroot/verify",
			PulseAuthToken:           "mQZQ6WmxURxWz5ch", // #nosec G101
		},
		DefaultFromPaymail:      "from@domain.com",
		Domains:                 []string{"localhost"},
		DomainValidationEnabled: true,
		SenderValidationEnabled: false,
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

package config

import (
	"time"

	"github.com/google/uuid"
	"github.com/mrz1836/go-datastore"
)

// DefaultAdminXpub is the default admin xpub used for authenticate requests.
const DefaultAdminXpub = "xpub661MyMwAqRbcFgfmdkPgE2m5UjHXu9dj124DbaGLSjaqVESTWfCD4VuNmEbVPkbYLCkykwVZvmA8Pbf8884TQr1FgdG2nPoHR8aB36YdDQh"

func getDefaultAppConfig() *AppConfig {
	return &AppConfig{
		Authentication:     getAuthConfigDefaults(),
		Cache:              getCacheDefaults(),
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
		Metrics:            getMetricsDefaults(),
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
			Prefix:      "spv_cluster_",
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
			DatabasePath:       "./spv.db",
			ExistingConnection: nil,
			Shared:             true,
		},
	}
}

func getLoggingDefaults() *LoggingConfig {
	return &LoggingConfig{
		Level:        "info",
		Format:       "console",
		InstanceName: "spv-wallet",
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
	depIDSufix, _ := uuid.NewUUID()
	return &NodesConfig{
		DeploymentID: "spv-" + depIDSufix.String(),
		Protocol:     NodesProtocolArc,
		Callback:     getCallbackDefaults(),
		Apis: []*MinerAPI{
			{
				ArcURL:  "https://api.taal.com/arc",
				Token:   "mainnet_06770f425eb00298839a24a49cbdc02c",
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
			UseBeef:                               true,
			BlockHeaderServiceHeaderValidationURL: "http://localhost:8080/api/v1/chain/merkleroot/verify",
			BlockHeaderServiceAuthToken:           "mQZQ6WmxURxWz5ch", // #nosec G101
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

func getMetricsDefaults() *MetricsConfig {
	return &MetricsConfig{
		Enabled: false,
	}
}

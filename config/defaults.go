package config

import (
	"time"

	"github.com/bitcoin-sv/spv-wallet/engine/datastore"
	"github.com/google/uuid"
)

// DefaultAdminXpub is the default admin xpub used for authenticate requests.
const DefaultAdminXpub = "xpub661MyMwAqRbcFgfmdkPgE2m5UjHXu9dj124DbaGLSjaqVESTWfCD4VuNmEbVPkbYLCkykwVZvmA8Pbf8884TQr1FgdG2nPoHR8aB36YdDQh"

// TaskManagerQueueName is the default queue name for the task manager.
const TaskManagerQueueName = "spv_wallet_queue"

func GetDefaultAppConfig() *AppConfig {
	return &AppConfig{
		Version:              "development",
		Authentication:       getAuthConfigDefaults(),
		Cache:                getCacheDefaults(),
		Db:                   getDbDefaults(),
		DebugProfiling:       true,
		DisableITC:           true,
		ImportBlockHeaders:   "",
		Logging:              getLoggingDefaults(),
		ARC:                  getARCDefaults(),
		Notifications:        getNotificationDefaults(),
		Paymail:              getPaymailDefaults(),
		BHS:                  getBHSDefaults(),
		RequestLogging:       true,
		Server:               getServerDefaults(),
		TaskManager:          getTaskManagerDefault(),
		Metrics:              getMetricsDefaults(),
		ExperimentalFeatures: getExperimentalFeaturesConfig(),
	}
}

func getAuthConfigDefaults() *AuthenticationConfig {
	return &AuthenticationConfig{
		AdminKey:       DefaultAdminXpub,
		RequireSigning: false,
		Scheme:         "xpub",
	}
}

func getCacheDefaults() *CacheConfig {
	return &CacheConfig{
		Engine: "freecache",
		Cluster: &ClusterConfig{
			Coordinator: "memory",
			Prefix:      "spv_wallet_cluster_",
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
		SQL: &datastore.SQLConfig{
			Driver:             "postgresql",
			ExistingConnection: nil,
			Host:               "localhost",
			Name:               "xapi",
			Password:           "",
			Port:               "5432",
			Replica:            false,
			TimeZone:           "UTC",
			TxTimeout:          10 * time.Second,
			User:               "postgres",
			SslMode:            "disable",
		},
		SQLite: &datastore.SQLiteConfig{
			DatabasePath:       "./spv-wallet.db",
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

func getARCDefaults() *ARCConfig {
	depIDSufix, _ := uuid.NewUUID()
	return &ARCConfig{
		DeploymentID: "spv-wallet-" + depIDSufix.String(),
		URL:          "https://arc.taal.com",
		Token:        "mainnet_06770f425eb00298839a24a49cbdc02c",
		UseFeeQuotes: true,
		Callback: &CallbackConfig{
			Enabled: false,
			Host:    "https://example.com",
			Token:   "",
		},
	}
}

func getNotificationDefaults() *NotificationsConfig {
	return &NotificationsConfig{
		Enabled: true,
	}
}

func getPaymailDefaults() *PaymailConfig {
	return &PaymailConfig{
		Beef: &BeefConfig{
			UseBeef:                                true,
			BlockHeadersServiceHeaderValidationURL: "http://localhost:8080/api/v1/chain/merkleroot/verify",
			BlockHeadersServiceAuthToken:           "mQZQ6WmxURxWz5ch", // #nosec G101
		},
		DefaultFromPaymail:      "from@domain.com",
		Domains:                 []string{"localhost"},
		DomainValidationEnabled: true,
		SenderValidationEnabled: false,
	}
}

func getBHSDefaults() *BHSConfig {
	return &BHSConfig{
		AuthToken: "mQZQ6WmxURxWz5ch",
		URL:       "http://localhost:8080",
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

func getExperimentalFeaturesConfig() *ExperimentalConfig {
	return &ExperimentalConfig{
		PikeContactsEnabled: false,
		PikePaymentEnabled:  false,
	}
}

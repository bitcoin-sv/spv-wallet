// Package config provides a configuration for the API
package config

import (
	"errors"
	"time"

	"github.com/BuxOrg/bux/cachestore"
	"github.com/BuxOrg/bux/datastore"
	"github.com/BuxOrg/bux/taskmanager"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/newrelic/go-agent/v3/newrelic"
)

// Config constants used for optimization and value testing
const (
	ApplicationName                = "BuxServer"
	CurrentMajorVersion            = "v1"
	DefaultHTTPRequestIdleTimeout  = 60 * time.Second
	DefaultHTTPRequestReadTimeout  = 15 * time.Second
	DefaultHTTPRequestWriteTimeout = 15 * time.Second
	DefaultNewRelicShutdown        = 10 * time.Second
	EnvironmentDevelopment         = "development"
	EnvironmentKey                 = "BUX_ENVIRONMENT"
	EnvironmentPrefix              = "bux"
	EnvironmentProduction          = "production"
	EnvironmentStaging             = "staging"
	EnvironmentTest                = "test"
	HealthRequestPath              = "health"
	Version                        = "v0.1.4"
)

// Local variables for configuration
var (
	environments = []interface{}{
		EnvironmentDevelopment,
		EnvironmentProduction,
		EnvironmentStaging,
		EnvironmentTest,
	}
)

// The global configuration settings
type (

	// AppConfig is the configuration values and associated env vars
	AppConfig struct {
		Authentication   *AuthenticationConfig    `json:"authentication" mapstructure:"authentication"`
		Cachestore       *CachestoreConfig        `json:"cache" mapstructure:"cache"`
		Datastore        *DatastoreConfig         `json:"datastore" mapstructure:"datastore"`
		Debug            bool                     `json:"debug" mapstructure:"debug"`
		DebugProfiling   bool                     `json:"debug_profiling" mapstructure:"debug_profiling"`
		DisableITC       bool                     `json:"disable_itc" mapstructure:"disable_itc"`
		Environment      string                   `json:"environment" mapstructure:"environment"`
		GDPRCompliance   bool                     `json:"gdpr_compliance" mapstructure:"gdpr_compliance"`
		GraphQL          *GraphqlConfig           `json:"graphql" mapstructure:"graphql"`
		Mongo            *datastore.MongoDBConfig `json:"mongodb" mapstructure:"mongodb"`
		Monitor          *MonitorOptions          `json:"monitor" mapstructure:"monitor"`
		NewRelic         *NewRelicConfig          `json:"new_relic" mapstructure:"new_relic"`
		Paymail          *PaymailConfig           `json:"paymail" mapstructure:"paymail"`
		Redis            *RedisConfig             `json:"redis" mapstructure:"redis"`
		RequestLogging   bool                     `json:"request_logging" mapstructure:"request_logging"`
		Server           *ServerConfig            `json:"server" mapstructure:"server"`
		SQL              *datastore.SQLConfig     `json:"sql" mapstructure:"sql"`
		SQLite           *datastore.SQLiteConfig  `json:"sqlite" mapstructure:"sqlite"`
		TaskManager      *TaskManagerConfig       `json:"task_manager" mapstructure:"task_manager"`
		WorkingDirectory string                   `json:"working_directory" mapstructure:"working_directory"`
	}

	// AuthenticationConfig is the configuration for Authentication
	AuthenticationConfig struct {
		AdminKey        string `json:"admin_key" mapstructure:"admin_key"`               // key that is used for administrative requests
		RequireSigning  bool   `json:"require_signing" mapstructure:"require_signing"`   // if the signing is required
		Scheme          string `json:"scheme" mapstructure:"scheme"`                     // authentication scheme to use (default is: xpub)
		SigningDisabled bool   `json:"signing_disabled" mapstructure:"signing_disabled"` // NOTE: Only for development (turns off signing)
	}

	// CachestoreConfig is a configuration for cachestore
	CachestoreConfig struct {
		Engine cachestore.Engine `json:"engine" mapstructure:"engine"` // Cache engine to use (redis, mcache, ristretto)
	}

	// DatastoreConfig is a configuration for the datastore
	DatastoreConfig struct {
		AutoMigrate bool             `json:"auto_migrate" mapstructure:"auto_migrate"` // loads a blank database
		Debug       bool             `json:"debug" mapstructure:"debug"`               // true for sql statements
		Engine      datastore.Engine `json:"engine" mapstructure:"engine"`             // mysql, sqlite
		TablePrefix string           `json:"table_prefix" mapstructure:"table_prefix"` // pre_users (pre)
	}

	// GraphqlConfig is the configuration for the GraphQL server
	GraphqlConfig struct {
		Enabled        bool   `json:"enabled" mapstructure:"enabled"`                 // true/false
		PlaygroundPath string `json:"playground_path" mapstructure:"playground_path"` // playground path i.e. "/graphiql"
		ServerPath     string `json:"server_path" mapstructure:"server_path"`         // server path i.e. "/graphql"
	}

	// MonitorOptions is the configuration for blockchain monitoring
	MonitorOptions struct {
		Enabled                 bool    `json:"enabled" mapstructure:"enabled"`                                       // true/false
		CentrifugeServer        string  `json:"centrifuge_server" mapstructure:"centrifuge_server"`                   // how many days in the past should we monitor an address (default: 7)
		MonitorDays             int     `json:"monitor_days" mapstructure:"monitor_days"`                             // how many days in the past should we monitor an address (default: 7)
		FalsePositiveRate       float64 `json:"false_positive_rate" mapstructure:"false_positive_rate"`               // how many false positives do we except (default: 0.01)
		MaxNumberOfDestinations int     `json:"max_number_of_destinations" mapstructure:"max_number_of_destinations"` // how many destinations can the filter hold (default: 100,000)
	}

	// NewRelicConfig is the configuration for New Relic
	NewRelicConfig struct {
		DomainName string `json:"domain_name" mapstructure:"domain_name"` // used for hostname display
		Enabled    bool   `json:"enabled" mapstructure:"enabled"`         // true/false
		LicenseKey string `json:"license_key" mapstructure:"license_key"` // 2342-3423523-62
	}

	// PaymailConfig is the configuration for the built-in Paymail server
	PaymailConfig struct {
		DefaultFromPaymail      string   `json:"default_from_paymail" mapstructure:"default_from_paymail"`           // IE: from@domain.com
		DefaultNote             string   `json:"default_note" mapstructure:"default_note"`                           // IE: message needed for address resolution
		Domains                 []string `json:"domains" mapstructure:"domains"`                                     // List of allowed domains
		DomainValidationEnabled bool     `json:"domain_validation_enabled" mapstructure:"domain_validation_enabled"` // Turn off if hosted domain is not paymail related
		Enabled                 bool     `json:"enabled" mapstructure:"enabled"`                                     // Flag for enabling the Paymail Server Service
		SenderValidationEnabled bool     `json:"sender_validation_enabled" mapstructure:"sender_validation_enabled"` // Turn on extra security
	}

	// RedisConfig is a configuration for Redis cachestore or taskmanager
	RedisConfig struct {
		DependencyMode        bool          `json:"dependency_mode" mapstructure:"dependency_mode"`                 // Only in Redis with script enabled
		MaxActiveConnections  int           `json:"max_active_connections" mapstructure:"max_active_connections"`   // Max active connections
		MaxConnectionLifetime time.Duration `json:"max_connection_lifetime" mapstructure:"max_connection_lifetime"` // Max connection lifetime
		MaxIdleConnections    int           `json:"max_idle_connections" mapstructure:"max_idle_connections"`       // Max idle connections
		MaxIdleTimeout        time.Duration `json:"max_idle_timeout" mapstructure:"max_idle_timeout"`               // Max idle timeout
		URL                   string        `json:"url" mapstructure:"url"`                                         // Redis URL connection string
		UseTLS                bool          `json:"use_tls" mapstructure:"use_tls"`                                 // Flag for using TLS
	}

	// TaskManagerConfig is a configuration for the taskmanager
	TaskManagerConfig struct {
		// QueueOptions *taskq.QueueOptions
		Engine    taskmanager.Engine  `json:"engine" mapstructure:"engine"`         // taskq, machinery
		Factory   taskmanager.Factory `json:"factory" mapstructure:"factory"`       // Factory (memory, redis)
		QueueName string              `json:"queue_name" mapstructure:"queue_name"` // test_queue
	}

	// ServerConfig is a configuration for the HTTP Server
	ServerConfig struct {
		IdleTimeout  time.Duration `json:"idle_timeout" mapstructure:"idle_timeout"`   // 60s
		ReadTimeout  time.Duration `json:"read_timeout" mapstructure:"read_timeout"`   // 15s
		WriteTimeout time.Duration `json:"write_timeout" mapstructure:"write_timeout"` // 15s
		Port         string        `json:"port" mapstructure:"port"`                   // 3003
	}
)

// GetUserAgent will return the outgoing user agent
func (a *AppConfig) GetUserAgent() string {
	return "BUX-Server " + a.Environment + " " + Version
}

// Validate checks the configuration for specific rules
func (a *AppConfig) Validate(txn *newrelic.Transaction) error {

	var err error
	defer txn.StartSegment("config_validation").End()

	if err = a.Authentication.Validate(); err != nil {
		return err
	}

	if err = a.Cachestore.Validate(); err != nil {
		return err
	}

	if err = a.Datastore.Validate(); err != nil {
		return err
	}

	if err = a.NewRelic.Validate(); err != nil {
		return err
	}

	if err = a.Paymail.Validate(); err != nil {
		return err
	}

	if err = a.Server.Validate(); err != nil {
		return err
	}

	if err = a.validateCachestore(); err != nil {
		return err
	}

	if err = a.validateDatastore(); err != nil {
		return err
	}

	return validation.ValidateStruct(a,
		validation.Field(&a.Environment, validation.Required, validation.In(environments...)),
		validation.Field(&a.WorkingDirectory, validation.Required),
	)
}

// validateDatastore will check the datastore and validate basic requirements
func (a *AppConfig) validateDatastore() error {
	if a.Datastore.Engine == datastore.SQLite {
		if a.SQLite == nil {
			return errors.New("missing sqlite config")
		}
	} else if a.Datastore.Engine == datastore.MySQL || a.Datastore.Engine == datastore.PostgreSQL {
		if a.SQL == nil {
			return errors.New("missing sql config")
		} else if len(a.SQL.Host) == 0 {
			return errors.New("missing sql host")
		} else if len(a.SQL.User) == 0 {
			return errors.New("missing sql username")
		} else if len(a.SQL.Name) == 0 {
			return errors.New("missing sql db name")
		}
	} else if a.Datastore.Engine == datastore.MongoDB {
		if a.Mongo == nil {
			return errors.New("missing mongo config")
		} else if len(a.Mongo.URI) == 0 {
			return errors.New("missing mongo uri")
		} else if len(a.Mongo.DatabaseName) == 0 {
			return errors.New("missing mongo database name")
		}
	}
	return nil
}

// validateCachestore will check the cachestore and validate basic requirements
func (a *AppConfig) validateCachestore() error {
	if a.Cachestore.Engine == cachestore.Redis {
		if a.Redis == nil {
			return errors.New("missing redis config")
		} else if len(a.Redis.URL) == 0 {
			return errors.New("missing redis url")
		}
	}
	return nil
}

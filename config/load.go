package config

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/BuxOrg/bux-server/dictionary"
	"github.com/spf13/viper"
)

// isValidEnvironment will return true if the testEnv is a known valid environment
func isValidEnvironment(testEnv string) bool {
	testEnv = strings.ToLower(testEnv)
	for _, env := range environments {
		if env == testEnv {
			return true
		}
	}
	return false
}

// getWorkingDirectory will get the current working directory
func getWorkingDirectory() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	return dir
}

// Added a mutex lock for a race-condition
var viperLock sync.Mutex

// Load all environment variables
func Load(customWorkingDirectory string) (_appConfig *AppConfig, err error) {
	// Check the environment we are running
	environment := os.Getenv(EnvironmentKey)
	if !isValidEnvironment(environment) {
		err = fmt.Errorf(dictionary.GetInternalMessage(dictionary.ErrorInvalidEnv), environment)
		return
	}

	// Get the working directory
	var workingDirectory string
	if len(customWorkingDirectory) > 0 {
		workingDirectory = customWorkingDirectory
	} else {
		workingDirectory = getWorkingDirectory()
	}

	viperLock.Lock()

	// Load configuration from json based on the environment from our working directory
	viper.SetConfigFile(workingDirectory + "/config/envs/" + environment + ".json") // For production (aws)

	// Set a replacer for replacing double underscore with nested period
	replacer := strings.NewReplacer(".", "__")
	viper.SetEnvKeyReplacer(replacer)

	// Set the prefix
	viper.SetEnvPrefix(EnvironmentPrefix)

	// Use env vars
	viper.AutomaticEnv()

	// Read the configuration
	if err = viper.ReadInConfig(); err != nil {
		err = fmt.Errorf(dictionary.GetInternalMessage(dictionary.ErrorReadingConfig), err.Error())
		return
	}

	// Initialize
	_appConfigVal := AppConfig{
		Authentication: AuthenticationConfig{},
		Cachestore:     CachestoreConfig{},
		ClusterConfig:  &ClusterConfig{},
		Datastore:      DatastoreConfig{},
		GraphQL:        GraphqlConfig{},
		Mongo:          datastore.MongoDBConfig{},
		Monitor:        MonitorOptions{},
		NewRelic:       NewRelicConfig{},
		Notifications:  NotificationsConfig{},
		Paymail:        PaymailConfig{},
		Redis:          RedisConfig{},
		Server:         ServerConfig{},
		TaskManager:    TaskManagerConfig{},
		Pulse:          PulseConfig{},
	}

	// Unmarshal into values struct
	if err = viper.Unmarshal(&_appConfigVal); err != nil {
		err = fmt.Errorf(dictionary.GetInternalMessage(dictionary.ErrorViper), err.Error())
		return
	}

	viperLock.Unlock()

	// Set working directory
	_appConfigVal.WorkingDirectory = workingDirectory
	_appConfig = &_appConfigVal

	return
}

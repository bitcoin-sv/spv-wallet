package config

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/spf13/viper"

	"github.com/BuxOrg/bux-server/dictionary"
)

// Added a mutex lock for a race-condition
var viperLock sync.Mutex

// Load all AppConfig
func Load(configFilePath string) (appConfig *AppConfig, err error) {
	setDefaults(configFilePath)

	loadFlags()

	envConfig()

	viperLock.Lock()

	if err = loadFromFile(); err != nil {
		return nil, err
	}

	appConfig = new(AppConfig)
	if err = unmarshallToAppConfig(appConfig); err != nil {
		return nil, err
	}

	viperLock.Unlock()

	return appConfig, nil
}

func envConfig() {
	viper.SetEnvPrefix("BUX")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
}

func loadFromFile() error {
	configFilePath := viper.GetString(ConfigFilePathKey)

	if configFilePath == "" {
		_, err := os.Stat(DefaultConfigFilePath)
		if os.IsNotExist(err) {
			// if the config is not specified and no default config file exists, use defaults
			logger.Data(2, logger.DEBUG, "Config file not specified. Using defaults")
			return nil
		}
		configFilePath = DefaultConfigFilePath
	}

	viper.SetConfigFile(configFilePath)
	if err := viper.ReadInConfig(); err != nil {
		err = fmt.Errorf(dictionary.GetInternalMessage(dictionary.ErrorReadingConfig), err.Error())
		logger.Data(2, logger.ERROR, err.Error())
		return err
	}

	return nil
}

func unmarshallToAppConfig(appConfig *AppConfig) error {
	if err := viper.Unmarshal(&appConfig); err != nil {
		err = fmt.Errorf(dictionary.GetInternalMessage(dictionary.ErrorViper), err.Error())
		return err
	}
	return nil
}

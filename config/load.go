package config

import (
	"fmt"
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
	viper.SetConfigFile(configFilePath)

	if configFilePath != "" {
		viper.SetConfigFile(configFilePath)
		if err := viper.ReadInConfig(); err != nil {
			switch err.(type) {
			case *viper.ConfigFileNotFoundError:
				logger.Data(2, logger.INFO, "Config file not specified. Using defaults")
			default:
				err = fmt.Errorf(dictionary.GetInternalMessage(dictionary.ErrorReadingConfig), err.Error())
				logger.Data(2, logger.ERROR, err.Error())
				return err
			}
		}
	}
	return nil
}

func unmarshallToAppConfig(appConfig *AppConfig) error {
	appConfig = new(AppConfig) // TODO: does it need to be here?

	if err := viper.Unmarshal(&appConfig); err != nil {
		err = fmt.Errorf(dictionary.GetInternalMessage(dictionary.ErrorViper), err.Error())
		return err
	}
	return nil
}

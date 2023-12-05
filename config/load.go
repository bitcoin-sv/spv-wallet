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

// Load all environment variables
func Load(customWorkingDirectory string) (appConfig *AppConfig, err error) {
	setDefaults()

	// set flags

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
	configFilePath := viper.GetString("configFilePath") // TODO: extract to CONST
	viper.SetConfigFile(configFilePath)

	if configFilePath != "" {
		viper.SetConfigFile(configFilePath)
		if err := viper.ReadInConfig(); err != nil {
			err = fmt.Errorf(dictionary.GetInternalMessage(dictionary.ErrorReadingConfig), err.Error())
			return err
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

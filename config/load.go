package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/bitcoin-sv/spv-wallet/dictionary"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

// Added a mutex lock for a race-condition
var viperLock sync.Mutex

// Load all AppConfig
func Load(logger *zerolog.Logger) (appConfig *AppConfig, err error) {
	viperLock.Lock()
	defer viperLock.Unlock()

	if err = setDefaults(); err != nil {
		return nil, err
	}

	envConfig()

	if err = loadFlags(); err != nil {
		return nil, err
	}

	if err = loadFromFile(logger); err != nil {
		return nil, err
	}

	appConfig = getDefaultAppConfig()
	if err = unmarshallToAppConfig(appConfig); err != nil {
		return nil, err
	}

	if appConfig.Debug {
		cfg, err := json.MarshalIndent(appConfig, "", "  ")
		if err != nil {
			logger.Error().Msg("Unable to decode App Config to json")
		} else {
			fmt.Printf("loaded config: %s", cfg)
		}
	}

	return appConfig, nil
}

func setDefaults() error {
	viper.SetDefault(ConfigFilePathKey, DefaultConfigFilePath)

	defaultsMap := make(map[string]interface{})
	if err := mapstructure.Decode(getDefaultAppConfig(), &defaultsMap); err != nil {
		return err
	}

	for key, value := range defaultsMap {
		viper.SetDefault(key, value)
	}

	return nil
}

func envConfig() {
	viper.SetEnvPrefix("SPV")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
}

func loadFromFile(logger *zerolog.Logger) error {
	configFilePath := viper.GetString(ConfigFilePathKey)

	if configFilePath == DefaultConfigFilePath {
		_, err := os.Stat(DefaultConfigFilePath)
		if os.IsNotExist(err) {
			logger.Debug().Msg("Config file not specified. Using defaults")
			return nil
		}
		configFilePath = DefaultConfigFilePath
	}

	viper.SetConfigFile(configFilePath)
	if err := viper.ReadInConfig(); err != nil {
		err = fmt.Errorf(dictionary.GetInternalMessage(dictionary.ErrorReadingConfig), err.Error())
		logger.Error().Msg(err.Error())
		return err
	}

	return nil
}

func unmarshallToAppConfig(appConfig *AppConfig) error {
	if err := viper.Unmarshal(appConfig); err != nil {
		err = fmt.Errorf(dictionary.GetInternalMessage(dictionary.ErrorViper), err.Error())
		return err
	}
	return nil
}

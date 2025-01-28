package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/bitcoin-sv/spv-wallet/dictionary"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/go-viper/mapstructure/v2"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

// Added a mutex lock for a race-condition
var viperLock sync.Mutex

// Load all AppConfig
func Load(versionToSet string, logger zerolog.Logger) (appConfig *AppConfig, err error) {
	viperLock.Lock()
	defer viperLock.Unlock()

	appConfig = GetDefaultAppConfig()
	appConfig.Version = versionToSet
	if err = setDefaults(appConfig); err != nil {
		return nil, err
	}

	envConfig()

	if err = loadFlags(appConfig); err != nil {
		return nil, err
	}

	if err = loadFromFile(logger); err != nil {
		return nil, err
	}

	if err = unmarshallToAppConfig(appConfig); err != nil {
		return nil, err
	}

	logger.Debug().MsgFunc(func() string {
		cfg, err := json.MarshalIndent(appConfig, "", "  ")
		if err != nil {
			return "Unable to decode App Config to json"
		}
		return fmt.Sprintf("loaded config: %s", cfg)
	})

	return appConfig, nil
}

func setDefaults(config *AppConfig) error {
	viper.SetDefault(ConfigFilePathKey, DefaultConfigFilePath)

	defaultsMap := make(map[string]interface{})
	if err := mapstructure.Decode(config, &defaultsMap); err != nil {
		err = spverrors.Wrapf(err, "error occurred while setting defaults")
		return err
	}

	for key, value := range defaultsMap {
		viper.SetDefault(key, value)
	}

	return nil
}

func envConfig() {
	viper.SetEnvPrefix(envPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
}

func loadFromFile(logger zerolog.Logger) error {
	configFilePath := viper.GetString(ConfigFilePathKey)

	if configFilePath == DefaultConfigFilePath {
		_, err := os.Stat(configFilePath)
		if os.IsNotExist(err) {
			logger.Debug().Msg("Config file not specified. Using defaults")
			return nil
		}
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
		err = spverrors.Wrapf(err, "error when unmarshalling config to App Config")
		return err
	}
	return nil
}

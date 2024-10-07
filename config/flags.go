package config

import (
	"fmt"
	"os"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type cliFlags struct {
	showVersion bool `mapstructure:"version"`
	showHelp    bool `mapstructure:"help"`
	dumpConfig  bool `mapstructure:"dump_config"`
}

func loadFlags(cfg *AppConfig) error {
	if !anyFlagsPassed() {
		return nil
	}

	cli := &cliFlags{}
	appFlags := pflag.NewFlagSet("appFlags", pflag.ContinueOnError)

	initFlags(appFlags, cli)

	err := appFlags.Parse(os.Args[1:])
	if err != nil {
		fmt.Printf("Flags can't be parsed: %v\n", err)
		os.Exit(1)
	}

	err = viper.BindPFlag(ConfigFilePathKey, appFlags.Lookup(ConfigFilePathKey))
	if err != nil {
		err = spverrors.Wrapf(err, "error while binding flags to viper")
		return err
	}

	parseCliFlags(cfg, appFlags, cli)

	return nil
}

func anyFlagsPassed() bool {
	return len(os.Args) > 1
}

func initFlags(fs *pflag.FlagSet, cliFlags *cliFlags) {
	fs.StringP(ConfigFilePathKey, "C", "", "custom config file path")

	fs.BoolVarP(&cliFlags.showHelp, "help", "h", false, "show help")
	fs.BoolVarP(&cliFlags.showVersion, "version", "v", false, "show version")
	fs.BoolVarP(&cliFlags.dumpConfig, "dump_config", "d", false, "dump config to file, specified by config_file flag")
}

func parseCliFlags(cfg *AppConfig, fs *pflag.FlagSet, cli *cliFlags) {
	if cli.showHelp {
		fs.PrintDefaults()
		os.Exit(0)
	}

	if cli.showVersion {
		fmt.Println(cfg.GetUserAgent())
		os.Exit(0)
	}

	if cli.dumpConfig {
		configPath := viper.GetString(ConfigFilePathKey)
		if configPath == "" {
			configPath = DefaultConfigFilePath
		}

		err := viper.SafeWriteConfigAs(configPath)
		if err != nil {
			fmt.Printf("error while dumping config: %v", err.Error())
		}
		os.Exit(0)
	}
}

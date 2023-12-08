package cli

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/BuxOrg/bux-server/cli/flags"
)

func ParseCliFlags(cli *flags.CliFlags, version, configFilePathKey string) {
	if cli.ShowHelp {
		pflag.Usage()
		os.Exit(0)
	}

	if cli.ShowVersion {
		fmt.Println("bux-sever", "version", version)
		os.Exit(0)
	}

	if cli.DumpConfig {
		viper.SafeWriteConfigAs(viper.GetString(configFilePathKey))
		os.Exit(0)
	}
}

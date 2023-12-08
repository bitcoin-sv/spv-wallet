package flags

type CliFlags struct {
	ShowVersion bool `mapstructure:"version"`
	ShowHelp    bool `mapstructure:"help"`
	DumpConfig  bool `mapstructure:"dump_config"`
}

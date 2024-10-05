package main

import (
	"flag"
	"github.com/cristalhq/aconfig"
	"github.com/cristalhq/aconfig/aconfigyaml"
)

type Config struct {
	Currency CurrencyConfig
}

type CurrencyConfig struct {
	AppID string `required:"true"`
}

type requiredFlags struct {
	ConfigFile string
}

func (f *requiredFlags) FlagSet() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.StringVar(&f.ConfigFile, "config", "config.yaml", "path to config file")
	return fs
}

func parseConfig(args []string) (*Config, error) {
	var flags requiredFlags
	if err := flags.FlagSet().Parse(args); err != nil {
		return nil, err
	}

	var cfg Config
	loader := aconfig.LoaderFor(&cfg, aconfig.Config{
		SkipFlags: true,
		Files:     []string{flags.ConfigFile},
		FileDecoders: map[string]aconfig.FileDecoder{
			".yaml": aconfigyaml.New(),
		},
	})

	return &cfg, loader.Load()
}

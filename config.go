package main

import (
	"os"

	"github.com/samber/lo"
	"github.com/spf13/viper"
)

type Config struct {
	LogLevel    string `mapstructure:"log_level"`
	Listen      string `mapstructure:"listen"`
	Watchexec   string `mapstructure:"watchexec"`
	ServicesDir string `mapstructure:"services_dir"`
	APIPrefix   string `mapstructure:"api_prefix"`
}

var (
	cfg Config
)

func loadConfig() {
	viper.SetConfigFile(lo.CoalesceOrEmpty(os.Getenv("CFG_FILE"), "config.toml"))
	viper.SetConfigType(lo.CoalesceOrEmpty(os.Getenv("CFG_FORMAT"), "toml"))
	lo.Must0(viper.ReadInConfig())
	lo.Must0(viper.Unmarshal(&cfg))
}

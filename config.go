package main

import (
	"log/slog"
	"os"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/samber/lo"
	"github.com/spf13/viper"
)

type Config struct {
	LogLevel    string `mapstructure:"log_level"`
	Listen      string `mapstructure:"listen"`
	Watchexec   string `mapstructure:"watchexec"`
	ServicesDir string `mapstructure:"services_dir"`
	APIPrefix   string `mapstructure:"api_prefix"`
	PostgresDSN string `mapstructure:"postgres_dsn"`

	Vars     map[string]string `mapstructure:"vars"`
	VarsDesc map[string]string `mapstructure:"vars_desc"`
}

var (
	cfgLock sync.RWMutex
	cfg     Config
)

func loadConfig() {
	viper.SetConfigFile(lo.CoalesceOrEmpty(os.Getenv("CFG_FILE"), "config.toml"))
	viper.SetConfigType(lo.CoalesceOrEmpty(os.Getenv("CFG_FORMAT"), "toml"))

	reloadConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		reloadConfig()
	})

	viper.WatchConfig()
}

func reloadConfig() {
	cfgLock.Lock()
	defer cfgLock.Unlock()

	lo.Must0(viper.ReadInConfig())
	lo.Must0(viper.Unmarshal(&cfg))
	setLogLevel()
}

func setLogLevel() {
	var l slog.Level
	l.UnmarshalText([]byte(cfg.LogLevel))
	slog.SetLogLoggerLevel(l)
}

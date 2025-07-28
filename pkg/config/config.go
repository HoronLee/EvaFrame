package config

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/google/wire"
	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port int    `mapstructure:"port"`
		Mode string `mapstructure:"mode"`
	} `mapstructure:"server"`

	Database struct {
		Type string `mapstructure:"type"` // 新增数据库类型字段
		DSN  string `mapstructure:"dsn"`
	} `mapstructure:"database"`

	JWT struct {
		Secret string `mapstructure:"secret"`
	} `mapstructure:"jwt"`

	Logger struct {
		Level   string `mapstructure:"level"`
		LogPath string `mapstructure:"log_path"`
	} `mapstructure:"logger"`

	DevChoice struct {
		DAO string `mapstructure:"dao"`
	} `mapstructure:"dev_choice"`
}

func NewConfig(path string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// 设置热更新
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})

	var c Config
	if err := v.Unmarshal(&c); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &c, nil
}

var ProviderSet = wire.NewSet(NewConfig)

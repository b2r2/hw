package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

var (
	ErrHostRequired = errors.New("http app host is required")
	ErrPortRequired = errors.New("http app port is required")
	ErrStorageName  = errors.New("storage name is required")
	ErrStorageHost  = errors.New("storage host is required")
	ErrStoragePort  = errors.New("storage port is required")
)

type Config struct {
	Logger  LoggerConf
	Server  ServerConf
	Storage StorageConf
}

type LoggerConf struct {
	Level string `toml:"level"`
	Path  string `toml:"filepath"`
}

type ServerConf struct {
	Host string `toml:"host"`
	Port string `toml:"port"`
}

type StorageConf struct {
	IsMem    bool   `toml:"is_mem"`
	Host     string `toml:"host"`
	Port     string `toml:"port"`
	Database string `toml:"database"`
	SSL      string `toml:"ssl"`
}

func NewConfig(filepath string) (*Config, error) {
	var config Config

	v := viper.New()

	if filepath == "" {
		v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
		v.AutomaticEnv()
		v.SetDefault("logger.level", "DEBUG")
		v.SetDefault("http.host", "localhost")
		v.SetDefault("http.port", "8080")
		v.SetDefault("database.ismem", true)
	}

	v.SetConfigFile(filepath)
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read configuration: %w", err)
	}

	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to read configuration: %w", err)
	}

	if err := config.validate(); err != nil {
		return nil, fmt.Errorf("failed to read configuration: %w", err)
	}

	return &config, nil
}

func (c *Config) validate() error {
	if c.Server.Host == "" {
		return ErrHostRequired
	}

	if c.Server.Port == "" {
		return ErrPortRequired
	}

	if !c.Storage.IsMem && c.Storage.Database == "" {
		return ErrStorageName
	}

	if !c.Storage.IsMem && c.Storage.Host == "" {
		return ErrStorageHost
	}

	if !c.Storage.IsMem && c.Storage.Port == "" {
		return ErrStoragePort
	}

	return nil
}

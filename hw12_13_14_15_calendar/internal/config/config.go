package config

import (
	"errors"
	"fmt"

	"github.com/BurntSushi/toml"
)

var (
	ErrAddrHTTPRequired = errors.New("http address is required")
	ErrAddrGRPCRequired = errors.New("grpc address is required")
	ErrStorageName      = errors.New("storage name is required")
	ErrStorageHost      = errors.New("storage host is required")
	ErrStoragePort      = errors.New("storage port is required")
)

type Config struct {
	Logger  LoggerConf  `toml:"logger"`
	Server  ServerConf  `toml:"server"`
	Storage StorageConf `toml:"storage"`
}

type LoggerConf struct {
	Level string `toml:"level"`
	Path  string `toml:"file_path"`
}

type ServerConf struct {
	AddrHTTP string `toml:"addr_http"`
	AddrGRPC string `toml:"addr_grpc"`
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

	if _, err := toml.DecodeFile(filepath, &config); err != nil {
		return nil, err
	}

	if err := config.validate(); err != nil {
		return nil, fmt.Errorf("failed to read configuration: %w", err)
	}

	return &config, nil
}

func (c *Config) validate() error {
	if c.Server.AddrHTTP == "" {
		return ErrAddrHTTPRequired
	}
	if c.Server.AddrGRPC == "" {
		return ErrAddrGRPCRequired
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

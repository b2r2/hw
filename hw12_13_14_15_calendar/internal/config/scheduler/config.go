package config

import (
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"

	"github.com/b2r2/hw/hw12_13_14_15_calendar/internal/config"
)

type duration struct {
	time.Duration
}

type Config struct {
	Scheduler SchedulerConf      `toml:"scheduler"`
	RabbitMQ  RabbitMQConf       `toml:"rabbit"`
	Storage   config.StorageConf `toml:"storage"`
	Logger    config.LoggerConf  `toml:"logger"`
}

type SchedulerConf struct {
	Duration duration `json:"duration"`
}

type RabbitMQConf struct {
	TTL int    `json:"ttl"`
	DSN string `json:"dsn"`
}

var (
	ErrRMQ = errors.New("DSN is empty")
	ErrTTL = errors.New("TTL must be bigger then 0")
)

func New(path string) (*Config, error) {
	if path == "" {
		path = filepath.Join("configs", "scheduler.toml")
	}
	var c Config
	if _, err := toml.DecodeFile(path, &c); err != nil {
		return nil, err
	}

	if err := c.validate(); err != nil {
		return nil, fmt.Errorf("failed to read configuration: %w", err)
	}
	return &c, nil
}

func (c *Config) validate() error {
	if c.RabbitMQ.DSN == "" {
		return ErrRMQ
	}
	if c.RabbitMQ.TTL <= 0 {
		return ErrTTL
	}
	return nil
}

func (d *duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}

package config

import (
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/b2r2/hw/hw12_13_14_15_calendar/internal/config"
	scheduler "github.com/b2r2/hw/hw12_13_14_15_calendar/internal/config/scheduler"
)

type SenderConf struct {
	Logger   config.LoggerConf
	RabbitMQ scheduler.RabbitMQConf
}

func New(path string) (*SenderConf, error) {
	if path == "" {
		path = filepath.Join("configs", "sender_config.toml")
	}
	var c SenderConf
	if _, err := toml.DecodeFile(path, &c); err != nil {
		return nil, err
	}

	return &c, nil
}

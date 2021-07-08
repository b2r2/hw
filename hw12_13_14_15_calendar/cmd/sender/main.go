package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/b2r2/hw/hw12_13_14_15_calendar/internal/storage"

	"github.com/b2r2/hw/hw12_13_14_15_calendar/internal/rmq"

	config "github.com/b2r2/hw/hw12_13_14_15_calendar/internal/config/sender"
	"github.com/b2r2/hw/hw12_13_14_15_calendar/internal/logger"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "", "Path to configuration file")
}

func main() {
	flag.Parse()

	mainContext, cancel := context.WithCancel(context.Background())

	go func(cancel context.CancelFunc) {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
		<-signals
		cancel()
	}(cancel)

	c, err := config.New(configFile)
	if err != nil {
		log.Fatal(err)
	}

	logg, err := logger.New(c.Logger.Level, c.Logger.Path, nil)
	if err != nil {
		log.Fatal(err)
	}

	rabbit, err := rmq.New(logg, c.RabbitMQ.DSN, c.RabbitMQ.TTL)
	if err != nil {
		log.Fatal(err)
	}

	if err := rabbit.Send(mainContext, prepare(logg)); err != nil {
		log.Fatal(err)
	}

	logg.Info("sender running")

	<-mainContext.Done()

	logg.Info("stopping sender")

	if err := rabbit.Close(); err != nil {
		log.Fatal(err)
	}
}

func prepare(log logger.Logger) func([]byte) {
	return func(body []byte) {
		n := storage.Notification{}
		if err := json.Unmarshal(body, &n); err != nil {
			log.Errorln("failed to decode a message", string(body), "err", err)
		}
		log.Infoln(n.String())
	}
}

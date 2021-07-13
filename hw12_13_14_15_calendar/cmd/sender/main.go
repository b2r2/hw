package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

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

	consumer, err := rmq.NewConsumer(logg, c.RabbitMQ.DSN, "events")
	if err != nil {
		log.Fatal(err)
	}

	messages, err := consumer.Consumer(mainContext, "events")
	if err != nil {
		log.Fatal(err)
	}

	logg.Info("sender running")

	for msg := range messages {
		fmt.Println("receive new message:", string(msg.Data))
	}

	logg.Info("sender stopped")
}

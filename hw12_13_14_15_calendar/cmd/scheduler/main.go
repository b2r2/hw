package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/b2r2/hw/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/b2r2/hw/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/b2r2/hw/hw12_13_14_15_calendar/internal/storage/sql"

	"github.com/b2r2/hw/hw12_13_14_15_calendar/internal/rmq"

	"github.com/b2r2/hw/hw12_13_14_15_calendar/internal/logger"

	config "github.com/b2r2/hw/hw12_13_14_15_calendar/internal/config/scheduler"
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

	scheduler := time.NewTicker(c.Scheduler.Duration.Duration)
	defer scheduler.Stop()

	var db storage.Storage
	if c.Storage.IsMem {
		db = memorystorage.New(logg)
	} else {
		dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			c.Storage.Host,
			c.Storage.Port,
			"calendar",
			"calendar",
			c.Storage.Database,
			c.Storage.SSL)

		if db, err = sqlstorage.New(logg, mainContext, dsn); err != nil {
			log.Fatal(err)
		}
	}

	go func() {
		for {
			select {
			case <-mainContext.Done():
				return
			case <-scheduler.C:
				events, err := db.ListNotifyEvents(mainContext)
				if err != nil {
					logg.Errorln("failed to retrieve notification events:", err)
				}
				rabbit.Notify(events)
			}
		}
	}()

	logg.Info("scheduler running")

	<-mainContext.Done()

	logg.Info("stopping scheduler")

	ctx, stop := context.WithTimeout(context.Background(), time.Second*5)
	if err := db.Close(ctx); err != nil {
		stop()
		log.Fatal(err)
	}
	if err := rabbit.Close(); err != nil {
		stop()
		log.Fatal(err)
	}
}

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/b2r2/hw/hw12_13_14_15_calendar/internal/app"
	"github.com/b2r2/hw/hw12_13_14_15_calendar/internal/config"
	"github.com/b2r2/hw/hw12_13_14_15_calendar/internal/logger"
	httpserver "github.com/b2r2/hw/hw12_13_14_15_calendar/internal/server/http"
	"github.com/b2r2/hw/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/b2r2/hw/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/b2r2/hw/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./configs/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		os.Exit(0)
	}

	mainContext, cancel := context.WithCancel(context.Background())

	go func(cancel context.CancelFunc) {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
		<-signals
		cancel()
	}(cancel)

	conf, err := config.NewConfig(configFile)
	if err != nil {
		log.Fatal(err)
	}
	logg, err := logger.New(conf.Logger.Level, conf.Logger.Path, nil)
	if err != nil {
		log.Fatal(err)
	}
	logg.Info("calendar app started")

	var db storage.Storage
	if conf.Storage.IsMem {
		db = memorystorage.New(logg)
	} else {
		db = sqlstorage.New(logg)
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		conf.Storage.Host,
		conf.Storage.Port,
		"calendar",
		"calendar",
		conf.Storage.Database,
		conf.Storage.SSL)

	if err := db.Connect(mainContext, dsn); err != nil {
		log.Fatal(err)
	}

	calendar := app.New(logg, db)
	server := httpserver.NewServer(calendar, logg)

	go func(cancel context.CancelFunc) {
		if err := server.Start(net.JoinHostPort(conf.Server.Host, conf.Server.Port)); err != nil {
			cancel()
			log.Fatal(err)
		}
	}(cancel)

	logg.Info("calendar app is running")

	<-mainContext.Done()

	logg.Info("stopping calendar app")
	cancel()

	ctx, newCancel := context.WithTimeout(context.Background(), time.Second*5)

	if err := db.Close(ctx); err != nil {
		newCancel()
		log.Fatal(err)
	}

	if err := server.Stop(ctx); err != nil {
		newCancel()
		log.Fatal(err)
	}
}

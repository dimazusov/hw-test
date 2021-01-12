package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/dimazusov/hw-test/hw12_13_14_15_calendar/internal/app"
	"github.com/dimazusov/hw-test/hw12_13_14_15_calendar/internal/config"
	"github.com/dimazusov/hw-test/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/dimazusov/hw-test/hw12_13_14_15_calendar/internal/server/http"
	"github.com/dimazusov/hw-test/hw12_13_14_15_calendar/internal/storage"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", `configs/config.yaml`, "Path to configuration file")
}

func main() {
	config, err := config.New(configFile)
	if err != nil {
		log.Fatalln(err)
	}

	logger, err := logger.New(config.Logger.Path, config.Logger.Level)
	if err != nil {
		log.Fatalln(err)
	}

	rep, err := storage.NewRepository(config)
	if err != nil {
		log.Fatalln(err)
	}

	calendar := app.New(logger, rep.(app.Repository))
	server := internalhttp.NewServer(config, calendar)

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals)

		<-signals
		signal.Stop(signals)

		if err := server.Stop(); err != nil {
			err = logger.Error("failed to stop http server: " + err.Error())
			if err != nil {
				log.Println(err)
			}
		}
	}()

	if err = server.Start(); err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
}

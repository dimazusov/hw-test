package main

import (
	"flag"
	"github.com/dimazusov/hw-test/hw12_13_14_15_calendar/internal/app"
	"github.com/dimazusov/hw-test/hw12_13_14_15_calendar/internal/config"
	"github.com/dimazusov/hw-test/hw12_13_14_15_calendar/internal/logger"
	"github.com/dimazusov/hw-test/hw12_13_14_15_calendar/internal/server"
	"github.com/dimazusov/hw-test/hw12_13_14_15_calendar/internal/storage"
	"log"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", `confisdgs/config.yaml`, "Path to configuration file")
	flag.Parse()
}

func main() {
	cfg, err := config.New(configFile)
	if err != nil {
		log.Fatalln(err)
	}

	lg, err := logger.New(cfg.Logger.Path, cfg.Logger.Level)
	if err != nil {
		log.Fatalln(err)
	}

	rep, err := storage.NewRepository(cfg)
	if err != nil {
		log.Fatalln(err)
	}

	calendar := app.New(lg, rep.(app.Repository))

	servers := server.NewServers(cfg, calendar)
	servers.StartServers()
}

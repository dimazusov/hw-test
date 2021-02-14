package main

import (
	"flag"
	"github.com/dimazusov/hw-test/hw12_13_14_15_calendar/internal/config"
	"github.com/dimazusov/hw-test/hw12_13_14_15_calendar/internal/rabbitmq"
	"log"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", `configs/config.yaml`, "Path to configuration file")
	flag.Parse()
}

func main() {
	cfg, err := config.New(configFile)
	if err != nil {
		log.Fatalln(err)
	}

	consumer := rabbitmq.NewConsumer(cfg.Rabbitmq.Url, rabbitmq.QueueEvents)
	err = consumer.Connect()
	if err != nil {
		log.Fatalln(err)
	}

	err = consumer.Consume(func(b []byte) error {
		log.Println("message ", string(b))

		return nil
	})
	if err != nil {
		log.Fatalln(err)
	}
}

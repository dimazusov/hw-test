package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"time"

	"github.com/dimazusov/hw-test/hw12_13_14_15_calendar/internal/app"
	"github.com/dimazusov/hw-test/hw12_13_14_15_calendar/internal/config"
	"github.com/dimazusov/hw-test/hw12_13_14_15_calendar/internal/logger"
	"github.com/dimazusov/hw-test/hw12_13_14_15_calendar/internal/rabbitmq"
	"github.com/dimazusov/hw-test/hw12_13_14_15_calendar/internal/scheduler"
	"github.com/dimazusov/hw-test/hw12_13_14_15_calendar/internal/storage"
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

	lg, err := logger.New(cfg.Logger.Path, cfg.Logger.Level)
	if err != nil {
		log.Fatalln(err)
	}

	rep, err := storage.NewRepository(cfg)
	if err != nil {
		log.Fatalln(err)
	}

	calendar := app.New(lg, rep.(app.Repository))

	sc := scheduler.New()
	sc.AddTask(func() error {
		log.Println("task run")

		params := make(map[string]interface{})

		params["notificationTimeFrom"] = uint(time.Now().Unix())
		params["isNotificationSend"] = "false"

		events, err := calendar.GetEventsByParams(context.Background(), params)
		if err != nil {
			return err
		}

		sendedIDs := make([]uint, 0, len(events))
		defer func() {
			for _, eventID := range sendedIDs {
				event, err := calendar.GetEventByID(context.Background(), eventID)
				if err != nil {
					log.Fatalln(err)
				}

				event.IsNotificationSend = true

				err = calendar.Update(context.Background(), event)
				if err != nil {
					log.Fatalln(err)
				}
			}
		}()

		mqProducer := rabbitmq.NewProducer(cfg.Rabbitmq.Url, rabbitmq.ExchangeEvents)

		err = mqProducer.Connect()
		if err != nil {
			return err
		}

		for _, event := range events {
			b, err := json.Marshal(event)
			if err != nil {
				return err
			}

			err = mqProducer.SendMessage(b)
			if err != nil {
				return err
			}

			sendedIDs = append(sendedIDs, event.ID)
		}

		err = calendar.DeleteOldEvents(context.Background(), cfg.EventTimeExpired)
		if err != nil {
			return err
		}

		return nil
	}, 5*time.Second)

	err = sc.Run()
	if err != nil {
		log.Fatalln(err)
	}
}

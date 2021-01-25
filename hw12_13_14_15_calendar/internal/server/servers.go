package server

import (
	"context"
	"fmt"
	internalgrpc "github.com/dimazusov/hw-test/hw12_13_14_15_calendar/internal/server/grpc"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/dimazusov/hw-test/hw12_13_14_15_calendar/internal/app"
	"github.com/dimazusov/hw-test/hw12_13_14_15_calendar/internal/config"
	internalhttp "github.com/dimazusov/hw-test/hw12_13_14_15_calendar/internal/server/http"
)

type Server interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

type StartServers interface {
	StartServers()
}

type servers struct {
	cfg   *config.Config
	app   app.App
	servs []Server
}

func NewServers(cfg *config.Config, app app.App) StartServers {
	return &servers{
		cfg: cfg,
		app: app,
	}
}

func (m servers) StartServers() {
	m.servs = []Server{
		internalhttp.NewServer(m.cfg, m.app),
		internalgrpc.NewServer(m.cfg, m.app),
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGHUP)

		<-signals
		signal.Stop(signals)
		cancel()

		for _, s := range m.servs {
			if err := s.Stop(context.Background()); err != nil {
				err = m.app.LogInfo("failed to stop http server: " + err.Error())
				if err != nil {
					log.Println(err)
				}
			}
		}
	}()

	wg := sync.WaitGroup{}
	wg.Add(len(m.servs))
	for i := range m.servs {
		go func(s Server, app app.App) {
			defer wg.Done()

			fmt.Println("start")
			if err := s.Start(ctx); err != nil {
				err = app.LogError("failed to start http server: " + err.Error())
				if err != nil {
					log.Println(err)
				}

				os.Exit(1)
			}
		}(m.servs[i], m.app)
	}
	wg.Wait()
}

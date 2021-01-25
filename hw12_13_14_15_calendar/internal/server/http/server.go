package internalhttp

import (
	"context"
	"net/http"

	"github.com/dimazusov/hw-test/hw12_13_14_15_calendar/internal/config"
	"github.com/dimazusov/hw-test/hw12_13_14_15_calendar/internal/domain"
	"github.com/pkg/errors"
)

type Server interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

type server struct {
	app Application
	cfg *config.Config
	srv *http.Server
}

type Application interface {
	LogInfo(data interface{}) error
	LogError(data interface{}) error
	Create(ctx context.Context, event domain.Event) (newID uint, err error)
	Update(ctx context.Context, event domain.Event) (err error)
	Delete(ctx context.Context, eventID uint) (err error)
	GetEventByID(ctx context.Context, eventID uint) (event domain.Event, err error)
	GetEventsByParams(ctx context.Context, params map[string]interface{}) (events []domain.Event, err error)
}

func NewServer(cfg *config.Config, app Application) Server {
	return &server{
		cfg: cfg,
		app: app,
	}
}

func (m *server) Start(ctx context.Context) error {
	router := NewGinRouter(m.app)

	m.srv = &http.Server{}
	m.srv.Addr = m.cfg.Server.HTTP.Host + ":" + m.cfg.Server.HTTP.Port
	m.srv.Handler = router

	err := m.srv.ListenAndServe()
	if err != nil {
		return errors.Wrap(err, "cannot listen and serve")
	}

	return nil
}

func (m *server) Stop(ctx context.Context) error {
	err := m.srv.Shutdown(ctx)
	if err != nil {
		return errors.Wrap(err, "cannot shutdown server")
	}

	return nil
}

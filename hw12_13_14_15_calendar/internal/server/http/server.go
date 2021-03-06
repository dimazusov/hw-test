package internalhttp

import (
	"context"
	"net/http"

	"github.com/dimazusov/hw-test/hw12_13_14_15_calendar/internal/config"
	"github.com/pkg/errors"
)

type Server struct {
	host string
	port string
	app  Application
	srv  *http.Server
}

type Application interface {
	LogInfo(data interface{}) error
}

func NewServer(cfg *config.Config, app Application) *Server {
	return &Server{
		host: cfg.Server.Host,
		port: cfg.Server.Port,
		app:  app,
	}
}

func (m *Server) Start(ctx context.Context) error {
	router := NewGinRouter(m.app)

	m.srv = &http.Server{}
	m.srv.Addr = m.host + ":" + m.port
	m.srv.Handler = router

	err := m.srv.ListenAndServe()
	if err != nil {
		return errors.Wrap(err, "cannot listen and serve")
	}

	return nil
}

//func (s *Server) Start(ctx context.Context) error {
//	// TODO
//	select {
//	case <-ctx.Done():
//		return nil
//	}
//}

func (m *Server) Stop(ctx context.Context) error {
	err := m.srv.Shutdown(ctx)
	if err != nil {
		return errors.Wrap(err, "cannot shutdown server")
	}

	return nil
}

package grpc

import (
	"context"
	"fmt"
	"github.com/dimazusov/hw-test/hw12_13_14_15_calendar/internal/config"
	"github.com/dimazusov/hw-test/hw12_13_14_15_calendar/internal/domain"
	"github.com/dimazusov/hw-test/hw12_13_14_15_calendar/internal/server/grpc/pb"
	"google.golang.org/grpc"
	"log"
	"net"
)

type Server struct {
	addr string
	srv  *grpc.Server
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

func NewServer(cfg *config.Config, app Application) *Server {
	grpcSrv := grpc.NewServer()
	eventsServer := newEventService(app)

	pb.RegisterEventsServer(grpcSrv, eventsServer)

	return &Server{
		srv:  grpcSrv,
		addr: fmt.Sprintf("%s:%s", cfg.Server.Grpc.Host, cfg.Server.Grpc.Port),
	}
}

func (m *Server) Start(ctx context.Context) error {
	lsn, err := net.Listen("tcp", m.addr)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("starting server on %s", lsn.Addr().String())
	if err := m.srv.Serve(lsn); err != nil {
		log.Fatal(err)
	}

	return nil
}

func (m *Server) Stop(ctx context.Context) error {
	m.srv.Stop()
	return nil
}

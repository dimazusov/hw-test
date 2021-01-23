package app

import (
	"context"

	"github.com/dimazusov/hw-test/hw12_13_14_15_calendar/internal/domain"
)

type app struct {
	logger Logger
	rep    Repository
}

type Logger interface {
	Debug(data interface{}) error
	Info(data interface{}) error
	Warn(data interface{}) error
	Error(data interface{}) error
	Close() error
}

type Repository interface {
	Create(ctx context.Context, event domain.Event) (newID uint, err error)
	Update(ctx context.Context, event domain.Event) (err error)
	Delete(ctx context.Context, eventID uint) (err error)
	GetEventByID(ctx context.Context, eventID uint) (event domain.Event, err error)
	GetEventsByParams(ctx context.Context, params map[string]interface{}) (events []domain.Event, err error)
}

type App interface {
	LogInfo(interface{}) error
}

func New(logger Logger, repository Repository) App {
	return &app{
		logger: logger,
		rep:    repository,
	}
}

func (m *app) LogInfo(data interface{}) error {
	return m.logger.Info(data)
}

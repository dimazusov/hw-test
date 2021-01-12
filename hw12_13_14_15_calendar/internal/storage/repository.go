package storage

import (
	"context"

	"github.com/dimazusov/hw-test/hw12_13_14_15_calendar/internal/config"
	"github.com/dimazusov/hw-test/hw12_13_14_15_calendar/internal/domain"
	memorystorage "github.com/dimazusov/hw-test/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/dimazusov/hw-test/hw12_13_14_15_calendar/internal/storage/postgres"
	"github.com/pkg/errors"
)

var ErrRepositoryTypeNotExists = errors.New("repository type not exists")

const RepTypeMemory = "memory"
const RepTypePostgres = "postgres"

type Repository interface {
	Create(ctx context.Context, event domain.Event) (newID uint, err error)
	Update(ctx context.Context, event domain.Event) (err error)
	Delete(ctx context.Context, eventID uint) (err error)
	GetEventByID(ctx context.Context, eventID uint) (event domain.Event, err error)
	GetEventsByParams(ctx context.Context, params map[string]interface{}) (events []domain.Event, err error)
}

func NewRepository(cfg *config.Config) (Repository, error) {
	var storage interface{}
	var err error

	switch cfg.Repository.Type {
	case RepTypeMemory:
		storage, err = memorystorage.New(cfg.DB.Memory.MaxSize)
	case RepTypePostgres:
		storage, err = sqlstorage.New(cfg.DB.Postgres.Dialect, cfg.DB.Postgres.Dsn)
	default:
		err = ErrRepositoryTypeNotExists
	}

	if err != nil {
		return nil, errors.Wrap(err, "cannot create repository")
	}

	return storage.(Repository), nil
}

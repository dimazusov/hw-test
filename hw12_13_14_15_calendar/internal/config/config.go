package config

import (
	"context"

	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/file"
	"github.com/pkg/errors"
)

type Config struct {
	Server struct {
		HTTP struct {
			Host string `config:"host"`
			Port string `config:"host"`
		} `config:"http"`
		Grpc struct {
			Host string `config:"host"`
			Port string `config:"host"`
		} `config:"grpc"`
	} `config:"server"`
	Logger struct {
		Path  string `config:"path"`
		Level string `config:"level"`
	} `config:"logger"`
	DB struct {
		Postgres struct {
			Dialect string `config:"dialect"`
			Dsn     string `config:"dsn"`
		} `config:"postgres"`
		Memory struct {
			MaxSize uint `config:"maxsize"`
		} `config:"memory"`
	} `config:"db"`
	Repository struct {
		Type string `config:"type"`
	} `config:"repository"`
	Rabbitmq struct {
		Url string `config:"url"`
	} `config:"rabbitmq"`
	EventTimeExpired uint `config:"eventTimeExpired"`
}

func New(filePath string) (*Config, error) {
	loader := confita.NewLoader(
		file.NewBackend(filePath),
	)

	cfg := &Config{}
	err := loader.Load(context.Background(), cfg)
	if err != nil {
		return nil, errors.Wrap(err, "cannot load config")
	}

	return cfg, nil
}

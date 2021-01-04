package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewConfig(t *testing.T) {
	cfg, err := New("config_test.yaml")
	require.Nil(t, err)
	require.NotNil(t, cfg)

	require.Equal(t, "localhost", cfg.Server.Host)
	require.Equal(t, "81", cfg.Server.Port)

	require.Equal(t, "debug", cfg.Logger.Level)
	require.Equal(t, "log/log.txt", cfg.Logger.Path)

	require.Equal(t, "postgres", cfg.Db.Postgres.Dialect)
	require.Equal(t, "host=localhost port=5401 dbname=postgres user=postgres password=postgres sslmode=disable", cfg.Db.Postgres.Dsn)

	require.Equal(t, 20, cfg.Db.Memory.MaxCount)

	require.Equal(t, "memory", cfg.Repository.Type)
}

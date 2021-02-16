//nolint:errcheck
package integration_tests

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func createDB() (testcontainers.Container, string, error) {
	var env = map[string]string{
		"POSTGRES_USER":     "db",
		"POSTGRES_PASSWORD": "db",
		"POSTGRES_DB":       "db",
	}

	var port = "5432/tcp"
	dbURL := func(port nat.Port) string {
		return fmt.Sprintf("postgres://db:@localhost:%s/%s?sslmode=disable", port.Port(), env["POSTGRES_DB"])
	}
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "postgres",
		ExposedPorts: []string{port},
		Env:          env,
		WaitingFor:   wait.ForSQL(nat.Port(port), "postgres", dbURL).Timeout(time.Second * 3),
	}
	postgresC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          false,
	})
	if err != nil {
		return nil, "", err
	}

	postgresC.Start(context.Background())

	ip, err := postgresC.Host(ctx)
	if err != nil {
		return nil, "", err
	}
	p, err := postgresC.MappedPort(ctx, nat.Port(port))
	if err != nil {
		return nil, "", err
	}

	dsn := "host=" + ip + " port=" + p.Port() + " dbname=db user=db password=db sslmode=disable"

	return postgresC, dsn, nil
}

func getConnDb(dsn string) (*sqlx.DB, error) {
	return sqlx.Connect("postgres", dsn)
}

//nolint:errcheck
package integration_tests

import (
	"context"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/jmoiron/sqlx"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"time"
)

func InitDB() error {
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
		WaitingFor:   wait.ForSQL(nat.Port(port), "postgres", dbURL).Timeout(time.Second * 15),
	}
	postgresC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          false,
	})
	if err != nil {
		return err
	}

	postgresC.Start(context.Background())
	defer postgresC.Terminate(ctx)

	ip, err := postgresC.Host(ctx)
	if err != nil {
		return err
	}
	p, err := postgresC.MappedPort(ctx, nat.Port(port))
	if err != nil {
		return err
	}

	fmt.Println("ip", ip)
	fmt.Println("port", p.Port())

	conn, err := sqlx.Connect("postgres", "host="+ip+" port="+p.Port()+" dbname=db user=db password=db sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}

	var a int
	err = conn.QueryRow("select 1").Scan(&a)
	if err != nil {
		log.Fatalln(err)
	}
	conn.Close()

	return nil
}

package services_test

import (
	"context"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"os"
	"testing"
)

var endpoint string

func setup() func() {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "redis:latest",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor: wait.ForAll(
			wait.ForLog("Ready to accept connections"),
			wait.ForExposedPort(),
		),
	}

	redisC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatal(err)
	}

	endpoint, err = redisC.Endpoint(ctx, "")
	if err != nil {
		log.Fatal(err)
	}

	return func() {
		if err := redisC.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate container: %s", err.Error())
		}
	}
}

func TestMain(m *testing.M) {
	teardown := setup()
	exitVal := m.Run()
	teardown()

	os.Exit(exitVal)
}

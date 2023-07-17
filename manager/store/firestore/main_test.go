//SPDX-License-Identifier: Apache-2.0

//go:build integration

package firestore_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var endpoint string

func setup() func() {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image: "google/cloud-sdk",
		Cmd: []string{
			"gcloud",
			"emulators",
			"firestore",
			"start",
			"--host-port=0.0.0.0:8080",
		},
		ExposedPorts: []string{"8080/tcp"},
		WaitingFor: wait.ForAll(
			wait.ForLog("Dev App Server is now running"),
			wait.ForExposedPort(),
		),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatal(err)
	}

	endpoint, err = container.Endpoint(ctx, "")
	if err != nil {
		log.Fatal(err)
	}

	err = os.Setenv("FIRESTORE_EMULATOR_HOST", endpoint)
	if err != nil {
		log.Fatal(err)
	}

	return func() {
		if err := container.Terminate(ctx); err != nil {
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

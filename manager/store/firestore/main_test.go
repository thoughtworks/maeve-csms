//SPDX-License-Identifier: Apache-2.0

//go:build integration

package firestore_test

import (
	firestoreapi "cloud.google.com/go/firestore"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/iterator"
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

func cleanupAllCollections(t *testing.T, gcloudProject string) {
	cleanupCollection(t, gcloudProject, "Certificate")
	cleanupCollection(t, gcloudProject, "ChargeStation")
	cleanupCollection(t, gcloudProject, "ChargeStationSettings")
	cleanupCollection(t, gcloudProject, "ChargeStationInstallCertificates")
	cleanupCollection(t, gcloudProject, "ChargeStationRuntimeDetails")
	cleanupCollection(t, gcloudProject, "Location")
	cleanupCollection(t, gcloudProject, "OcpiParty")
	cleanupCollection(t, gcloudProject, "OcpiRegistration")
	cleanupCollection(t, gcloudProject, "Token")
	cleanupCollection(t, gcloudProject, "Transaction")
}

func cleanupCollection(t *testing.T, gcloudProject, collection string) {
	ctx := context.Background()

	client, err := firestoreapi.NewClient(ctx, gcloudProject)
	assert.NoError(t, err)

	col := client.Collection(collection)
	bulkwriter := client.BulkWriter(ctx)

	numDeleted := 0
	iter := col.Documents(ctx)
	for {
		doc, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		assert.NoError(t, err)
		_, err = bulkwriter.Delete(doc.Ref)
		assert.NoError(t, err)
		numDeleted++
	}

	if numDeleted == 0 {
		bulkwriter.End()
	}

	bulkwriter.Flush()
}

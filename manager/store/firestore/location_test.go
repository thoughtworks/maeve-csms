// SPDX-License-Identifier: Apache-2.0

//go:build integration

package firestore_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thoughtworks/maeve-csms/manager/store"
	"github.com/thoughtworks/maeve-csms/manager/store/firestore"
	"golang.org/x/net/context"
	"math/rand"
	"testing"
)

func TestSetAndLookupLocation(t *testing.T) {
	ctx := context.Background()
	locationStore, err := firestore.NewStore(ctx, "myproject")
	require.NoError(t, err)

	want := &store.Location{
		Address: "F.Rooseveltlaan 3A",
		City:    "Gent",
		Coordinates: store.GeoLocation{
			Latitude:  "51.047599",
			Longitude: "3.729944",
		},
		Country:     "BEL",
		Id:          "loc001",
		Name:        strPointer("Gent Zuid"),
		ParkingType: strPointer("ON_STREET"),
		PostalCode:  strPointer("9000"),
	}
	err = locationStore.SetLocation(ctx, want)
	require.NoError(t, err)

	got, err := locationStore.LookupLocation(ctx, "loc001")
	require.NoError(t, err)

	assert.Regexp(t, `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z`, got.LastUpdated)
	got.LastUpdated = ""

	assert.Equal(t, want, got)
}

func TestListLocations(t *testing.T) {
	ctx := context.Background()
	locationStore, err := firestore.NewStore(ctx, "myproject")
	require.NoError(t, err)

	locations := make([]*store.Location, 20)
	for i := 0; i < 20; i++ {
		locations[i] = &store.Location{
			Address: "Randomlaan 3A",
			City:    "Randomtown",
			Coordinates: store.GeoLocation{
				Latitude:  fmt.Sprintf("%f", rand.Float32()*90),
				Longitude: fmt.Sprintf("%f", rand.Float32()*180),
			},
			Country:     "RAND",
			Id:          fmt.Sprintf("loc%03d", i),
			Name:        strPointer("Random Location"),
			ParkingType: strPointer("ON_STREET"),
			PostalCode:  strPointer("12345"),
		}
	}

	for _, loc := range locations {
		err = locationStore.SetLocation(ctx, loc)
		require.NoError(t, err)
	}

	got, err := locationStore.ListLocations(ctx, 0, 10)
	require.NoError(t, err)

	assert.Equal(t, 10, len(got))
	assert.Equal(t, locations[:10], got)
}

func strPointer(s string) *string {
	return &s
}

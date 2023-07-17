// SPDX-License-Identifier: Apache-2.0

//go:build integration

package firestore_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thoughtworks/maeve-csms/manager/store"
	"github.com/thoughtworks/maeve-csms/manager/store/firestore"
)

func TestSetAndLookupChargeStationAuth(t *testing.T) {
	ctx := context.Background()

	authStore, err := firestore.NewStore(ctx, "myproject")
	require.NoError(t, err)

	want := &store.ChargeStationAuth{
		SecurityProfile:      store.TLSWithClientSideCertificates,
		Base64SHA256Password: "DEADBEEF",
	}

	err = authStore.SetChargeStationAuth(ctx, "cs001", want)
	require.NoError(t, err)

	got, err := authStore.LookupChargeStationAuth(ctx, "cs001")
	require.NoError(t, err)

	assert.Equal(t, want, got)
}

func TestLookupChargeStationAuthWithUnregisteredChargeStation(t *testing.T) {
	ctx := context.Background()

	authStore, err := firestore.NewStore(ctx, "myproject")
	require.NoError(t, err)

	got, err := authStore.LookupChargeStationAuth(ctx, "not-created")
	require.NoError(t, err)
	assert.Nil(t, got)
}

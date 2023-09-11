// SPDX-License-Identifier: Apache-2.0

//go:build integration

package firestore_test

import (
	"context"
	"fmt"
	"k8s.io/utils/clock"
	clockTest "k8s.io/utils/clock/testing"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thoughtworks/maeve-csms/manager/store"
	"github.com/thoughtworks/maeve-csms/manager/store/firestore"
)

func TestSetAndLookupChargeStationAuth(t *testing.T) {
	ctx := context.Background()

	authStore, err := firestore.NewStore(ctx, "myproject", clock.RealClock{})
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

	authStore, err := firestore.NewStore(ctx, "myproject", clock.RealClock{})
	require.NoError(t, err)

	got, err := authStore.LookupChargeStationAuth(ctx, "not-created")
	require.NoError(t, err)
	assert.Nil(t, got)
}

func TestUpdateAndLookupChargeStationSettingsWithNewSettings(t *testing.T) {
	ctx := context.Background()

	now := time.Now()
	settingsStore, err := firestore.NewStore(ctx, "myproject", clockTest.NewFakePassiveClock(now))
	require.NoError(t, err)

	want := &store.ChargeStationSettings{
		ChargeStationId: "cs001",
		Settings: map[string]*store.ChargeStationSetting{
			"foo": {Value: "bar", Status: store.ChargeStationSettingStatusPending, LastUpdated: now.UTC()},
			"baz": {Value: "qux", Status: store.ChargeStationSettingStatusPending, LastUpdated: now.UTC()},
		},
	}

	err = settingsStore.UpdateChargeStationSettings(context.Background(), "cs001", want)
	require.NoError(t, err)

	got, err := settingsStore.LookupChargeStationSettings(context.Background(), "cs001")
	require.NoError(t, err)

	assert.Equal(t, want, got)
}

func TestUpdateAndLookupChargeStationSettingsWithUpdatedSettings(t *testing.T) {
	ctx := context.Background()

	settingsStore, err := firestore.NewStore(ctx, "myproject", clock.RealClock{})
	require.NoError(t, err)

	want := &store.ChargeStationSettings{
		ChargeStationId: "cs001",
		Settings: map[string]*store.ChargeStationSetting{
			"foo": {Value: "bar", Status: store.ChargeStationSettingStatusPending},
			"baz": {Value: "qux", Status: store.ChargeStationSettingStatusAccepted},
		},
	}

	err = settingsStore.UpdateChargeStationSettings(context.Background(), "cs001", &store.ChargeStationSettings{
		Settings: map[string]*store.ChargeStationSetting{
			"foo": {Value: "bar", Status: store.ChargeStationSettingStatusPending},
			"baz": {Value: "qux", Status: store.ChargeStationSettingStatusPending},
		},
	})
	require.NoError(t, err)

	err = settingsStore.UpdateChargeStationSettings(context.Background(), "cs001", &store.ChargeStationSettings{
		Settings: map[string]*store.ChargeStationSetting{
			"baz": {Value: "qux", Status: store.ChargeStationSettingStatusAccepted},
		},
	})
	require.NoError(t, err)

	got, err := settingsStore.LookupChargeStationSettings(context.Background(), "cs001")
	require.NoError(t, err)

	assert.Equal(t, want.ChargeStationId, got.ChargeStationId)
	assert.Len(t, got.Settings, len(want.Settings))
	assert.Equal(t, store.ChargeStationSettingStatusPending, got.Settings["foo"].Status)
	assert.Equal(t, store.ChargeStationSettingStatusAccepted, got.Settings["baz"].Status)
	assert.True(t, got.Settings["foo"].LastUpdated.Before(got.Settings["baz"].LastUpdated))
}

func TestListChargeStationSettings(t *testing.T) {
	ctx := context.Background()

	now := time.Now()
	settingsStore, err := firestore.NewStore(ctx, "myproject", clockTest.NewFakePassiveClock(now))
	require.NoError(t, err)

	want := &store.ChargeStationSettings{
		Settings: map[string]*store.ChargeStationSetting{
			"foo": {Value: "bar", Status: store.ChargeStationSettingStatusPending, LastUpdated: now.UTC()},
			"baz": {Value: "qux", Status: store.ChargeStationSettingStatusPending, LastUpdated: now.UTC()},
		},
	}
	for i := 0; i < 25; i++ {
		csId := fmt.Sprintf("cs%03d", i)
		err := settingsStore.UpdateChargeStationSettings(ctx, csId, want)
		require.NoError(t, err)
	}

	csIds := make(map[string]struct{})

	page1, err := settingsStore.ListChargeStationSettings(ctx, 10, "")
	require.NoError(t, err)
	require.Len(t, page1, 10)
	for _, got := range page1 {
		csIds[got.ChargeStationId] = struct{}{}
		assert.Equal(t, want.Settings, got.Settings)
	}

	page2, err := settingsStore.ListChargeStationSettings(ctx, 10, page1[len(page1)-1].ChargeStationId)
	require.NoError(t, err)
	require.Len(t, page2, 10)
	for _, got := range page2 {
		csIds[got.ChargeStationId] = struct{}{}
		assert.Equal(t, want.Settings, got.Settings)
	}

	page3, err := settingsStore.ListChargeStationSettings(ctx, 10, page2[len(page2)-1].ChargeStationId)
	require.NoError(t, err)
	require.Len(t, page3, 5)
	for _, got := range page3 {
		csIds[got.ChargeStationId] = struct{}{}
		assert.Equal(t, want.Settings, got.Settings)
	}

	assert.Len(t, csIds, 25)
}

func TestSetAndLookupChargeStationRuntimeDetails(t *testing.T) {
	ctx := context.Background()

	detailsStore, err := firestore.NewStore(ctx, "myproject", clock.RealClock{})
	require.NoError(t, err)

	want := &store.ChargeStationRuntimeDetails{
		OcppVersion: "1.6",
	}

	err = detailsStore.SetChargeStationRuntimeDetails(ctx, "cs001", want)
	require.NoError(t, err)

	got, err := detailsStore.LookupChargeStationRuntimeDetails(ctx, "cs001")
	require.NoError(t, err)

	assert.Equal(t, want, got)
}

func TestLookupChargeStationRuntimeDetailsWithUnregisteredChargeStation(t *testing.T) {
	ctx := context.Background()

	detailsStore, err := firestore.NewStore(ctx, "myproject", clock.RealClock{})
	require.NoError(t, err)

	got, err := detailsStore.LookupChargeStationRuntimeDetails(ctx, "not-created")
	require.NoError(t, err)
	assert.Nil(t, got)
}

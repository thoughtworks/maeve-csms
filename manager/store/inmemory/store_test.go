// SPDX-License-Identifier: Apache-2.0

package inmemory_test

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thoughtworks/maeve-csms/manager/store"
	"github.com/thoughtworks/maeve-csms/manager/store/inmemory"
	"k8s.io/utils/clock"
	clockTest "k8s.io/utils/clock/testing"
	"testing"
	"time"
)

func TestUpdateChargeStationSettingsWithNewSettings(t *testing.T) {
	now := time.Now()
	engine := inmemory.NewStore(clockTest.NewFakePassiveClock(now))

	want := &store.ChargeStationSettings{
		ChargeStationId: "cs001",
		Settings: map[string]*store.ChargeStationSetting{
			"foo": {Value: "bar", Status: store.ChargeStationSettingStatusPending, LastUpdated: now.UTC()},
			"baz": {Value: "qux", Status: store.ChargeStationSettingStatusPending, LastUpdated: now.UTC()},
		},
	}

	err := engine.UpdateChargeStationSettings(context.Background(), "cs001", want)
	require.NoError(t, err)

	got, err := engine.LookupChargeStationSettings(context.Background(), "cs001")
	require.NoError(t, err)

	assert.Equal(t, want, got)
}

func TestUpdateChargeStationSettingsWithExistingSettings(t *testing.T) {
	want := &store.ChargeStationSettings{
		ChargeStationId: "cs001",
		Settings: map[string]*store.ChargeStationSetting{
			"foo": {Value: "bar", Status: store.ChargeStationSettingStatusPending},
			"baz": {Value: "qux", Status: store.ChargeStationSettingStatusAccepted},
		},
	}

	engine := inmemory.NewStore(clock.RealClock{})
	err := engine.UpdateChargeStationSettings(context.Background(), "cs001", &store.ChargeStationSettings{
		Settings: map[string]*store.ChargeStationSetting{
			"foo": {Value: "bar", Status: store.ChargeStationSettingStatusPending},
			"baz": {Value: "qux", Status: store.ChargeStationSettingStatusPending},
		},
	})
	require.NoError(t, err)

	err = engine.UpdateChargeStationSettings(context.Background(), "cs001", &store.ChargeStationSettings{
		Settings: map[string]*store.ChargeStationSetting{
			"baz": {Value: "qux", Status: store.ChargeStationSettingStatusAccepted},
		},
	})
	require.NoError(t, err)

	got, err := engine.LookupChargeStationSettings(context.Background(), "cs001")
	require.NoError(t, err)

	assert.Equal(t, want.ChargeStationId, got.ChargeStationId)
	assert.Len(t, got.Settings, len(want.Settings))
	assert.Equal(t, store.ChargeStationSettingStatusPending, got.Settings["foo"].Status)
	assert.Equal(t, store.ChargeStationSettingStatusAccepted, got.Settings["baz"].Status)
	assert.True(t, got.Settings["foo"].LastUpdated.Before(got.Settings["baz"].LastUpdated))
}

func TestListChargeStationSettingsReturnsDataInPages(t *testing.T) {
	now := time.Now()
	engine := inmemory.NewStore(clockTest.NewFakePassiveClock(now))

	want := &store.ChargeStationSettings{
		ChargeStationId: "cs001",
		Settings: map[string]*store.ChargeStationSetting{
			"foo": {Value: "bar", Status: store.ChargeStationSettingStatusPending, LastUpdated: now.UTC()},
			"baz": {Value: "qux", Status: store.ChargeStationSettingStatusPending, LastUpdated: now.UTC()},
		},
	}
	for i := 0; i < 25; i++ {
		csId := fmt.Sprintf("cs%03d", i)
		err := engine.UpdateChargeStationSettings(context.Background(), csId, want)
		require.NoError(t, err)
	}

	csIds := make(map[string]struct{})

	page1, err := engine.ListChargeStationSettings(context.Background(), 10, "")
	require.NoError(t, err)
	require.Len(t, page1, 10)
	for _, got := range page1 {
		csIds[got.ChargeStationId] = struct{}{}
		assert.Equal(t, want.Settings, got.Settings)
	}

	page2, err := engine.ListChargeStationSettings(context.Background(), 10, page1[len(page1)-1].ChargeStationId)
	require.NoError(t, err)
	require.Len(t, page2, 10)
	for _, got := range page2 {
		csIds[got.ChargeStationId] = struct{}{}
		assert.Equal(t, want.Settings, got.Settings)
	}

	page3, err := engine.ListChargeStationSettings(context.Background(), 10, page2[len(page2)-1].ChargeStationId)
	require.NoError(t, err)
	require.Len(t, page3, 5)
	for _, got := range page3 {
		csIds[got.ChargeStationId] = struct{}{}
		assert.Equal(t, want.Settings, got.Settings)
	}

	assert.Len(t, csIds, 25)
}

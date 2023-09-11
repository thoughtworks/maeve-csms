// SPDX-License-Identifier: Apache-2.0

package ocpp16_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	handlers "github.com/thoughtworks/maeve-csms/manager/handlers/ocpp16"
	types "github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp16"
	"github.com/thoughtworks/maeve-csms/manager/store"
	"github.com/thoughtworks/maeve-csms/manager/store/inmemory"
	"k8s.io/utils/clock"
	clockTest "k8s.io/utils/clock/testing"
	"testing"
	"time"
)

func TestBootNotificationHandler(t *testing.T) {
	now, err := time.Parse(time.RFC3339, "2023-06-15T15:05:00+01:00")
	require.NoError(t, err)

	engine := inmemory.NewStore(clock.RealClock{})

	err = engine.UpdateChargeStationSettings(context.Background(), "cs001", &store.ChargeStationSettings{
		Settings: map[string]*store.ChargeStationSetting{
			"foo": {Value: "bar", Status: store.ChargeStationSettingStatusAccepted},
			"baz": {Value: "qux", Status: store.ChargeStationSettingStatusRebootRequired},
		},
	})
	require.NoError(t, err)

	handler := handlers.BootNotificationHandler{
		Clock:               clockTest.NewFakePassiveClock(now),
		RuntimeDetailsStore: engine,
		SettingsStore:       engine,
		HeartbeatInterval:   10,
	}

	serialNumber := "cs001-1234"
	req := &types.BootNotificationJson{
		ChargePointSerialNumber: &serialNumber,
	}

	got, err := handler.HandleCall(context.Background(), "cs001", req)
	assert.NoError(t, err)

	want := &types.BootNotificationResponseJson{
		CurrentTime: "2023-06-15T15:05:00+01:00",
		Status:      types.BootNotificationResponseJsonStatusAccepted,
		Interval:    10,
	}

	assert.Equal(t, want, got)

	details, err := engine.LookupChargeStationRuntimeDetails(context.Background(), "cs001")
	require.NoError(t, err)
	assert.Equal(t, store.ChargeStationRuntimeDetails{
		OcppVersion: "1.6",
	}, *details)

	settings, err := engine.LookupChargeStationSettings(context.Background(), "cs001")
	require.NoError(t, err)
	for _, v := range settings.Settings {
		assert.NotEqual(t, store.ChargeStationSettingStatusRebootRequired, v.Status)
	}
}

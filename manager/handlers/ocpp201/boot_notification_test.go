// SPDX-License-Identifier: Apache-2.0

package ocpp201_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	handlers "github.com/thoughtworks/maeve-csms/manager/handlers/ocpp201"
	types "github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"github.com/thoughtworks/maeve-csms/manager/store"
	"github.com/thoughtworks/maeve-csms/manager/store/inmemory"
	"k8s.io/utils/clock"
	clockTest "k8s.io/utils/clock/testing"
	"testing"
	"time"
)

func makePtr[T any](t T) *T {
	p := t
	return &p
}

func TestBootNotificationHandler(t *testing.T) {
	now, err := time.Parse(time.RFC3339, "2023-06-15T15:05:00+01:00")
	require.NoError(t, err)
	engine := inmemory.NewStore(clock.RealClock{})

	handler := handlers.BootNotificationHandler{
		Clock:               clockTest.NewFakePassiveClock(now),
		RuntimeDetailsStore: engine,
		HeartbeatInterval:   10,
	}

	req := &types.BootNotificationRequestJson{
		ChargingStation: types.ChargingStationType{
			Model:        "testy",
			SerialNumber: makePtr("cs001"),
		},
		Reason: types.BootReasonEnumTypePowerUp,
	}

	got, err := handler.HandleCall(context.Background(), "cs001", req)
	assert.NoError(t, err)

	want := &types.BootNotificationResponseJson{
		CurrentTime: "2023-06-15T15:05:00+01:00",
		Status:      types.RegistrationStatusEnumTypeAccepted,
		Interval:    10,
	}

	assert.Equal(t, want, got)

	details, err := engine.LookupChargeStationRuntimeDetails(context.Background(), "cs001")
	require.NoError(t, err)
	assert.Equal(t, store.ChargeStationRuntimeDetails{
		OcppVersion: "2.0.1",
	}, *details)
}

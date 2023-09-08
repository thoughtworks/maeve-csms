// SPDX-License-Identifier: Apache-2.0

package mqtt_test

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thoughtworks/maeve-csms/manager/mqtt"
	"github.com/thoughtworks/maeve-csms/manager/ocpp"
	"github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp16"
	"github.com/thoughtworks/maeve-csms/manager/store"
	"github.com/thoughtworks/maeve-csms/manager/store/inmemory"
	"k8s.io/utils/clock"
	"testing"
	"time"
)

type callEvent struct {
	chargeStationId string
	request         ocpp.Request
}

type mockCallMaker struct {
	engine     store.Engine
	callEvents []callEvent
	updateFn   updateFn
}

type updateFn func(ctx context.Context, engine store.Engine, chargeStationId string, req ocpp.Request) error

func (m *mockCallMaker) Send(ctx context.Context, chargeStationId string, request ocpp.Request) error {
	m.callEvents = append(m.callEvents, callEvent{
		chargeStationId: chargeStationId,
		request:         request,
	})
	return m.updateFn(ctx, m.engine, chargeStationId, request)
}

func updateV16StatusToAccepted(ctx context.Context, engine store.Engine, chargeStationId string, request ocpp.Request) error {
	req := request.(*ocpp16.ChangeConfigurationJson)
	return engine.UpdateChargeStationSettings(ctx, chargeStationId, &store.ChargeStationSettings{
		Settings: map[string]*store.ChargeStationSetting{
			req.Key: {Value: req.Value, Status: store.ChargeStationSettingStatusAccepted},
		},
	})
}

type updateV16ReqWithNoResponse struct {
	updateAttempts []time.Time
}

func (u *updateV16ReqWithNoResponse) update(ctx context.Context, engine store.Engine, chargeStationId string, req ocpp.Request) error {
	u.updateAttempts = append(u.updateAttempts, time.Now())
	return nil
}

func TestSyncV16Settings(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	engine := inmemory.NewStore(clock.RealClock{})

	err := engine.SetChargeStationRuntimeDetails(ctx, "cs001", &store.ChargeStationRuntimeDetails{
		OcppVersion: "1.6",
	})
	require.NoError(t, err)
	err = engine.SetChargeStationRuntimeDetails(ctx, "cs002", &store.ChargeStationRuntimeDetails{
		OcppVersion: "1.6",
	})
	require.NoError(t, err)
	err = engine.UpdateChargeStationSettings(ctx, "cs001", &store.ChargeStationSettings{
		Settings: map[string]*store.ChargeStationSetting{
			"foo": {Value: "bar", Status: "Pending"},
			"baz": {Value: "qux", Status: "Pending"},
		},
	})
	require.NoError(t, err)
	err = engine.UpdateChargeStationSettings(ctx, "cs002", &store.ChargeStationSettings{
		Settings: map[string]*store.ChargeStationSetting{
			"foo": {Value: "bar", Status: "Pending"},
		},
	})
	require.NoError(t, err)

	v16CallMaker := &mockCallMaker{engine: engine, updateFn: updateV16StatusToAccepted}
	mqtt.SyncSettings(ctx, engine, v16CallMaker, nil, 100*time.Millisecond, 500*time.Millisecond)

	settings, err := engine.LookupChargeStationSettings(ctx, "cs001")
	require.NoError(t, err)
	assert.Len(t, settings.Settings, 2)
	for _, v := range settings.Settings {
		assert.Equal(t, store.ChargeStationSettingStatusAccepted, v.Status)
	}

	settings, err = engine.LookupChargeStationSettings(ctx, "cs002")
	require.NoError(t, err)
	assert.Len(t, settings.Settings, 1)
	for _, v := range settings.Settings {
		assert.Equal(t, store.ChargeStationSettingStatusAccepted, v.Status)
	}
}

func TestSyncV16SettingsRetryAfterDelay(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	engine := inmemory.NewStore(clock.RealClock{})

	err := engine.SetChargeStationRuntimeDetails(ctx, "cs001", &store.ChargeStationRuntimeDetails{
		OcppVersion: "1.6",
	})
	require.NoError(t, err)

	err = engine.UpdateChargeStationSettings(ctx, "cs001", &store.ChargeStationSettings{
		Settings: map[string]*store.ChargeStationSetting{
			"foo": {Value: "bar", Status: "Pending"},
		},
	})
	require.NoError(t, err)

	updater := updateV16ReqWithNoResponse{}
	v16CallMaker := &mockCallMaker{engine: engine, updateFn: updater.update}
	mqtt.SyncSettings(ctx, engine, v16CallMaker, nil, 100*time.Millisecond, 400*time.Millisecond)

	require.Equal(t, 2, len(updater.updateAttempts))
	assert.True(t, updater.updateAttempts[1].After(updater.updateAttempts[0].Add(400*time.Millisecond)))
}

func TestSyncSettingsWithManyChargeStations(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	engine := inmemory.NewStore(clock.RealClock{})

	for i := 0; i < 300; i++ {
		csId := fmt.Sprintf("cs%3d", i)
		err := engine.SetChargeStationRuntimeDetails(ctx, csId, &store.ChargeStationRuntimeDetails{
			OcppVersion: "1.6",
		})
		require.NoError(t, err)
		err = engine.UpdateChargeStationSettings(ctx, csId, &store.ChargeStationSettings{
			Settings: map[string]*store.ChargeStationSetting{
				"foo": {Value: "bar", Status: "Accepted"},
			},
		})
		require.NoError(t, err)
	}

	err := engine.UpdateChargeStationSettings(ctx, "cs275", &store.ChargeStationSettings{
		Settings: map[string]*store.ChargeStationSetting{
			"foo": {Value: "bar2", Status: "Pending"},
		},
	})
	require.NoError(t, err)

	v16CallMaker := &mockCallMaker{engine: engine, updateFn: updateV16StatusToAccepted}
	mqtt.SyncSettings(ctx, engine, v16CallMaker, nil, 100*time.Millisecond, 500*time.Millisecond)

	settings, err := engine.LookupChargeStationSettings(ctx, "cs275")
	require.NoError(t, err)
	assert.Len(t, settings.Settings, 1)
	for _, v := range settings.Settings {
		assert.Equal(t, store.ChargeStationSettingStatusAccepted, v.Status)
	}
}

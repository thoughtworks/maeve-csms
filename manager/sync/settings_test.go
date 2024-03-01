// SPDX-License-Identifier: Apache-2.0

package sync_test

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thoughtworks/maeve-csms/manager/ocpp"
	"github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp16"
	"github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"github.com/thoughtworks/maeve-csms/manager/store"
	"github.com/thoughtworks/maeve-csms/manager/store/inmemory"
	"github.com/thoughtworks/maeve-csms/manager/sync"
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
	if m.updateFn != nil {
		return m.updateFn(ctx, m.engine, chargeStationId, request)
	}
	return nil
}

func updateV16StatusToAccepted(ctx context.Context, engine store.Engine, chargeStationId string, request ocpp.Request) error {
	req := request.(*ocpp16.ChangeConfigurationJson)
	return engine.UpdateChargeStationSettings(ctx, chargeStationId, &store.ChargeStationSettings{
		Settings: map[string]*store.ChargeStationSetting{
			req.Key: {Value: req.Value, Status: store.ChargeStationSettingStatusAccepted},
		},
	})
}

func updateV201StatusToAccepted(ctx context.Context, engine store.Engine, chargeStationId string, request ocpp.Request) error {
	req := request.(*ocpp201.SetVariablesRequestJson)
	updatedSettings := make(map[string]*store.ChargeStationSetting)
	for _, variable := range req.SetVariableData {
		name := variableDataToName(variable)
		updatedSettings[name] = &store.ChargeStationSetting{
			Value:  variable.AttributeValue,
			Status: store.ChargeStationSettingStatusAccepted,
		}
	}
	return engine.UpdateChargeStationSettings(ctx, chargeStationId, &store.ChargeStationSettings{
		Settings: updatedSettings,
	})
}

func variableDataToName(variable ocpp201.SetVariableDataType) string {
	s := variable.Component.Name
	if variable.Component.Instance != nil || variable.Component.Evse != nil {
		s += ";"
		if variable.Component.Instance != nil {
			s += *variable.Component.Instance
		}
		if variable.Component.Evse != nil {
			s += fmt.Sprintf(";%d", variable.Component.Evse.Id)
		}
	}
	s += fmt.Sprintf("/%s", variable.Variable.Name)
	if variable.Variable.Instance != nil || variable.AttributeType != nil {
		s += ";"
		if variable.Variable.Instance != nil {
			s += *variable.Variable.Instance
		}
		if variable.AttributeType != nil {
			s += fmt.Sprintf(";%s", *variable.AttributeType)
		}
	}
	return s
}

type updateWithNoResponse struct {
	updateAttempts []time.Time
}

func (u *updateWithNoResponse) update(context.Context, store.Engine, string, ocpp.Request) error {
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
	sync.SyncSettings(ctx, engine, clock.RealClock{}, v16CallMaker, nil, 100*time.Millisecond, 500*time.Millisecond)

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

	updater := updateWithNoResponse{}
	v16CallMaker := &mockCallMaker{engine: engine, updateFn: updater.update}
	sync.SyncSettings(ctx, engine, clock.RealClock{}, v16CallMaker, nil, 100*time.Millisecond, 400*time.Millisecond)

	require.Equal(t, 3, len(updater.updateAttempts))
	assert.True(t, updater.updateAttempts[1].After(updater.updateAttempts[0].Add(400*time.Millisecond)))
	assert.True(t, updater.updateAttempts[2].After(updater.updateAttempts[1].Add(400*time.Millisecond)))
}

func TestSyncV201Variables(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	engine := inmemory.NewStore(clock.RealClock{})

	err := engine.SetChargeStationRuntimeDetails(ctx, "cs001", &store.ChargeStationRuntimeDetails{
		OcppVersion: "2.0.1",
	})
	require.NoError(t, err)
	err = engine.SetChargeStationRuntimeDetails(ctx, "cs002", &store.ChargeStationRuntimeDetails{
		OcppVersion: "2.0.1",
	})
	require.NoError(t, err)
	err = engine.UpdateChargeStationSettings(ctx, "cs001", &store.ChargeStationSettings{
		Settings: map[string]*store.ChargeStationSetting{
			"foo/fam": {Value: "bar", Status: "Pending"},
			"baz/bam": {Value: "qux", Status: "Pending"},
		},
	})
	require.NoError(t, err)
	err = engine.UpdateChargeStationSettings(ctx, "cs002", &store.ChargeStationSettings{
		Settings: map[string]*store.ChargeStationSetting{
			"foo/fam": {Value: "bar", Status: "Pending"},
		},
	})
	require.NoError(t, err)

	v201CallMaker := &mockCallMaker{engine: engine, updateFn: updateV201StatusToAccepted}
	sync.SyncSettings(ctx, engine, clock.RealClock{}, nil, v201CallMaker, 100*time.Millisecond, 500*time.Millisecond)

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

func TestSyncV201SettingsRetryAfterDelay(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	engine := inmemory.NewStore(clock.RealClock{})

	err := engine.SetChargeStationRuntimeDetails(ctx, "cs001", &store.ChargeStationRuntimeDetails{
		OcppVersion: "2.0.1",
	})
	require.NoError(t, err)

	err = engine.UpdateChargeStationSettings(ctx, "cs001", &store.ChargeStationSettings{
		Settings: map[string]*store.ChargeStationSetting{
			"foo/fam": {Value: "bar", Status: "Pending"},
		},
	})
	require.NoError(t, err)

	updater := updateWithNoResponse{}
	v201CallMaker := &mockCallMaker{engine: engine, updateFn: updater.update}
	sync.SyncSettings(ctx, engine, clock.RealClock{}, nil, v201CallMaker, 100*time.Millisecond, 400*time.Millisecond)

	require.Equal(t, 3, len(updater.updateAttempts))
	assert.True(t, updater.updateAttempts[1].After(updater.updateAttempts[0].Add(400*time.Millisecond)))
	assert.True(t, updater.updateAttempts[2].After(updater.updateAttempts[1].Add(400*time.Millisecond)))
}

func TestSyncSettingsWithManyChargeStations(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	engine := inmemory.NewStore(clock.RealClock{})

	for i := 0; i < 150; i++ {
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
	for i := 150; i < 300; i++ {
		csId := fmt.Sprintf("cs%3d", i)
		err := engine.SetChargeStationRuntimeDetails(ctx, csId, &store.ChargeStationRuntimeDetails{
			OcppVersion: "2.0.1",
		})
		require.NoError(t, err)
		err = engine.UpdateChargeStationSettings(ctx, csId, &store.ChargeStationSettings{
			Settings: map[string]*store.ChargeStationSetting{
				"foo/fam": {Value: "bar", Status: "Accepted"},
			},
		})
		require.NoError(t, err)
	}

	err := engine.UpdateChargeStationSettings(ctx, "cs133", &store.ChargeStationSettings{
		Settings: map[string]*store.ChargeStationSetting{
			"foo": {Value: "bar2", Status: "Pending"},
		},
	})
	require.NoError(t, err)

	err = engine.UpdateChargeStationSettings(ctx, "cs275", &store.ChargeStationSettings{
		Settings: map[string]*store.ChargeStationSetting{
			"foo/fam": {Value: "bar2", Status: "Pending"},
		},
	})
	require.NoError(t, err)

	v16CallMaker := &mockCallMaker{engine: engine, updateFn: updateV16StatusToAccepted}
	v201CallMaker := &mockCallMaker{engine: engine, updateFn: updateV201StatusToAccepted}
	sync.SyncSettings(ctx, engine, clock.RealClock{}, v16CallMaker, v201CallMaker, 100*time.Millisecond, 500*time.Millisecond)

	settings, err := engine.LookupChargeStationSettings(ctx, "cs133")
	require.NoError(t, err)
	assert.Len(t, settings.Settings, 1)
	for _, v := range settings.Settings {
		assert.Equal(t, store.ChargeStationSettingStatusAccepted, v.Status)
	}

	settings, err = engine.LookupChargeStationSettings(ctx, "cs275")
	require.NoError(t, err)
	assert.Len(t, settings.Settings, 1)
	for _, v := range settings.Settings {
		assert.Equal(t, store.ChargeStationSettingStatusAccepted, v.Status)
	}
}

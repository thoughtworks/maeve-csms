// SPDX-License-Identifier: Apache-2.0

package ocpp201_test

import (
	"context"
	"github.com/thoughtworks/maeve-csms/manager/store"
	"k8s.io/utils/clock"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	handlers "github.com/thoughtworks/maeve-csms/manager/handlers/ocpp201"
	types "github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"github.com/thoughtworks/maeve-csms/manager/services"
	"github.com/thoughtworks/maeve-csms/manager/store/inmemory"
)

func TestTransactionEventHandlerWithStartedEvent(t *testing.T) {
	ctx := context.Background()
	engine := inmemory.NewStore(clock.RealClock{})
	tariffService := services.BasicKwhTariffService{}

	err := engine.SetToken(ctx, &store.Token{
		Uid: "SOMERFID",
	})
	require.NoError(t, err)

	handler := handlers.TransactionEventHandler{
		Store:         engine,
		TariffService: tariffService,
	}

	req := &types.TransactionEventRequestJson{
		EventType:     types.TransactionEventEnumTypeStarted,
		TriggerReason: types.TriggerReasonEnumTypeCablePluggedIn,
		Timestamp:     "2023-05-05T12:00:00+01:00",
		IdToken: &types.IdTokenType{
			Type:    types.IdTokenEnumTypeISO14443,
			IdToken: "SOMERFID",
		},
		MeterValue: []types.MeterValueType{
			{
				Timestamp: "2023-05-05T12:00:00+01:00",
				SampledValue: []types.SampledValueType{
					{
						Measurand: makePtr(types.MeasurandEnumTypeEnergyActiveImportRegister),
						Location:  makePtr(types.LocationEnumTypeOutlet),
						Value:     100,
					},
				},
			},
		},
		SeqNo: 0,
		TransactionInfo: types.TransactionType{
			TransactionId: "5555",
			ChargingState: makePtr(types.ChargingStateEnumTypeCharging),
		},
	}

	got, err := handler.HandleCall(ctx, "cs001", req)
	assert.NoError(t, err)

	want := &types.TransactionEventResponseJson{
		IdTokenInfo: &types.IdTokenInfoType{
			Status: types.AuthorizationStatusEnumTypeAccepted,
		},
	}
	assert.Equal(t, want, got)

	transaction, err := engine.FindTransaction(ctx, "cs001", "5555")
	require.NoError(t, err)
	assert.NotNil(t, transaction)
}

func TestTransactionEventHandlerWithStartedEventWithInvalidToken(t *testing.T) {
	ctx := context.Background()
	engine := inmemory.NewStore(clock.RealClock{})
	tariffService := services.BasicKwhTariffService{}

	handler := handlers.TransactionEventHandler{
		Store:         engine,
		TariffService: tariffService,
	}

	req := &types.TransactionEventRequestJson{
		EventType:     types.TransactionEventEnumTypeStarted,
		TriggerReason: types.TriggerReasonEnumTypeCablePluggedIn,
		Timestamp:     "2023-05-05T12:00:00+01:00",
		IdToken: &types.IdTokenType{
			Type:    types.IdTokenEnumTypeISO14443,
			IdToken: "SOMERFID",
		},
		MeterValue: []types.MeterValueType{
			{
				Timestamp: "2023-05-05T12:00:00+01:00",
				SampledValue: []types.SampledValueType{
					{
						Measurand: makePtr(types.MeasurandEnumTypeEnergyActiveImportRegister),
						Location:  makePtr(types.LocationEnumTypeOutlet),
						Value:     100,
					},
				},
			},
		},
		SeqNo: 0,
		TransactionInfo: types.TransactionType{
			TransactionId: "5555",
			ChargingState: makePtr(types.ChargingStateEnumTypeCharging),
		},
	}

	got, err := handler.HandleCall(ctx, "cs001", req)
	assert.NoError(t, err)

	want := &types.TransactionEventResponseJson{
		IdTokenInfo: &types.IdTokenInfoType{
			Status: types.AuthorizationStatusEnumTypeUnknown,
		},
	}
	assert.Equal(t, want, got)

	transaction, err := engine.FindTransaction(ctx, "cs001", "5555")
	require.NoError(t, err)
	assert.NotNil(t, transaction)
}

func TestTransactionEventHandlerWithUpdatedEvent(t *testing.T) {
	ctx := context.Background()
	engine := inmemory.NewStore(clock.RealClock{})
	tariffService := services.BasicKwhTariffService{}

	handler := handlers.TransactionEventHandler{
		Store:         engine,
		TariffService: tariffService,
	}

	req := &types.TransactionEventRequestJson{
		EventType:     types.TransactionEventEnumTypeUpdated,
		TriggerReason: types.TriggerReasonEnumTypeMeterValuePeriodic,
		Timestamp:     "2023-05-05T12:00:00+01:00",
		MeterValue: []types.MeterValueType{
			{
				Timestamp: "2023-05-05T12:00:00+01:00",
				SampledValue: []types.SampledValueType{
					{
						Measurand: makePtr(types.MeasurandEnumTypeEnergyActiveImportRegister),
						Location:  makePtr(types.LocationEnumTypeOutlet),
						Value:     100,
					},
				},
			},
		},
		SeqNo: 0,
		TransactionInfo: types.TransactionType{
			TransactionId: "5555",
			ChargingState: makePtr(types.ChargingStateEnumTypeCharging),
		},
	}

	got, err := handler.HandleCall(ctx, "cs001", req)
	assert.NoError(t, err)

	want := &types.TransactionEventResponseJson{}
	assert.Equal(t, want, got)

	transaction, err := engine.FindTransaction(ctx, "cs001", "5555")
	require.NoError(t, err)
	assert.NotNil(t, transaction)
}

func TestTransactionEventHandlerWithEndedEvent(t *testing.T) {
	ctx := context.Background()
	engine := inmemory.NewStore(clock.RealClock{})
	tariffService := services.BasicKwhTariffService{}

	handler := handlers.TransactionEventHandler{
		Store:         engine,
		TariffService: tariffService,
	}

	req := &types.TransactionEventRequestJson{
		EventType:     types.TransactionEventEnumTypeEnded,
		TriggerReason: types.TriggerReasonEnumTypeStopAuthorized,
		Timestamp:     "2023-05-05T12:00:00+01:00",
		MeterValue: []types.MeterValueType{
			{
				Timestamp: "2023-05-05T12:00:00+01:00",
				SampledValue: []types.SampledValueType{
					{
						Context:   makePtr(types.ReadingContextEnumTypeTransactionEnd),
						Measurand: makePtr(types.MeasurandEnumTypeEnergyActiveImportRegister),
						Location:  makePtr(types.LocationEnumTypeOutlet),
						Value:     100,
					},
				},
			},
		},
		SeqNo: 0,
		TransactionInfo: types.TransactionType{
			TransactionId: "5555",
			ChargingState: makePtr(types.ChargingStateEnumTypeCharging),
		},
	}

	got, err := handler.HandleCall(ctx, "cs001", req)
	assert.NoError(t, err)

	want := &types.TransactionEventResponseJson{
		TotalCost: makePtr(0.055),
	}
	assert.Equal(t, want, got)

	transaction, err := engine.FindTransaction(ctx, "cs001", "5555")
	require.NoError(t, err)
	assert.NotNil(t, transaction)
}

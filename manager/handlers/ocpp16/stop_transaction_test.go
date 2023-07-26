// SPDX-License-Identifier: Apache-2.0

package ocpp16_test

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	handlers "github.com/thoughtworks/maeve-csms/manager/handlers/ocpp16"
	types "github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp16"
	"github.com/thoughtworks/maeve-csms/manager/store"
	"github.com/thoughtworks/maeve-csms/manager/store/inmemory"
	clockTest "k8s.io/utils/clock/testing"
)

func TestStopTransactionHandler(t *testing.T) {
	chargingStationId := fmt.Sprintf("cs%03d", rand.Intn(1000))
	engine := inmemory.NewStore()

	err := engine.SetToken(context.Background(), &store.Token{
		CountryCode: "GB",
		PartyId:     "TWK",
		Type:        "RFID",
		Uid:         "MYRFIDTAG",
		ContractId:  "GBTWK012345678V",
		Issuer:      "Thoughtworks",
		Valid:       true,
		CacheMode:   "NEVER",
		LastUpdated: time.Now().Format(time.RFC3339),
	})
	require.NoError(t, err)

	now, err := time.Parse(time.RFC3339, "2023-06-15T15:06:00+01:00")
	require.NoError(t, err)

	transactionStore := inmemory.NewStore()

	startContext := "Transaction.Begin"
	startMeasurand := "MeterValue"
	startLocation := "Outlet"
	err = transactionStore.CreateTransaction(context.TODO(), chargingStationId, handlers.ConvertToUUID(42), "MYRFIDTAG", "ISO14443",
		[]store.MeterValue{
			{
				SampledValues: []store.SampledValue{
					{
						Context:   &startContext,
						Measurand: &startMeasurand,
						Location:  &startLocation,
						Value:     50,
					},
				},
				Timestamp: now.Format(time.RFC3339),
			},
		}, 0, false)
	require.NoError(t, err)

	handler := handlers.StopTransactionHandler{
		Clock:            clockTest.NewFakePassiveClock(now),
		TokenStore:       engine,
		TransactionStore: transactionStore,
	}

	idTag := "MYRFIDTAG"
	reason := types.StopTransactionJsonReasonEVDisconnected
	periodicSampleContext := types.StopTransactionJsonTransactionDataElemSampledValueElemContextSamplePeriodic
	energyRegisterMeasurand := types.StopTransactionJsonTransactionDataElemSampledValueElemMeasurandEnergyActiveImportRegister
	outletLocation := types.StopTransactionJsonTransactionDataElemSampledValueElemLocationOutlet
	req := &types.StopTransactionJson{
		IdTag:     &idTag,
		MeterStop: 200,
		Reason:    &reason,
		Timestamp: now.Format(time.RFC3339),
		TransactionData: []types.StopTransactionJsonTransactionDataElem{
			{
				SampledValue: []types.StopTransactionJsonTransactionDataElemSampledValueElem{
					{
						Context:   &periodicSampleContext,
						Measurand: &energyRegisterMeasurand,
						Location:  &outletLocation,
						Value:     "100",
					},
				},
				Timestamp: now.Format(time.RFC3339),
			},
		},
		TransactionId: 42,
	}

	got, err := handler.HandleCall(context.Background(), chargingStationId, req)
	require.NoError(t, err)

	want := &types.StopTransactionResponseJson{
		IdTagInfo: &types.StopTransactionResponseJsonIdTagInfo{
			Status: types.StopTransactionResponseJsonIdTagInfoStatusAccepted,
		},
	}

	assert.Equal(t, want, got)

	found, err := transactionStore.FindTransaction(context.TODO(), chargingStationId, handlers.ConvertToUUID(42))
	require.NoError(t, err)

	expectedTransactionEndContext := "Transaction.End"
	expectedPeriodicContext := "Sample.Periodic"
	expectedOutletLocation := "Outlet"
	expectedMeasurand := "Energy.Active.Import.Register"
	expected := &store.Transaction{
		ChargeStationId: chargingStationId,
		TransactionId:   handlers.ConvertToUUID(42),
		IdToken:         "MYRFIDTAG",
		TokenType:       "ISO14443",
		MeterValues: []store.MeterValue{
			{
				Timestamp: now.Format(time.RFC3339),
				SampledValues: []store.SampledValue{
					{
						Context:   &startContext,
						Location:  &startLocation,
						Measurand: &startMeasurand,
						Value:     50,
					},
				},
			},
			{
				Timestamp: now.Format(time.RFC3339),
				SampledValues: []store.SampledValue{
					{
						Context:   &expectedPeriodicContext,
						Location:  &expectedOutletLocation,
						Measurand: &expectedMeasurand,
						Value:     100,
					},
				},
			},
			{
				Timestamp: now.Format(time.RFC3339),
				SampledValues: []store.SampledValue{
					{
						Context:   &expectedTransactionEndContext,
						Location:  &expectedOutletLocation,
						Measurand: &expectedMeasurand,
						Value:     150,
					},
				},
			},
		},
		StartSeqNo:        0,
		EndedSeqNo:        1,
		UpdatedSeqNoCount: 0,
		Offline:           false,
	}

	assert.Equal(t, expected, found)
}

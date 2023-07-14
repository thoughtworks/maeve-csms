// SPDX-License-Identifier: Apache-2.0

package ocpp16_test

import (
	"context"
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

func TestStartTransaction(t *testing.T) {
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

	transactionStore := inmemory.NewStore()

	now, err := time.Parse(time.RFC3339, "2023-06-15T15:05:00+01:00")
	require.NoError(t, err)

	handler := handlers.StartTransactionHandler{
		Clock:            clockTest.NewFakePassiveClock(now),
		TokenStore:       engine,
		TransactionStore: transactionStore,
	}

	req := &types.StartTransactionJson{
		ConnectorId:   1,
		IdTag:         "MYRFIDTAG",
		MeterStart:    100,
		ReservationId: nil,
		Timestamp:     now.Format(time.RFC3339),
	}

	ctx := context.Background()
	resp, err := handler.HandleCall(ctx, "cs001", req)
	require.NoError(t, err)
	got := resp.(*types.StartTransactionResponseJson)

	want := &types.StartTransactionResponseJson{
		IdTagInfo: types.StartTransactionResponseJsonIdTagInfo{
			Status: types.StartTransactionResponseJsonIdTagInfoStatusAccepted,
		},
		TransactionId: 0,
	}

	assert.Equal(t, want.IdTagInfo, got.IdTagInfo)
	assert.Equal(t, 0, want.TransactionId>>31)

	require.NoError(t, err)
	transactionId := handlers.ConvertToUUID(got.TransactionId)
	found, err := transactionStore.FindTransaction(ctx, "cs001", transactionId)
	require.NoError(t, err)

	expectedContext := "Transaction.Begin"
	expectedMeasurand := "MeterValue"
	expected := &store.Transaction{
		ChargeStationId: "cs001",
		TransactionId:   transactionId,
		IdToken:         "MYRFIDTAG",
		TokenType:       "ISO14443",
		MeterValues: []store.MeterValue{
			{
				Timestamp: now.Format(time.RFC3339),
				SampledValues: []store.SampledValue{
					{
						Context:   &expectedContext,
						Measurand: &expectedMeasurand,
						UnitOfMeasure: &store.UnitOfMeasure{
							Unit:      "Wh",
							Multipler: 1,
						},
						Value: 100,
					},
				},
			},
		},
		StartSeqNo:        0,
		EndedSeqNo:        0,
		UpdatedSeqNoCount: 0,
		Offline:           false,
	}

	assert.Equal(t, expected, found)
}

func TestStartTransactionWithInvalidRFID(t *testing.T) {
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

	transactionStore := inmemory.NewStore()

	now, err := time.Parse(time.RFC3339, "2023-06-15T15:05:00+01:00")
	require.NoError(t, err)

	handler := handlers.StartTransactionHandler{
		Clock:            clockTest.NewFakePassiveClock(now),
		TokenStore:       engine,
		TransactionStore: transactionStore,
	}

	req := &types.StartTransactionJson{
		ConnectorId:   1,
		IdTag:         "BADRFIDTAG",
		MeterStart:    0,
		ReservationId: nil,
		Timestamp:     now.Format(time.RFC3339),
	}

	got, err := handler.HandleCall(context.Background(), "cs001", req)
	require.NoError(t, err)

	want := &types.StartTransactionResponseJson{
		IdTagInfo: types.StartTransactionResponseJsonIdTagInfo{
			Status: types.StartTransactionResponseJsonIdTagInfoStatusInvalid,
		},
		TransactionId: -1,
	}

	assert.Equal(t, want, got)
}

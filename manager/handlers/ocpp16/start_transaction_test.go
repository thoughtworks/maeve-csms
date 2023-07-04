package ocpp16_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	handlers "github.com/twlabs/ocpp2-broker-core/manager/handlers/ocpp16"
	types "github.com/twlabs/ocpp2-broker-core/manager/ocpp/ocpp16"
	"github.com/twlabs/ocpp2-broker-core/manager/services"
	clockTest "k8s.io/utils/clock/testing"
	"testing"
	"time"
)

func TestStartTransaction(t *testing.T) {
	tokenStore := services.InMemoryTokenStore{
		Tokens: map[string]*services.Token{
			"ISO14443:MYRFIDTAG": {
				Type: "ISO14443",
				Uid:  "MYRFIDTAG",
			},
		},
	}
	transactionStore := services.NewInMemoryTransactionStore()

	now, err := time.Parse(time.RFC3339, "2023-06-15T15:05:00+01:00")
	require.NoError(t, err)

	handler := handlers.StartTransactionHandler{
		Clock:            clockTest.NewFakePassiveClock(now),
		TokenStore:       tokenStore,
		TransactionStore: transactionStore,
	}

	req := &types.StartTransactionJson{
		ConnectorId:   1,
		IdTag:         "MYRFIDTAG",
		MeterStart:    100,
		ReservationId: nil,
		Timestamp:     now.Format(time.RFC3339),
	}

	resp, err := handler.HandleCall(context.Background(), "cs001", req)
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
	found, err := transactionStore.FindTransaction("cs001", transactionId)
	require.NoError(t, err)

	expectedContext := "Transaction.Begin"
	expectedMeasurand := "MeterValue"
	expected := &services.Transaction{
		ChargeStationId: "cs001",
		TransactionId:   transactionId,
		IdToken:         "MYRFIDTAG",
		TokenType:       "ISO14443",
		MeterValues: []services.MeterValue{
			{
				Timestamp: now.Format(time.RFC3339),
				SampledValues: []services.SampledValue{
					{
						Context:   &expectedContext,
						Measurand: &expectedMeasurand,
						UnitOfMeasure: &services.UnitOfMeasure{
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
	tokenStore := services.InMemoryTokenStore{
		Tokens: map[string]*services.Token{
			"ISO14443:MYRFIDTAG": {
				Type: "ISO14443",
				Uid:  "MYRFIDTAG",
			},
		},
	}
	transactionStore := services.NewInMemoryTransactionStore()

	now, err := time.Parse(time.RFC3339, "2023-06-15T15:05:00+01:00")
	require.NoError(t, err)

	handler := handlers.StartTransactionHandler{
		Clock:            clockTest.NewFakePassiveClock(now),
		TokenStore:       tokenStore,
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

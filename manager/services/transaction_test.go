package services_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twlabs/maeve-csms/manager/services"
)

const idToken = "SOMERFID"
const tokenType = "ISO14443"

func NewMeterValues(energyReactiveExportValue float64) []services.MeterValue {
	return []services.MeterValue{
		{
			Timestamp: time.Now().Format(time.RFC3339),
			SampledValues: []services.SampledValue{
				{
					Measurand: makePtr("Energy.Active.Import.Register"),
					Value:     energyReactiveExportValue,
				},
			},
		},
	}
}

func TestTransactionStoreCreateTransactionWithNewTransaction(t *testing.T) {
	tests := map[string]services.TransactionStore{
		"InMemory": services.NewInMemoryTransactionStore(),
		"Redis":    services.NewRedisTransactionStore(endpoint),
	}

	for name, transactionStore := range tests {
		t.Run(name, func(t *testing.T) {
			assert.NotNil(t, transactionStore)

			meterValues := NewMeterValues(100)

			err := transactionStore.CreateTransaction("cs001", "1234", idToken, tokenType, meterValues, 0, false)
			assert.NoError(t, err)

			got, err := transactionStore.FindTransaction("cs001", "1234")
			assert.NoError(t, err)

			want := &services.Transaction{
				ChargeStationId: "cs001",
				TransactionId:   "1234",
				IdToken:         idToken,
				TokenType:       tokenType,
				MeterValues:     meterValues,
				StartSeqNo:      0,
			}

			assert.Equal(t, want, got)
		})
	}
}

func TestTransactionStoreCreateTransactionWithExistingTransaction(t *testing.T) {
	tests := map[string]services.TransactionStore{
		"InMemory": services.NewInMemoryTransactionStore(),
		"Redis":    services.NewRedisTransactionStore(endpoint),
	}

	for name, transactionStore := range tests {
		t.Run(name, func(t *testing.T) {
			updateMeterValues := NewMeterValues(150)

			err := transactionStore.UpdateTransaction("cs002", "1234", updateMeterValues)
			require.NoError(t, err)

			createMeterValues := NewMeterValues(100)

			err = transactionStore.CreateTransaction("cs002", "1234", idToken, tokenType, createMeterValues, 2, false)
			require.NoError(t, err)

			got, err := transactionStore.FindTransaction("cs002", "1234")
			require.NoError(t, err)

			want := &services.Transaction{
				ChargeStationId:   "cs002",
				TransactionId:     "1234",
				IdToken:           idToken,
				TokenType:         tokenType,
				MeterValues:       append(updateMeterValues, createMeterValues...),
				StartSeqNo:        2,
				EndedSeqNo:        0,
				UpdatedSeqNoCount: 1,
				Offline:           false,
			}

			assert.Equal(t, want, got)
		})
	}
}

func TestTransactionStoreUpdateCreatedTransaction(t *testing.T) {
	tests := map[string]services.TransactionStore{
		"InMemory": services.NewInMemoryTransactionStore(),
		"Redis":    services.NewRedisTransactionStore(endpoint),
	}

	for name, transactionStore := range tests {
		t.Run(name, func(t *testing.T) {
			createMeterValues := NewMeterValues(100)

			err := transactionStore.CreateTransaction("cs003", "1234", idToken, tokenType, createMeterValues, 2, false)
			require.NoError(t, err)

			updateMeterValues := NewMeterValues(150)

			err = transactionStore.UpdateTransaction("cs003", "1234", updateMeterValues)
			require.NoError(t, err)

			got, err := transactionStore.FindTransaction("cs003", "1234")
			require.NoError(t, err)

			want := &services.Transaction{
				ChargeStationId:   "cs003",
				TransactionId:     "1234",
				IdToken:           idToken,
				TokenType:         tokenType,
				MeterValues:       append(createMeterValues, updateMeterValues...),
				StartSeqNo:        2,
				EndedSeqNo:        0,
				UpdatedSeqNoCount: 1,
				Offline:           false,
			}

			assert.Equal(t, want, got)
		})
	}
}

func TestTransactionStoreEndTransaction(t *testing.T) {
	tests := map[string]services.TransactionStore{
		"InMemory": services.NewInMemoryTransactionStore(),
		"Redis":    services.NewRedisTransactionStore(endpoint),
	}

	for name, transactionStore := range tests {
		t.Run(name, func(t *testing.T) {
			createMeterValues := NewMeterValues(100)

			err := transactionStore.CreateTransaction("cs004", "1234", idToken, tokenType, createMeterValues, 0, false)
			require.NoError(t, err)

			updateMeterValues := NewMeterValues(150)

			err = transactionStore.UpdateTransaction("cs004", "1234", updateMeterValues)
			require.NoError(t, err)

			endMeterValues := NewMeterValues(200)

			err = transactionStore.EndTransaction("cs004", "1234", idToken, tokenType, endMeterValues, 2)
			require.NoError(t, err)

			got, err := transactionStore.FindTransaction("cs004", "1234")
			require.NoError(t, err)

			want := &services.Transaction{
				ChargeStationId:   "cs004",
				TransactionId:     "1234",
				IdToken:           idToken,
				TokenType:         tokenType,
				MeterValues:       append(createMeterValues, append(updateMeterValues, endMeterValues...)...),
				StartSeqNo:        0,
				EndedSeqNo:        2,
				UpdatedSeqNoCount: 1,
				Offline:           false,
			}

			assert.Equal(t, want, got)
		})
	}
}

func TestTransactionStoreEndTransactionNoCreate(t *testing.T) {
	tests := map[string]services.TransactionStore{
		"InMemory": services.NewInMemoryTransactionStore(),
		"Redis":    services.NewRedisTransactionStore(endpoint),
	}

	for name, transactionStore := range tests {
		t.Run(name, func(t *testing.T) {
			endMeterValues := NewMeterValues(200)

			err := transactionStore.EndTransaction("cs005", "1234", idToken, tokenType, endMeterValues, 2)
			require.NoError(t, err)

			got, err := transactionStore.FindTransaction("cs005", "1234")
			require.NoError(t, err)

			want := &services.Transaction{
				ChargeStationId:   "cs005",
				TransactionId:     "1234",
				IdToken:           idToken,
				TokenType:         tokenType,
				MeterValues:       endMeterValues,
				StartSeqNo:        0,
				EndedSeqNo:        2,
				UpdatedSeqNoCount: 0,
				Offline:           false,
			}

			assert.Equal(t, want, got)
		})
	}
}

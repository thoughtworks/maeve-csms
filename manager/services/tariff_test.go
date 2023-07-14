// SPDX-License-Identifier: Apache-2.0

package services_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/thoughtworks/maeve-csms/manager/services"
	"github.com/thoughtworks/maeve-csms/manager/store"
)

func makePtr[T any](t T) *T {
	v := t
	return &v
}

func TestBasicKwhTariffServiceCanCalculateCost(t *testing.T) {
	transaction := &store.Transaction{
		MeterValues: []store.MeterValue{
			{
				Timestamp: time.Now().Format(time.RFC3339),
				SampledValues: []store.SampledValue{
					{
						Context:   makePtr("Transaction.End"),
						Measurand: makePtr("Energy.Active.Import.Register"),
						Location:  makePtr("Outlet"),
						Value:     100,
					},
				},
			},
		},
	}
	tariffService := services.BasicKwhTariffService{}
	cost, err := tariffService.CalculateCost(transaction)
	assert.NoError(t, err)
	assert.Equal(t, 0.055, cost)
}

func TestBasicKwhTariffServiceErrorsWithNilTransaction(t *testing.T) {
	tariffService := services.BasicKwhTariffService{}
	cost, err := tariffService.CalculateCost(nil)
	assert.ErrorContains(t, err, "no transaction provided")
	var zero float64
	assert.Equal(t, zero, cost)
}

func TestBasicKwhTariffServiceErrorsWhenNoKwhReading(t *testing.T) {
	transaction := &store.Transaction{}
	tariffService := services.BasicKwhTariffService{}
	cost, err := tariffService.CalculateCost(transaction)
	assert.ErrorContains(t, err, "no output energy reading found in transaction")
	var zero float64
	assert.Equal(t, zero, cost)
}

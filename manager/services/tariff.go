package services

import (
	"errors"
	"fmt"
	"sort"
	"time"
)

type TariffService interface {
	CalculateCost(transaction *Transaction) (float64, error)
}

type BasicKwhTariffService struct{}

func (BasicKwhTariffService) CalculateCost(transaction *Transaction) (float64, error) {
	var cost float64

	if transaction == nil {
		return cost, errors.New("no transaction provided")
	}

	costPerWh := 0.55 / 1000
	Wh, found := findMostRecentOutletEnergyReading(transaction)
	if !found {
		return cost, fmt.Errorf("no output energy reading found in transaction")
	}
	cost = costPerWh * Wh

	return cost, nil
}

func findMostRecentOutletEnergyReading(transaction *Transaction) (float64, bool) {
	sort.Slice(transaction.MeterValues, func(i, j int) bool {
		ts1, err := time.Parse(time.RFC3339, transaction.MeterValues[i].Timestamp)
		if err != nil {
			return false
		}
		ts2, err := time.Parse(time.RFC3339, transaction.MeterValues[j].Timestamp)
		if err != nil {
			return false
		}

		return ts2.After(ts1)
	})

	var totalWh float64
	found := false

	for _, mv := range transaction.MeterValues {
		for _, sv := range mv.SampledValues {
			if sv.Context != nil && *sv.Context == "Transaction.End" &&
				sv.Measurand != nil && *sv.Measurand == "Energy.Active.Import.Register" &&
				sv.Location != nil && *sv.Location == "Outlet" {
				totalWh = sv.Value
				found = true
			}
		}
	}

	return totalWh, found
}

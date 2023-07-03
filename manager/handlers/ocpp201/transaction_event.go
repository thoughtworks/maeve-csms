package ocpp201

import (
	"context"
	"github.com/twlabs/ocpp2-broker-core/manager/ocpp"
	types "github.com/twlabs/ocpp2-broker-core/manager/ocpp/ocpp201"
	"github.com/twlabs/ocpp2-broker-core/manager/services"
	"log"
)

type TransactionEventHandler struct {
	TransactionStore services.TransactionStore
	TariffService    services.TariffService
}

func (t TransactionEventHandler) HandleCall(ctx context.Context, chargeStationId string, request ocpp.Request) (ocpp.Response, error) {
	req := request.(*types.TransactionEventRequestJson)
	log.Printf("Charge station %s transaction %s event: %s %s %d", chargeStationId, req.TransactionInfo.TransactionId, req.EventType, req.TriggerReason, req.SeqNo)

	var idToken, tokenType string
	if req.IdToken != nil {
		idToken = req.IdToken.IdToken
		tokenType = string(req.IdToken.Type)
	}

	var err error
	switch req.EventType {
	case types.TransactionEventEnumTypeStarted:
		err = t.TransactionStore.CreateTransaction(chargeStationId,
			req.TransactionInfo.TransactionId,
			idToken,
			tokenType,
			convertMeterValues(req.MeterValue),
			req.SeqNo,
			req.Offline)
	case types.TransactionEventEnumTypeUpdated:
		err = t.TransactionStore.UpdateTransaction(chargeStationId,
			req.TransactionInfo.TransactionId,
			convertMeterValues(req.MeterValue))
	case types.TransactionEventEnumTypeEnded:
		err = t.TransactionStore.EndTransaction(chargeStationId,
			req.TransactionInfo.TransactionId,
			idToken,
			tokenType,
			convertMeterValues(req.MeterValue),
			req.SeqNo)
	}

	if err != nil {
		return nil, err
	}

	response := &types.TransactionEventResponseJson{}

	if req.EventType == types.TransactionEventEnumTypeEnded {
		transaction, err := t.TransactionStore.FindTransaction(chargeStationId, req.TransactionInfo.TransactionId)
		if err != nil {
			return nil, err
		}
		cost, err := t.TariffService.CalculateCost(transaction)
		if err != nil {
			log.Printf("error calculating tariff: %v", err)
		} else {
			log.Printf("total cost: %f", cost)
			response.TotalCost = &cost
		}
	}

	return response, nil
}

func convertMeterValues(meterValues []types.MeterValueType) []services.MeterValue {
	var converted []services.MeterValue
	for _, meterValue := range meterValues {
		converted = append(converted, convertMeterValue(meterValue))
	}
	return converted
}

func convertMeterValue(meterValue types.MeterValueType) services.MeterValue {
	return services.MeterValue{
		SampledValues: convertSampledValues(meterValue.SampledValue),
		Timestamp:     meterValue.Timestamp,
	}
}

func convertSampledValues(sampledValues []types.SampledValueType) []services.SampledValue {
	var converted []services.SampledValue
	for _, sampledValue := range sampledValues {
		converted = append(converted, convertSampledValue(sampledValue))
	}
	return converted
}

func convertSampledValue(sampledValue types.SampledValueType) services.SampledValue {
	return services.SampledValue{
		Context:       (*string)(sampledValue.Context),
		Location:      (*string)(sampledValue.Location),
		Measurand:     (*string)(sampledValue.Measurand),
		Phase:         (*string)(sampledValue.Phase),
		UnitOfMeasure: convertUnitOfMeasure(sampledValue.UnitOfMeasure),
		Value:         sampledValue.Value,
	}
}

func convertUnitOfMeasure(unitOfMeasure *types.UnitOfMeasureType) *services.UnitOfMeasure {
	if unitOfMeasure == nil {
		return nil
	}

	return &services.UnitOfMeasure{
		Unit:      unitOfMeasure.Unit,
		Multipler: unitOfMeasure.Multiplier,
	}
}

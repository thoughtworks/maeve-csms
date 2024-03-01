// SPDX-License-Identifier: Apache-2.0

package ocpp201

import (
	"context"

	"github.com/thoughtworks/maeve-csms/manager/ocpp"
	types "github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"github.com/thoughtworks/maeve-csms/manager/services"
	"github.com/thoughtworks/maeve-csms/manager/store"
	"golang.org/x/exp/slog"
)

type TransactionEventHandler struct {
	Store         store.Engine
	TariffService services.TariffService
}

func (t TransactionEventHandler) HandleCall(ctx context.Context, chargeStationId string, request ocpp.Request) (ocpp.Response, error) {
	req := request.(*types.TransactionEventRequestJson)
	slog.Info("transaction event", slog.String("chargeStationId", chargeStationId),
		slog.String("transactionId", req.TransactionInfo.TransactionId),
		slog.String("eventType", string(req.EventType)),
		slog.String("triggerReason", string(req.TriggerReason)),
		slog.Int("seqNo", req.SeqNo))
	response := &types.TransactionEventResponseJson{}

	var idToken, tokenType string
	if req.IdToken != nil {
		idToken = req.IdToken.IdToken
		tokenType = string(req.IdToken.Type)

		tok, err := t.Store.LookupToken(ctx, idToken)
		if err != nil {
			return nil, err
		}
		if tok != nil {
			response.IdTokenInfo = &types.IdTokenInfoType{
				Status: types.AuthorizationStatusEnumTypeAccepted,
			}
		} else {
			response.IdTokenInfo = &types.IdTokenInfoType{
				Status: types.AuthorizationStatusEnumTypeUnknown,
			}
		}
	}

	var err error
	switch req.EventType {
	case types.TransactionEventEnumTypeStarted:
		err = t.Store.CreateTransaction(
			ctx,
			chargeStationId,
			req.TransactionInfo.TransactionId,
			idToken,
			tokenType,
			convertMeterValues(req.MeterValue),
			req.SeqNo,
			req.Offline)
	case types.TransactionEventEnumTypeUpdated:
		err = t.Store.UpdateTransaction(
			ctx,
			chargeStationId,
			req.TransactionInfo.TransactionId,
			convertMeterValues(req.MeterValue))
	case types.TransactionEventEnumTypeEnded:
		err = t.Store.EndTransaction(
			ctx,
			chargeStationId,
			req.TransactionInfo.TransactionId,
			idToken,
			tokenType,
			convertMeterValues(req.MeterValue),
			req.SeqNo)
	}

	if err != nil {
		return nil, err
	}

	if req.EventType == types.TransactionEventEnumTypeEnded {
		transaction, err := t.Store.FindTransaction(ctx, chargeStationId, req.TransactionInfo.TransactionId)
		if err != nil {
			return nil, err
		}
		cost, err := t.TariffService.CalculateCost(transaction)
		if err != nil {
			slog.Error("error calculating tariff", "err", err)
		} else {
			slog.Info("total cost", slog.Float64("cost", cost))
			response.TotalCost = &cost
		}
	}

	return response, nil
}

func convertMeterValues(meterValues []types.MeterValueType) []store.MeterValue {
	var converted []store.MeterValue
	for _, meterValue := range meterValues {
		converted = append(converted, convertMeterValue(meterValue))
	}
	return converted
}

func convertMeterValue(meterValue types.MeterValueType) store.MeterValue {
	return store.MeterValue{
		SampledValues: convertSampledValues(meterValue.SampledValue),
		Timestamp:     meterValue.Timestamp,
	}
}

func convertSampledValues(sampledValues []types.SampledValueType) []store.SampledValue {
	var converted []store.SampledValue
	for _, sampledValue := range sampledValues {
		converted = append(converted, convertSampledValue(sampledValue))
	}
	return converted
}

func convertSampledValue(sampledValue types.SampledValueType) store.SampledValue {
	return store.SampledValue{
		Context:       (*string)(sampledValue.Context),
		Location:      (*string)(sampledValue.Location),
		Measurand:     (*string)(sampledValue.Measurand),
		Phase:         (*string)(sampledValue.Phase),
		UnitOfMeasure: convertUnitOfMeasure(sampledValue.UnitOfMeasure),
		Value:         sampledValue.Value,
	}
}

func convertUnitOfMeasure(unitOfMeasure *types.UnitOfMeasureType) *store.UnitOfMeasure {
	if unitOfMeasure == nil {
		return nil
	}

	return &store.UnitOfMeasure{
		Unit:      unitOfMeasure.Unit,
		Multipler: unitOfMeasure.Multiplier,
	}
}

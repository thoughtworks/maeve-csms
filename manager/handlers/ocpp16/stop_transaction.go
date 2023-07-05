package ocpp16

import (
	"context"
	"errors"
	"github.com/twlabs/maeve-csms/manager/ocpp"
	types "github.com/twlabs/maeve-csms/manager/ocpp/ocpp16"
	"github.com/twlabs/maeve-csms/manager/services"
	"k8s.io/utils/clock"
	"log"
	"strconv"
	"time"
)

type StopTransactionHandler struct {
	Clock            clock.PassiveClock
	TokenStore       services.TokenStore
	TransactionStore services.TransactionStore
}

func (s StopTransactionHandler) HandleCall(ctx context.Context, chargeStationId string, request ocpp.Request) (response ocpp.Response, err error) {
	req := request.(*types.StopTransactionJson)

	reason := "*unknown*"
	if req.Reason != nil {
		reason = string(*req.Reason)
	}
	transactionId := ConvertToUUID(req.TransactionId)
	log.Printf("Stop transaction %s for reason %s", transactionId, reason)

	var idTagInfo *types.StopTransactionResponseJsonIdTagInfo
	if req.IdTag != nil {
		status := types.StopTransactionResponseJsonIdTagInfoStatusInvalid
		tok, err := s.TokenStore.FindToken("ISO14443", *req.IdTag)
		if err != nil {
			return nil, err
		}
		if tok != nil {
			status = types.StopTransactionResponseJsonIdTagInfoStatusAccepted
		}
		idTagInfo = &types.StopTransactionResponseJsonIdTagInfo{
			Status: status,
		}
	}

	transaction, err := s.TransactionStore.FindTransaction("cs001", transactionId)
	if err != nil {
		return nil, err
	}
	seqNo := -1
	if transaction != nil {
		seqNo = transaction.StartSeqNo + transaction.UpdatedSeqNoCount + 1
	}

	var idToken, tokenType string
	if req.IdTag != nil {
		idToken = *req.IdTag
		tokenType = "ISO14443"
	}

	meterValues, err := convertMeterValues(req.TransactionData)
	if err != nil {
		return nil, err
	}

	var previousMeterValues []services.MeterValue
	if transaction != nil {
		previousMeterValues = transaction.MeterValues
	}
	meterValues = calculateTransactionEndOutletEnergy(s.Clock, meterValues, previousMeterValues, req.MeterStop)

	err = s.TransactionStore.EndTransaction(chargeStationId, transactionId, idToken, tokenType, meterValues, seqNo)
	if err != nil {
		return nil, err
	}

	return &types.StopTransactionResponseJson{
		IdTagInfo: idTagInfo,
	}, nil
}

func calculateTransactionEndOutletEnergy(clock clock.PassiveClock, transactionValues []services.MeterValue, previousValues []services.MeterValue, meterStop int) []services.MeterValue {
	if findOutletEnergyReading(transactionValues) {
		return transactionValues
	}

	if meterStart, ok := findTransactionBeginMeterValues(previousValues); ok {
		energyUsed := meterStop - meterStart
		transactionEndContext := "Transaction.End"
		outletLocation := "Outlet"
		energyRegisteredMeasurand := "Energy.Active.Import.Register"

		transactionValues = append(transactionValues, services.MeterValue{
			SampledValues: []services.SampledValue{
				{
					Context:   &transactionEndContext,
					Location:  &outletLocation,
					Measurand: &energyRegisteredMeasurand,
					Value:     float64(energyUsed),
				},
			},
			Timestamp: clock.Now().Format(time.RFC3339),
		})
	}

	return transactionValues
}

func findOutletEnergyReading(values []services.MeterValue) bool {
	for _, value := range values {
		for _, sv := range value.SampledValues {
			if sv.Context != nil && *sv.Context == "Transaction.End" &&
				sv.Measurand != nil && *sv.Measurand == "Energy.Active.Import.Register" &&
				sv.Location != nil && *sv.Location == "Outlet" {
				return true
			}
		}
	}
	return false
}

func findTransactionBeginMeterValues(values []services.MeterValue) (int, bool) {
	for _, value := range values {
		for _, sv := range value.SampledValues {
			if sv.Context != nil && *sv.Context == "Transaction.Begin" &&
				sv.Measurand != nil && *sv.Measurand == "MeterValue" &&
				sv.Location != nil && *sv.Location == "Outlet" {
				return int(sv.Value), true
			}
		}
	}

	return 0, false
}

func convertMeterValues(meterValues []types.StopTransactionJsonTransactionDataElem) ([]services.MeterValue, error) {
	var converted []services.MeterValue
	for _, meterValue := range meterValues {
		convertedMeterValue, err := convertMeterValue(meterValue)
		if err != nil {
			return nil, err
		}
		converted = append(converted, convertedMeterValue)
	}
	return converted, nil
}

func convertMeterValue(meterValue types.StopTransactionJsonTransactionDataElem) (services.MeterValue, error) {
	sampledValues, err := convertSampledValues(meterValue.SampledValue)
	if err != nil {
		return services.MeterValue{}, err
	}
	return services.MeterValue{
		SampledValues: sampledValues,
		Timestamp:     meterValue.Timestamp,
	}, nil
}

func convertSampledValues(sampledValues []types.StopTransactionJsonTransactionDataElemSampledValueElem) ([]services.SampledValue, error) {
	var converted []services.SampledValue
	for _, sampleValue := range sampledValues {
		convertedSampleValue, err := convertSampleValue(sampleValue)
		if err != nil {
			return nil, err
		}
		converted = append(converted, convertedSampleValue)
	}
	return converted, nil
}

func convertSampleValue(sampleValue types.StopTransactionJsonTransactionDataElemSampledValueElem) (services.SampledValue, error) {
	value, err := convertValue(sampleValue.Format, sampleValue.Value)
	if err != nil {
		return services.SampledValue{}, err
	}
	return services.SampledValue{
		Context:       (*string)(sampleValue.Context),
		Location:      (*string)(sampleValue.Location),
		Measurand:     (*string)(sampleValue.Measurand),
		Phase:         (*string)(sampleValue.Phase),
		UnitOfMeasure: convertUnitOfMeasure(sampleValue.Unit),
		Value:         value,
	}, nil
}

func convertValue(format *types.StopTransactionJsonTransactionDataElemSampledValueElemFormat, value string) (float64, error) {
	if format != nil && *format != types.StopTransactionJsonTransactionDataElemSampledValueElemFormatRaw {
		return 0, errors.New("conversion from signed data not implemented")
	}

	return strconv.ParseFloat(value, 64)
}

func convertUnitOfMeasure(unit *types.StopTransactionJsonTransactionDataElemSampledValueElemUnit) *services.UnitOfMeasure {
	if unit == nil {
		return nil
	}

	return &services.UnitOfMeasure{
		Unit:      string(*unit),
		Multipler: 1,
	}
}

package ocpp16

import (
	"context"
	"github.com/google/uuid"
	"github.com/twlabs/ocpp2-broker-core/manager/ocpp"
	types "github.com/twlabs/ocpp2-broker-core/manager/ocpp/ocpp16"
	"github.com/twlabs/ocpp2-broker-core/manager/services"
	"k8s.io/utils/clock"
	"log"
	"math/rand"
	"time"
)

type StartTransactionHandler struct {
	Clock            clock.PassiveClock
	TokenStore       services.TokenStore
	TransactionStore services.TransactionStore
}

func (t StartTransactionHandler) HandleCall(ctx context.Context, chargeStationId string, request ocpp.Request) (ocpp.Response, error) {
	req := request.(*types.StartTransactionJson)

	log.Printf("Start transaction: %v", req)

	transactionId := -1
	status := types.StartTransactionResponseJsonIdTagInfoStatusInvalid
	tok, err := t.TokenStore.FindToken("", req.IdTag)
	if err != nil {
		return nil, err
	}
	if tok != nil {
		status = types.StartTransactionResponseJsonIdTagInfoStatusAccepted
		transactionId = int(rand.Int31())
	}

	contextTransactionBegin := types.MeterValuesJsonMeterValueElemSampledValueElemContextTransactionBegin
	meterValueMeasurand := "MeterValue"
	transactionUuid := ConvertToUUID(transactionId)
	err = t.TransactionStore.CreateTransaction(chargeStationId, transactionUuid, req.IdTag, "ISO14443",
		[]services.MeterValue{
			{
				Timestamp: t.Clock.Now().Format(time.RFC3339),
				SampledValues: []services.SampledValue{
					{
						Context:   (*string)(&contextTransactionBegin),
						Measurand: &meterValueMeasurand,
						UnitOfMeasure: &services.UnitOfMeasure{
							Unit:      string(types.MeterValuesJsonMeterValueElemSampledValueElemUnitWh),
							Multipler: 1,
						},
						Value: float64(req.MeterStart),
					},
				},
			},
		}, 0, false)
	if err != nil {
		return nil, err
	}

	return &types.StartTransactionResponseJson{
		IdTagInfo: types.StartTransactionResponseJsonIdTagInfo{
			Status: status,
		},
		TransactionId: transactionId,
	}, nil
}

func ConvertToUUID(transactionId int) string {
	uuidBytes := []byte{
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		(byte)(transactionId >> 24),
		(byte)(transactionId >> 16),
		(byte)(transactionId >> 8),
		(byte)(transactionId),
	}
	return uuid.Must(uuid.FromBytes(uuidBytes)).String()
}

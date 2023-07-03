package ocpp16

import (
	"context"
	"github.com/twlabs/ocpp2-broker-core/manager/ocpp"
	types "github.com/twlabs/ocpp2-broker-core/manager/ocpp/ocpp16"
	"github.com/twlabs/ocpp2-broker-core/manager/services"
)

type MeterValuesHandler struct {
	TransactionStore services.TransactionStore
}

func (m MeterValuesHandler) HandleCall(ctx context.Context, chargeStationId string, request ocpp.Request) (response ocpp.Response, err error) {
	// TODO: store in transaction store

	return &types.MeterValuesResponseJson{}, nil
}

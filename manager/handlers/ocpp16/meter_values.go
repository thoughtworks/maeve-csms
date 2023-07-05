package ocpp16

import (
	"context"
	"github.com/twlabs/maeve-csms/manager/ocpp"
	types "github.com/twlabs/maeve-csms/manager/ocpp/ocpp16"
	"github.com/twlabs/maeve-csms/manager/services"
)

type MeterValuesHandler struct {
	TransactionStore services.TransactionStore
}

func (m MeterValuesHandler) HandleCall(ctx context.Context, chargeStationId string, request ocpp.Request) (response ocpp.Response, err error) {
	// TODO: store in transaction store

	return &types.MeterValuesResponseJson{}, nil
}

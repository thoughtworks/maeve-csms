package ocpp16

import (
	"context"
	"github.com/twlabs/ocpp2-broker-core/manager/ocpp"
	types "github.com/twlabs/ocpp2-broker-core/manager/ocpp/ocpp16"
	"log"
)

func StatusNotificationHandler(ctx context.Context, chargeStationId string, request ocpp.Request) (ocpp.Response, error) {
	req := request.(*types.StatusNotificationJson)
	log.Printf("Charge station %s, connection %d status: %s(%s)", chargeStationId, req.ConnectorId, req.Status, req.ErrorCode)
	return &types.StatusNotificationResponseJson{}, nil
}

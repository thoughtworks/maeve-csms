package ocpp16

import (
	"context"
	"github.com/twlabs/maeve-csms/manager/ocpp"
	types "github.com/twlabs/maeve-csms/manager/ocpp/ocpp16"
	"log"
)

func StatusNotificationHandler(ctx context.Context, chargeStationId string, request ocpp.Request) (ocpp.Response, error) {
	req := request.(*types.StatusNotificationJson)
	log.Printf("Charge station %s, connection %d status: %s(%s)", chargeStationId, req.ConnectorId, req.Status, req.ErrorCode)
	return &types.StatusNotificationResponseJson{}, nil
}

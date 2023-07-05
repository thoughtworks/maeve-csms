package ocpp201

import (
	"context"
	"github.com/twlabs/maeve-csms/manager/ocpp"
	types "github.com/twlabs/maeve-csms/manager/ocpp/ocpp201"
	"log"
)

func StatusNotificationHandler(ctx context.Context, chargeStationId string, request ocpp.Request) (ocpp.Response, error) {
	req := request.(*types.StatusNotificationRequestJson)
	log.Printf("Charge station %s, EVSE %d, connection %d status: %s", chargeStationId, req.EvseId, req.ConnectorId, req.ConnectorStatus)
	return &types.StatusNotificationResponseJson{}, nil
}

package ocpp201

import (
	"context"
	"github.com/twlabs/maeve-csms/manager/ocpp"
	types "github.com/twlabs/maeve-csms/manager/ocpp/ocpp201"
	"k8s.io/utils/clock"
	"log"
	"time"
)

type BootNotificationHandler struct {
	Clock             clock.PassiveClock
	HeartbeatInterval int
}

func (b BootNotificationHandler) HandleCall(ctx context.Context, chargeStationId string, request ocpp.Request) (ocpp.Response, error) {
	req := request.(*types.BootNotificationRequestJson)
	var serialNumber string
	if req.ChargingStation.SerialNumber != nil {
		serialNumber = *req.ChargingStation.SerialNumber
	} else {
		serialNumber = "*unknown*"
	}
	log.Printf("Charge station %s with serial number %s booting for reason %s", chargeStationId, serialNumber, req.Reason)
	return &types.BootNotificationResponseJson{
		CurrentTime: b.Clock.Now().Format(time.RFC3339),
		Interval:    b.HeartbeatInterval,
		Status:      types.RegistrationStatusEnumTypeAccepted,
	}, nil
}

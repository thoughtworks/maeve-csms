package ocpp16

import (
	"context"
	"github.com/twlabs/maeve-csms/manager/ocpp"
	types "github.com/twlabs/maeve-csms/manager/ocpp/ocpp16"
	"k8s.io/utils/clock"
	"log"
	"time"
)

type BootNotificationHandler struct {
	Clock             clock.PassiveClock
	HeartbeatInterval int
}

func (b BootNotificationHandler) HandleCall(ctx context.Context, chargeStationId string, request ocpp.Request) (ocpp.Response, error) {
	req := request.(*types.BootNotificationJson)

	var serialNumber string
	if req.ChargePointSerialNumber != nil {
		serialNumber = *req.ChargePointSerialNumber
	} else {
		serialNumber = "*unknown*"
	}
	log.Printf("Charge station %s with serial number %s booting", chargeStationId, serialNumber)
	return &types.BootNotificationResponseJson{
		CurrentTime: b.Clock.Now().Format(time.RFC3339),
		Interval:    b.HeartbeatInterval,
		Status:      types.BootNotificationResponseJsonStatusAccepted,
	}, nil
}

package ocpp16

import (
	"context"
	"github.com/thoughtworks/maeve-csms/manager/ocpp"
	types "github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp16"
	"k8s.io/utils/clock"
	"log"
	"time"
)

type HeartbeatHandler struct {
	Clock clock.PassiveClock
}

func (h HeartbeatHandler) HandleCall(ctx context.Context, chargeStationId string, request ocpp.Request) (ocpp.Response, error) {
	log.Printf("Charge station %s heartbeat", chargeStationId)
	return &types.HeartbeatResponseJson{
		CurrentTime: h.Clock.Now().Format(time.RFC3339),
	}, nil
}

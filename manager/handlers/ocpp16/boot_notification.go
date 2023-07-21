// SPDX-License-Identifier: Apache-2.0

package ocpp16

import (
	"context"
	"time"

	"github.com/thoughtworks/maeve-csms/manager/ocpp"
	types "github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp16"
	"golang.org/x/exp/slog"
	"k8s.io/utils/clock"
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
	slog.Info("booting", slog.String("chargeStationId", chargeStationId),
		slog.String("serialNumber", serialNumber))
	return &types.BootNotificationResponseJson{
		CurrentTime: b.Clock.Now().Format(time.RFC3339),
		Interval:    b.HeartbeatInterval,
		Status:      types.BootNotificationResponseJsonStatusAccepted,
	}, nil
}

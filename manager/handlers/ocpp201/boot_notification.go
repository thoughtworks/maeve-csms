// SPDX-License-Identifier: Apache-2.0

package ocpp201

import (
	"context"
	"time"

	"github.com/thoughtworks/maeve-csms/manager/ocpp"
	types "github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"golang.org/x/exp/slog"
	"k8s.io/utils/clock"
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
	slog.Info("booting", slog.String("chargeStationId", chargeStationId),
		slog.String("serialNumber", serialNumber),
		slog.String("reason", string(req.Reason)))
	return &types.BootNotificationResponseJson{
		CurrentTime: b.Clock.Now().Format(time.RFC3339),
		Interval:    b.HeartbeatInterval,
		Status:      types.RegistrationStatusEnumTypeAccepted,
	}, nil
}

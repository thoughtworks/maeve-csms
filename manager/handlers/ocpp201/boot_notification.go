// SPDX-License-Identifier: Apache-2.0

package ocpp201

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"time"

	"github.com/thoughtworks/maeve-csms/manager/ocpp"
	types "github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"k8s.io/utils/clock"
)

type BootNotificationHandler struct {
	Clock             clock.PassiveClock
	HeartbeatInterval int
}

func (b BootNotificationHandler) HandleCall(ctx context.Context, chargeStationId string, request ocpp.Request) (ocpp.Response, error) {
	span := trace.SpanFromContext(ctx)

	req := request.(*types.BootNotificationRequestJson)

	span.SetAttributes(
		attribute.String("request.status", string(types.RegistrationStatusEnumTypeAccepted)),
		attribute.String("boot.reason", string(req.Reason)),
		attribute.String("boot.vendor", req.ChargingStation.VendorName),
		attribute.String("boot.model", req.ChargingStation.Model))

	if req.ChargingStation.SerialNumber != nil {
		span.SetAttributes(attribute.String("boot.serial", *req.ChargingStation.SerialNumber))
	}
	if req.ChargingStation.FirmwareVersion != nil {
		span.SetAttributes(attribute.String("boot.firmware", *req.ChargingStation.FirmwareVersion))
	}

	return &types.BootNotificationResponseJson{
		CurrentTime: b.Clock.Now().Format(time.RFC3339),
		Interval:    b.HeartbeatInterval,
		Status:      types.RegistrationStatusEnumTypeAccepted,
	}, nil
}

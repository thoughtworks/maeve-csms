// SPDX-License-Identifier: Apache-2.0

package ocpp16

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
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
	span := trace.SpanFromContext(ctx)

	req := request.(*types.BootNotificationJson)

	var serialNumber string
	if req.ChargePointSerialNumber != nil {
		serialNumber = *req.ChargePointSerialNumber
	} else {
		serialNumber = "*unknown*"
	}

	span.SetAttributes(
		attribute.String("request.status", string(types.BootNotificationResponseJsonStatusAccepted)),
		attribute.String("boot.vendor", req.ChargePointVendor),
		attribute.String("boot.model", req.ChargePointModel))

	if req.ChargePointSerialNumber != nil {
		span.SetAttributes(attribute.String("boot.serial", *req.ChargePointSerialNumber))
	}
	if req.FirmwareVersion != nil {
		span.SetAttributes(attribute.String("boot.firmware", *req.FirmwareVersion))
	}

	slog.Info("booting", slog.String("chargeStationId", chargeStationId),
		slog.String("serialNumber", serialNumber))
	return &types.BootNotificationResponseJson{
		CurrentTime: b.Clock.Now().Format(time.RFC3339),
		Interval:    b.HeartbeatInterval,
		Status:      types.BootNotificationResponseJsonStatusAccepted,
	}, nil
}

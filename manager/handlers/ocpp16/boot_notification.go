// SPDX-License-Identifier: Apache-2.0

package ocpp16

import (
	"context"
	"github.com/thoughtworks/maeve-csms/manager/store"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"time"

	"github.com/thoughtworks/maeve-csms/manager/ocpp"
	types "github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp16"
	"k8s.io/utils/clock"
)

type BootNotificationHandler struct {
	Clock               clock.PassiveClock
	RuntimeDetailsStore store.ChargeStationRuntimeDetailsStore
	SettingsStore       store.ChargeStationSettingsStore
	HeartbeatInterval   int
}

func (b BootNotificationHandler) HandleCall(ctx context.Context, chargeStationId string, request ocpp.Request) (ocpp.Response, error) {
	span := trace.SpanFromContext(ctx)

	req := request.(*types.BootNotificationJson)

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

	err := b.RuntimeDetailsStore.SetChargeStationRuntimeDetails(ctx, chargeStationId, &store.ChargeStationRuntimeDetails{
		OcppVersion: "1.6",
	})
	if err != nil {
		return nil, err
	}

	// remove any reboot required settings
	settings, err := b.SettingsStore.LookupChargeStationSettings(ctx, chargeStationId)
	if err != nil {
		return nil, err
	}

	if settings != nil && settings.Settings != nil {
		updated := false
		for _, setting := range settings.Settings {
			if setting.Status == store.ChargeStationSettingStatusRebootRequired {
				setting.Status = store.ChargeStationSettingStatusAccepted
				updated = true
			}
		}

		if updated {
			err = b.SettingsStore.UpdateChargeStationSettings(ctx, chargeStationId, settings)
			if err != nil {
				return nil, err
			}
		}
	}

	return &types.BootNotificationResponseJson{
		CurrentTime: b.Clock.Now().Format(time.RFC3339),
		Interval:    b.HeartbeatInterval,
		Status:      types.BootNotificationResponseJsonStatusAccepted,
	}, nil
}

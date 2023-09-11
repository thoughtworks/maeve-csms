// SPDX-License-Identifier: Apache-2.0

package ocpp16

import (
	"context"
	"fmt"
	"github.com/thoughtworks/maeve-csms/manager/handlers"
	"github.com/thoughtworks/maeve-csms/manager/ocpp"
	"github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp16"
	"github.com/thoughtworks/maeve-csms/manager/store"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type ChangeConfigurationResultHandler struct {
	SettingsStore store.ChargeStationSettingsStore
	CallMaker     handlers.CallMaker
}

func (c ChangeConfigurationResultHandler) HandleCallResult(ctx context.Context, chargeStationId string, request ocpp.Request, response ocpp.Response, state any) error {
	req := request.(*ocpp16.ChangeConfigurationJson)
	resp := response.(*ocpp16.ChangeConfigurationResponseJson)

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("setting.key", req.Key),
		attribute.String("setting.value", req.Value),
		attribute.String("setting.status", string(resp.Status)))

	err := c.SettingsStore.UpdateChargeStationSettings(ctx, chargeStationId, &store.ChargeStationSettings{
		ChargeStationId: chargeStationId,
		Settings: map[string]*store.ChargeStationSetting{
			req.Key: {
				Value:  req.Value,
				Status: store.ChargeStationSettingStatus(resp.Status),
			},
		},
	})
	if err != nil {
		return fmt.Errorf("update charge station settings: %w", err)
	}

	settings, err := c.SettingsStore.LookupChargeStationSettings(ctx, chargeStationId)
	if err != nil {
		return fmt.Errorf("lookup charge station settings: %w", err)
	}

	// check if all settings are done and if so reboot the charge station if necessary
	allDone := true
	rebootRequired := false
	for _, setting := range settings.Settings {
		if setting.Status == store.ChargeStationSettingStatusRebootRequired {
			rebootRequired = true
		}
		if setting.Status == store.ChargeStationSettingStatusPending {
			allDone = false
			break
		}
	}
	if allDone && rebootRequired {
		err = c.CallMaker.Send(ctx, chargeStationId, &ocpp16.TriggerMessageJson{
			RequestedMessage: ocpp16.TriggerMessageJsonRequestedMessageBootNotification,
		})
	}

	return err
}

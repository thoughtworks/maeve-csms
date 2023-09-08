// SPDX-License-Identifier: Apache-2.0

package ocpp16

import (
	"context"
	"github.com/thoughtworks/maeve-csms/manager/ocpp"
	"github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp16"
	"github.com/thoughtworks/maeve-csms/manager/store"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type ChangeConfigurationResultHandler struct {
	SettingsStore store.ChargeStationSettingsStore
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

	return err
}

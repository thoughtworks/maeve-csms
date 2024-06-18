// SPDX-License-Identifier: Apache-2.0

package ocpp201

import (
	"context"

	"github.com/thoughtworks/maeve-csms/manager/ocpp"
	"github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"github.com/thoughtworks/maeve-csms/manager/store"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type TriggerMessageResultHandler struct {
	Store store.Engine
}

func (i TriggerMessageResultHandler) HandleCallResult(ctx context.Context, chargeStationId string, request ocpp.Request, response ocpp.Response, state any) error {
	req := request.(*ocpp201.TriggerMessageRequestJson)

	status := ocpp201.TriggerMessageStatusEnumTypeNotImplemented

	span := trace.SpanFromContext(ctx)
	if response != nil {
		resp := response.(*ocpp201.TriggerMessageResponseJson)

		span.SetAttributes(
			attribute.String("trigger_message.trigger", string(req.RequestedMessage)),
			attribute.String("trigger_message.status", string(resp.Status)))

		if resp.Status == ocpp201.TriggerMessageStatusEnumTypeAccepted {
			status = resp.Status
		}
	} else {
		span.SetAttributes(
			attribute.String("trigger_message.trigger", string(req.RequestedMessage)))

		status = ocpp201.TriggerMessageStatusEnumTypeAccepted
	}

	if status == ocpp201.TriggerMessageStatusEnumTypeAccepted {
		return i.Store.DeleteChargeStationTriggerMessage(ctx, chargeStationId)
	} else {
		err := i.Store.SetChargeStationTriggerMessage(ctx, chargeStationId, &store.ChargeStationTriggerMessage{
			TriggerMessage: store.TriggerMessage(req.RequestedMessage),
			TriggerStatus:  store.TriggerStatus(status),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

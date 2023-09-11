// SPDX-License-Identifier: Apache-2.0

package ocpp16

import (
	"context"
	"github.com/thoughtworks/maeve-csms/manager/ocpp"
	"github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp16"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type TriggerMessageResultHandler struct{}

func (c TriggerMessageResultHandler) HandleCallResult(ctx context.Context, chargeStationId string, request ocpp.Request, response ocpp.Response, state any) error {
	req := request.(*ocpp16.TriggerMessageJson)
	resp := response.(*ocpp16.TriggerMessageResponseJson)

	span := trace.SpanFromContext(ctx)

	span.SetAttributes(
		attribute.String("trigger.requested_message", string(req.RequestedMessage)),
		attribute.String("trigger.status", string(resp.Status)))

	if req.ConnectorId != nil {
		span.SetAttributes(attribute.Int("trigger.connector_id", *req.ConnectorId))
	}

	return nil
}

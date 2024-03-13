// SPDX-License-Identifier: Apache-2.0

package ocpp201

import (
	"context"
	"github.com/thoughtworks/maeve-csms/manager/ocpp"
	types "github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type ChangeAvailabilityResultHandler struct{}

func (h ChangeAvailabilityResultHandler) HandleCallResult(ctx context.Context, chargeStationId string, request ocpp.Request, response ocpp.Response, state any) error {
	req := request.(*types.ChangeAvailabilityRequestJson)
	resp := response.(*types.ChangeAvailabilityResponseJson)

	span := trace.SpanFromContext(ctx)

	span.SetAttributes(
		attribute.String("change_availability.operational_status", string(req.OperationalStatus)),
		attribute.String("change_availability.status", string(resp.Status)))

	if req.Evse != nil {
		span.SetAttributes(
			attribute.Int("change_availability.evse_id", req.Evse.Id))

		if req.Evse.ConnectorId != nil {
			span.SetAttributes(
				attribute.Int("change_availability.connector_id", *req.Evse.ConnectorId))
		}
	}

	return nil
}

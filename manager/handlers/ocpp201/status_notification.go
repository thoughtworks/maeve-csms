// SPDX-License-Identifier: Apache-2.0

package ocpp201

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/thoughtworks/maeve-csms/manager/ocpp"
	types "github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
)

func StatusNotificationHandler(ctx context.Context, chargeStationId string, request ocpp.Request) (ocpp.Response, error) {
	span := trace.SpanFromContext(ctx)

	req := request.(*types.StatusNotificationRequestJson)

	span.SetAttributes(
		attribute.Int("status.evse_id", req.EvseId),
		attribute.Int("status.connector_id", req.ConnectorId),
		attribute.String("status.connector_status", string(req.ConnectorStatus)))

	return &types.StatusNotificationResponseJson{}, nil
}

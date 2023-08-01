// SPDX-License-Identifier: Apache-2.0

package ocpp16

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/thoughtworks/maeve-csms/manager/ocpp"
	types "github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp16"
)

func StatusNotificationHandler(ctx context.Context, chargeStationId string, request ocpp.Request) (ocpp.Response, error) {
	span := trace.SpanFromContext(ctx)

	req := request.(*types.StatusNotificationJson)

	span.SetAttributes(
		attribute.Int("status.connector_id", req.ConnectorId),
		attribute.String("status.connector_status", string(req.Status)))

	return &types.StatusNotificationResponseJson{}, nil
}

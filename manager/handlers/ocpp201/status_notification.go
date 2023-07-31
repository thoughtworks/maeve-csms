// SPDX-License-Identifier: Apache-2.0

package ocpp201

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/thoughtworks/maeve-csms/manager/ocpp"
	types "github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"golang.org/x/exp/slog"
)

func StatusNotificationHandler(ctx context.Context, chargeStationId string, request ocpp.Request) (ocpp.Response, error) {
	span := trace.SpanFromContext(ctx)

	req := request.(*types.StatusNotificationRequestJson)
	slog.Info("status notification", slog.String("chargeStationId", chargeStationId),
		slog.Int("evseId", req.EvseId), slog.Int("connectorId", req.ConnectorId), slog.Any("connectorStatus", req.ConnectorStatus))

	span.SetAttributes(
		attribute.Int("status.evse_id", req.EvseId),
		attribute.Int("status.connector_id", req.ConnectorId),
		attribute.String("status.connector_status", string(req.ConnectorStatus)))

	return &types.StatusNotificationResponseJson{}, nil
}

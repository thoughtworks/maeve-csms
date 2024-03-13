// SPDX-License-Identifier: Apache-2.0

package ocpp201

import (
	"context"
	"github.com/thoughtworks/maeve-csms/manager/ocpp"
	"github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type UnlockConnectorResultHandler struct{}

func (h UnlockConnectorResultHandler) HandleCallResult(ctx context.Context, chargeStationId string, request ocpp.Request, response ocpp.Response, state any) error {
	req := request.(*ocpp201.UnlockConnectorRequestJson)
	resp := response.(*ocpp201.UnlockConnectorResponseJson)

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.Int("unlock_connector.evse_id", req.EvseId),
		attribute.Int("unlock_connector.connector_id", req.ConnectorId),
		attribute.String("unlock_connector.status", string(resp.Status)))

	return nil
}

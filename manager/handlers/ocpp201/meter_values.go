// SPDX-License-Identifier: Apache-2.0

package ocpp201

import (
	"context"
	"github.com/thoughtworks/maeve-csms/manager/ocpp"
	"github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type MeterValuesHandler struct{}

func (h MeterValuesHandler) HandleCall(ctx context.Context, chargeStationId string, request ocpp.Request) (response ocpp.Response, err error) {
	req := request.(*ocpp201.MeterValuesRequestJson)

	span := trace.SpanFromContext(ctx)

	span.SetAttributes(attribute.Int("meter_values.evse_id", req.EvseId))

	return &ocpp201.MeterValuesResponseJson{}, nil
}

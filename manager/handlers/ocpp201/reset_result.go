// SPDX-License-Identifier: Apache-2.0

package ocpp201

import (
	"context"
	"github.com/thoughtworks/maeve-csms/manager/ocpp"
	types "github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type ResetResultHandler struct{}

func (h ResetResultHandler) HandleCallResult(ctx context.Context, chargeStationId string, request ocpp.Request, response ocpp.Response, state any) error {
	req := request.(*types.ResetRequestJson)
	resp := response.(*types.ResetResponseJson)

	span := trace.SpanFromContext(ctx)

	span.SetAttributes(
		attribute.String("reset.type", string(req.Type)),
		attribute.String("reset.status", string(resp.Status)))

	return nil
}

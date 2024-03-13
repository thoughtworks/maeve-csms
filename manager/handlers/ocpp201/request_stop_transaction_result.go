// SPDX-License-Identifier: Apache-2.0

package ocpp201

import (
	"context"
	"github.com/thoughtworks/maeve-csms/manager/ocpp"
	types "github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type RequestStopTransactionResultHandler struct{}

func (h RequestStopTransactionResultHandler) HandleCallResult(ctx context.Context, chargeStationId string, request ocpp.Request, response ocpp.Response, state any) error {
	req := request.(*types.RequestStopTransactionRequestJson)
	resp := response.(*types.RequestStopTransactionResponseJson)

	span := trace.SpanFromContext(ctx)

	span.SetAttributes(
		attribute.String("request_stop.transaction_id", req.TransactionId),
		attribute.String("request_stop.status", string(resp.Status)))

	return nil
}

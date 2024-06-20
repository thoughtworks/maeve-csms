// SPDX-License-Identifier: Apache-2.0

package ocpp201

import (
	"context"
	"github.com/thoughtworks/maeve-csms/manager/ocpp"
	types "github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type RequestStartTransactionResultHandler struct{}

func (h RequestStartTransactionResultHandler) HandleCallResult(ctx context.Context, chargeStationId string, request ocpp.Request, response ocpp.Response, state any) error {
	req := request.(*types.RequestStartTransactionRequestJson)
	resp := response.(*types.RequestStartTransactionResponseJson)

	span := trace.SpanFromContext(ctx)

	span.SetAttributes(
		attribute.Int("request_start.remote_start_id", req.RemoteStartId),
		attribute.String("request_start.status", string(resp.Status)))

	if resp.TransactionId != nil {
		span.SetAttributes(
			attribute.String("request_start.transaction_id", *resp.TransactionId))
	}

	return nil
}

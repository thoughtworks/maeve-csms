// SPDX-License-Identifier: Apache-2.0

package ocpp201

import (
	"context"
	"github.com/thoughtworks/maeve-csms/manager/ocpp"
	types "github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type GetTransactionStatusResultHandler struct{}

func (h GetTransactionStatusResultHandler) HandleCallResult(ctx context.Context, chargeStationId string, request ocpp.Request, response ocpp.Response, state any) error {
	req := request.(*types.GetTransactionStatusRequestJson)
	resp := response.(*types.GetTransactionStatusResponseJson)

	span := trace.SpanFromContext(ctx)

	if req.TransactionId != nil {
		span.SetAttributes(
			attribute.String("get_transaction_status.transaction_id", *req.TransactionId))
	}

	span.SetAttributes(
		attribute.Bool("get_transaction_status.messages_in_queue", resp.MessagesInQueue))
	if resp.OngoingIndicator != nil {
		span.SetAttributes(
			attribute.Bool("get_transaction_status.ongoing", *resp.OngoingIndicator))
	}

	return nil
}

// SPDX-License-Identifier: Apache-2.0

package ocpp201

import (
	"context"

	"github.com/thoughtworks/maeve-csms/manager/ocpp"
	types "github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"golang.org/x/exp/slog"
)

type SetChargingProfileResultHandler struct{}

func (h SetChargingProfileResultHandler) HandleCallResult(ctx context.Context, chargeStationId string, request ocpp.Request, response ocpp.Response, state any) error {
	// req := request.(*types.RequestStartTransactionRequestJson)
	// resp := response.(*types.RequestStartTransactionResponseJson)
	resp := response.(*types.SetChargingProfileResponseJson)

	slog.Info("[TEST] in scp_result.go, got response:", resp)

	// span := trace.SpanFromContext(ctx)

	// span.SetAttributes(
	// 	attribute.Int("request_start.remote_start_id", req.RemoteStartId),
	// 	attribute.String("request_start.status", string(resp.Status)))

	// if resp.TransactionId != nil {
	// 	span.SetAttributes(
	// 		attribute.String("request_start.transaction_id", *resp.TransactionId))
	// }

	// Do something here ^^^ above is template from request_start_transaction_result.go

	return nil
}

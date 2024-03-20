// SPDX-License-Identifier: Apache-2.0

package ocpp201

import (
	"context"
	"github.com/thoughtworks/maeve-csms/manager/ocpp"
	types "github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type SendLocalListResultHandler struct{}

func (h SendLocalListResultHandler) HandleCallResult(ctx context.Context, chargeStationId string, request ocpp.Request, response ocpp.Response, state any) error {
	req := request.(*types.SendLocalListRequestJson)
	resp := response.(*types.SendLocalListResponseJson)

	span := trace.SpanFromContext(ctx)

	span.SetAttributes(
		attribute.String("send_local_list.update_type", string(req.UpdateType)),
		attribute.Int("send_local_list.version_number", req.VersionNumber),
		attribute.String("send_local_list.status", string(resp.Status)))

	return nil
}

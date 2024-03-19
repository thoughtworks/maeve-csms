// SPDX-License-Identifier: Apache-2.0

package ocpp201

import (
	"context"
	"github.com/thoughtworks/maeve-csms/manager/ocpp"
	types "github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type SetNetworkProfileResultHandler struct{}

func (h SetNetworkProfileResultHandler) HandleCallResult(ctx context.Context, chargeStationId string, request ocpp.Request, response ocpp.Response, state any) error {
	req := request.(*types.SetNetworkProfileRequestJson)
	resp := response.(*types.SetNetworkProfileResponseJson)

	span := trace.SpanFromContext(ctx)

	span.SetAttributes(
		attribute.Int("set_network_profile.config_slot", req.ConfigurationSlot),
		attribute.String("set_network_profile.status", string(resp.Status)))

	return nil
}

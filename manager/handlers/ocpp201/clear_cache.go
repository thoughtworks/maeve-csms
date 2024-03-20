// SPDX-License-Identifier: Apache-2.0

package ocpp201

import (
	"context"
	"github.com/thoughtworks/maeve-csms/manager/ocpp"
	types "github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type ClearCacheResultHandler struct{}

func (h ClearCacheResultHandler) HandleCallResult(ctx context.Context, chargeStationId string, request ocpp.Request, response ocpp.Response, state any) error {
	resp := response.(*types.ClearCacheResponseJson)

	span := trace.SpanFromContext(ctx)

	span.SetAttributes(
		attribute.String("clear_cache.status", string(resp.Status)))

	return nil
}

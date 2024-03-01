// SPDX-License-Identifier: Apache-2.0

package ocpp201

import (
	"context"
	"fmt"
	"github.com/thoughtworks/maeve-csms/manager/ocpp"
	"github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"github.com/thoughtworks/maeve-csms/manager/store"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/exp/slog"
)

type SetVariablesResultHandler struct {
	Store store.Engine
}

func (i SetVariablesResultHandler) HandleCallResult(ctx context.Context, chargeStationId string, request ocpp.Request, response ocpp.Response, state any) error {
	span := trace.SpanFromContext(ctx)
	if response != nil {
		resp := response.(*ocpp201.SetVariablesResponseJson)

		for _, variable := range resp.SetVariableResult {
			span.SetAttributes(
				attribute.String(fmt.Sprintf("set_variables.%s_%s.result", variable.Component.Name, variable.Variable.Name),
					string(variable.AttributeStatus)))
		}

		err := i.Store.DeleteChargeStationSettings(ctx, chargeStationId)
		if err != nil {
			slog.Error("failed to delete charge station settings", "err", err)
			span.AddEvent("failed to delete charge station settings", trace.WithAttributes(attribute.String("err", err.Error())))
		}
	}

	return nil
}

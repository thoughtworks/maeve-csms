// SPDX-License-Identifier: Apache-2.0

package ocpp201

import (
	"context"
	"github.com/thoughtworks/maeve-csms/manager/ocpp"
	"github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type SecurityEventNotificationHandler struct{}

func (s SecurityEventNotificationHandler) HandleCall(ctx context.Context, chargeStationId string, request ocpp.Request) (response ocpp.Response, err error) {
	req := request.(*ocpp201.SecurityEventNotificationRequestJson)

	span := trace.SpanFromContext(ctx)

	span.SetAttributes(attribute.String("security_event.timestamp", req.Timestamp), attribute.String("security_event.type", req.Type))
	if req.TechInfo != nil {
		span.SetAttributes(attribute.String("security_event.tech_info", *req.TechInfo))
	}

	return &ocpp201.SecurityEventNotificationResponseJson{}, nil
}

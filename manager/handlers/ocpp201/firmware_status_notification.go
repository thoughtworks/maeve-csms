// SPDX-License-Identifier: Apache-2.0

package ocpp201

import (
	"context"
	"github.com/thoughtworks/maeve-csms/manager/ocpp"
	"github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type FirmwareStatusNotificationHandler struct{}

func (h FirmwareStatusNotificationHandler) HandleCall(ctx context.Context, chargeStationId string, request ocpp.Request) (response ocpp.Response, err error) {
	req := request.(*ocpp201.FirmwareStatusNotificationRequestJson)

	span := trace.SpanFromContext(ctx)

	span.SetAttributes(attribute.String("firmware_status.status", string(req.Status)))
	if req.RequestId != nil {
		span.SetAttributes(attribute.Int("firmware_status.request_id", *req.RequestId))
	}

	return &ocpp201.FirmwareStatusNotificationResponseJson{}, nil
}

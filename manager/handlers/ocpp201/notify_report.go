// SPDX-License-Identifier: Apache-2.0

package ocpp201

import (
	"context"
	"github.com/thoughtworks/maeve-csms/manager/ocpp"
	"github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type NotifyReportHandler struct{}

func (h NotifyReportHandler) HandleCall(ctx context.Context, chargeStationId string, request ocpp.Request) (response ocpp.Response, err error) {
	req := request.(*ocpp201.NotifyReportRequestJson)

	span := trace.SpanFromContext(ctx)

	span.SetAttributes(
		attribute.String("notify_report.generated_at", req.GeneratedAt),
		attribute.Int("notify_report.request_id", req.RequestId),
		attribute.Int("notify_report.seq_no", req.SeqNo),
		attribute.Bool("notify_report.tbc", req.Tbc))

	return &ocpp201.NotifyReportResponseJson{}, nil
}

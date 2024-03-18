// SPDX-License-Identifier: Apache-2.0

package ocpp201

import (
	"context"
	"github.com/thoughtworks/maeve-csms/manager/ocpp"
	types "github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type GetBaseReportResultHandler struct{}

func (h GetBaseReportResultHandler) HandleCallResult(ctx context.Context, chargeStationId string, request ocpp.Request, response ocpp.Response, state any) error {
	req := request.(*types.GetBaseReportRequestJson)
	resp := response.(*types.GetBaseReportResponseJson)

	span := trace.SpanFromContext(ctx)

	span.SetAttributes(
		attribute.Int("get_base_report.request_id", req.RequestId),
		attribute.String("get_base_report.report_base", string(req.ReportBase)),
		attribute.String("get_base_report.status", string(resp.Status)))

	return nil
}

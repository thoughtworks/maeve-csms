// SPDX-License-Identifier: Apache-2.0

package ocpp201

import (
	"context"
	"github.com/thoughtworks/maeve-csms/manager/ocpp"
	types "github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type DeleteCertificateResultHandler struct{}

func (h DeleteCertificateResultHandler) HandleCallResult(ctx context.Context, chargeStationId string, request ocpp.Request, response ocpp.Response, state any) error {
	req := request.(*types.DeleteCertificateRequestJson)
	resp := response.(*types.DeleteCertificateResponseJson)

	span := trace.SpanFromContext(ctx)

	span.SetAttributes(
		attribute.String("delete_certificate.serial_number", req.CertificateHashData.SerialNumber),
		attribute.String("delete_certificate.status", string(resp.Status)))

	return nil
}

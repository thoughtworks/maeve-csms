// SPDX-License-Identifier: Apache-2.0

package ocpp201

import (
	"context"
	"github.com/thoughtworks/maeve-csms/manager/ocpp"
	types "github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"strings"
)

type GetInstalledCertificateIdsResultHandler struct{}

func (h GetInstalledCertificateIdsResultHandler) HandleCallResult(ctx context.Context, chargeStationId string, request ocpp.Request, response ocpp.Response, state any) error {
	req := request.(*types.GetInstalledCertificateIdsRequestJson)
	resp := response.(*types.GetInstalledCertificateIdsResponseJson)

	span := trace.SpanFromContext(ctx)

	var certTypes []string
	if req.CertificateType != nil {
		for _, ct := range req.CertificateType {
			certTypes = append(certTypes, string(ct))
		}
	}

	span.SetAttributes(
		attribute.String("get_installed_certificate.types", strings.Join(certTypes, ",")),
		attribute.String("get_installed_certificate.status", string(resp.Status)))

	return nil
}

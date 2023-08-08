// SPDX-License-Identifier: Apache-2.0

package ocpp201

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/thoughtworks/maeve-csms/manager/ocpp"
	types "github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"github.com/thoughtworks/maeve-csms/manager/services"
)

type Get15118EvCertificateHandler struct {
	ContractCertificateProvider services.ContractCertificateProvider
}

func (g Get15118EvCertificateHandler) HandleCall(ctx context.Context, _ string, request ocpp.Request) (ocpp.Response, error) {
	span := trace.SpanFromContext(ctx)

	req := request.(*types.Get15118EVCertificateRequestJson)

	status := types.Iso15118EVCertificateStatusEnumTypeFailed
	response := types.Get15118EVCertificateResponseJson{
		Status: status,
	}
	if g.ContractCertificateProvider != nil {
		res, err := g.ContractCertificateProvider.ProvideCertificate(ctx, req.ExiRequest)

		if err != nil {
			span.SetAttributes(attribute.String("get_ev_cert.error", err.Error()))
		} else {
			response = types.Get15118EVCertificateResponseJson{
				Status:      res.Status,
				ExiResponse: res.CertificateInstallationRes,
			}
		}
	}

	span.SetAttributes(attribute.String("request.status", string(status)))

	return &response, nil
}

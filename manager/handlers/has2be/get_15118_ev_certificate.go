package has2be

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/thoughtworks/maeve-csms/manager/ocpp"
	types "github.com/thoughtworks/maeve-csms/manager/ocpp/has2be"
	"github.com/thoughtworks/maeve-csms/manager/services"
)

type Get15118EvCertificateHandler struct {
	EvCertificateProvider services.EvCertificateProvider
}

func (g Get15118EvCertificateHandler) HandleCall(ctx context.Context, _ string, request ocpp.Request) (ocpp.Response, error) {
	span := trace.SpanFromContext(ctx)

	req := request.(*types.Get15118EVCertificateRequestJson)

	status := types.Iso15118EVCertificateStatusEnumTypeFailed
	response := types.Get15118EVCertificateResponseJson{
		Status: status,
	}
	if g.EvCertificateProvider != nil {
		res, err := g.EvCertificateProvider.ProvideCertificate(ctx, req.ExiRequest)

		if err != nil {
			span.SetAttributes(attribute.String("get_ev_cert.error", err.Error()))
		} else {
			response = types.Get15118EVCertificateResponseJson{
				Status:      types.Iso15118EVCertificateStatusEnumType(res.Status),
				ExiResponse: res.CertificateInstallationRes,
			}
		}
	}

	span.SetAttributes(attribute.String("request.status", string(status)))

	return &response, nil
}

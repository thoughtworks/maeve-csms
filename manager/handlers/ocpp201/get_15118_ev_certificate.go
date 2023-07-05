package ocpp201

import (
	"context"
	"github.com/twlabs/maeve-csms/manager/ocpp"
	types "github.com/twlabs/maeve-csms/manager/ocpp/ocpp201"
	"github.com/twlabs/maeve-csms/manager/services"
	"log"
)

type Get15118EvCertificateHandler struct {
	EvCertificateProvider services.EvCertificateProvider
}

func (g Get15118EvCertificateHandler) HandleCall(_ context.Context, _ string, request ocpp.Request) (ocpp.Response, error) {
	req := request.(*types.Get15118EVCertificateRequestJson)

	status := types.Iso15118EVCertificateStatusEnumTypeFailed
	response := types.Get15118EVCertificateResponseJson{
		Status: status,
	}
	if g.EvCertificateProvider != nil {
		res, err := g.EvCertificateProvider.ProvideCertificate(req.ExiRequest)

		if err != nil {
			log.Printf("failed to provide certificate: %v", err)
		} else {
			response = types.Get15118EVCertificateResponseJson{
				Status:      res.Status,
				ExiResponse: res.CertificateInstallationRes,
			}
		}
	}

	return &response, nil
}

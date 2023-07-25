// SPDX-License-Identifier: Apache-2.0

package ocpp201

import (
	"context"

	"github.com/thoughtworks/maeve-csms/manager/ocpp"
	types "github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"github.com/thoughtworks/maeve-csms/manager/services"
	"golang.org/x/exp/slog"
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
			slog.Error("failed to provide certificate", "err", err)
		} else {
			response = types.Get15118EVCertificateResponseJson{
				Status:      res.Status,
				ExiResponse: res.CertificateInstallationRes,
			}
		}
	}

	return &response, nil
}

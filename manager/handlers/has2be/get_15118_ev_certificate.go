// SPDX-License-Identifier: Apache-2.0

package has2be

import (
	"context"
	handlers201 "github.com/thoughtworks/maeve-csms/manager/handlers/ocpp201"
	"github.com/thoughtworks/maeve-csms/manager/ocpp"
	typesHasToBe "github.com/thoughtworks/maeve-csms/manager/ocpp/has2be"
	types201 "github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
)

type Get15118EvCertificateHandler struct {
	Handler201 handlers201.Get15118EvCertificateHandler
}

func (g Get15118EvCertificateHandler) HandleCall(ctx context.Context, _ string, request ocpp.Request) (ocpp.Response, error) {
	req := request.(*typesHasToBe.Get15118EVCertificateRequestJson)

	req201 := types201.Get15118EVCertificateRequestJson{
		ExiRequest:            req.ExiRequest,
		Iso15118SchemaVersion: *req.A15118SchemaVersion,
		// the only difference for install vs update is in the ExiRequest - so this field is only required for validation
		Action: types201.CertificateActionEnumTypeInstall,
	}

	res, err := g.Handler201.HandleCall(ctx, "", &req201)
	if err != nil {
		return nil, err
	}
	res201 := res.(*types201.Get15118EVCertificateResponseJson)

	return &typesHasToBe.Get15118EVCertificateResponseJson{
		Status:      typesHasToBe.Iso15118EVCertificateStatusEnumType(res201.Status),
		ExiResponse: res201.ExiResponse,
	}, nil
}

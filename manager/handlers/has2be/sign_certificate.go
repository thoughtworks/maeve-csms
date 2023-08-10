// SPDX-License-Identifier: Apache-2.0

package has2be

import (
	"context"
	handlers201 "github.com/thoughtworks/maeve-csms/manager/handlers/ocpp201"
	"github.com/thoughtworks/maeve-csms/manager/ocpp"
	typesHasToBe "github.com/thoughtworks/maeve-csms/manager/ocpp/has2be"
	types201 "github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
)

type SignCertificateHandler struct {
	Handler201 handlers201.SignCertificateHandler
}

func (s SignCertificateHandler) HandleCall(ctx context.Context, chargeStationId string, request ocpp.Request) (ocpp.Response, error) {
	req := request.(*typesHasToBe.SignCertificateRequestJson)

	req201 := &types201.SignCertificateRequestJson{
		Csr: req.Csr,
	}

	if req.TypeOfCertificate != nil {
		req201 = &types201.SignCertificateRequestJson{
			Csr:             req.Csr,
			CertificateType: (*types201.CertificateSigningUseEnumType)(req.TypeOfCertificate),
		}
	}

	res, err := s.Handler201.HandleCall(ctx, chargeStationId, req201)
	if err != nil {
		return nil, err
	}
	res201 := res.(*types201.SignCertificateResponseJson)

	return &typesHasToBe.SignCertificateResponseJson{
		Status: typesHasToBe.GenericStatusEnumType(res201.Status),
	}, nil
}

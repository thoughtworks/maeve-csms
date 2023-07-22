// SPDX-License-Identifier: Apache-2.0

package ocpp201

import (
	"context"

	"github.com/thoughtworks/maeve-csms/manager/handlers"
	"github.com/thoughtworks/maeve-csms/manager/ocpp"
	types "github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"github.com/thoughtworks/maeve-csms/manager/services"
	"golang.org/x/exp/slog"
)

type SignCertificateHandler struct {
	CertificateSignerService services.CertificateSignerService
	CallMaker                handlers.CallMaker
}

func (s SignCertificateHandler) HandleCall(ctx context.Context, chargeStationId string, request ocpp.Request) (ocpp.Response, error) {
	req := request.(*types.SignCertificateRequestJson)

	certificateType := types.CertificateSigningUseEnumTypeV2GCertificate
	if req.CertificateType != nil {
		certificateType = *req.CertificateType
	}

	slog.Info("sign certificate", slog.String("certificateType", string(certificateType)))

	status := types.GenericStatusEnumTypeRejected

	if s.CertificateSignerService != nil {
		status = types.GenericStatusEnumTypeAccepted

		go func() {
			var certType services.CertificateType
			if certificateType == types.CertificateSigningUseEnumTypeChargingStationCertificate {
				certType = services.CertificateTypeCSO
			} else {
				certType = services.CertificateTypeV2G
			}

			pemChain, err := s.CertificateSignerService.SignCertificate(certType, req.Csr)
			if err != nil {
				slog.Error("failed to sign certificate", err)
			} else {
				certSignedReq := &types.CertificateSignedRequestJson{
					CertificateChain: pemChain,
					CertificateType:  &certificateType,
				}

				err = s.CallMaker.Send(ctx, chargeStationId, certSignedReq)
				if err != nil {
					slog.Error("failed to send certificate signed request", err)
				}
			}
		}()
	}

	return &types.SignCertificateResponseJson{
		Status: status,
	}, nil
}

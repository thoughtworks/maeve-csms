package ocpp201

import (
	"context"
	"github.com/twlabs/maeve-csms/manager/handlers"
	"github.com/twlabs/maeve-csms/manager/ocpp"
	types "github.com/twlabs/maeve-csms/manager/ocpp/ocpp201"
	"github.com/twlabs/maeve-csms/manager/services"
	"log"
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

	log.Printf("Sign certificate: %s", certificateType)

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
				log.Printf("failed to sign certificate: %v", err)
			} else {
				certSignedReq := &types.CertificateSignedRequestJson{
					CertificateChain: pemChain,
					CertificateType:  &certificateType,
				}

				err = s.CallMaker.Send(ctx, chargeStationId, certSignedReq)
				if err != nil {
					log.Printf("failed to send certificate signed request: %v", err)
				}
			}
		}()
	}

	return &types.SignCertificateResponseJson{
		Status: status,
	}, nil
}

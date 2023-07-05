package ocpp201

import (
	"context"
	"github.com/twlabs/maeve-csms/manager/ocpp"
	types "github.com/twlabs/maeve-csms/manager/ocpp/ocpp201"
	"github.com/twlabs/maeve-csms/manager/services"
	"log"
)

type GetCertificateStatusHandler struct {
	CertificateValidationService services.CertificateValidationService
}

func (g GetCertificateStatusHandler) HandleCall(ctx context.Context, chargeStationId string, request ocpp.Request) (ocpp.Response, error) {
	req := request.(*types.GetCertificateStatusRequestJson)

	log.Printf("Get certificate status: %s", req.OcspRequestData.SerialNumber)

	status := types.GetCertificateStatusEnumTypeAccepted
	ocspResp, err := g.CertificateValidationService.ValidateHashedCertificateChain([]types.OCSPRequestDataType{req.OcspRequestData})
	if err != nil {
		log.Printf("validating hashed certificate chain: %v", err)
	}
	if ocspResp == nil {
		status = types.GetCertificateStatusEnumTypeFailed
	}

	return &types.GetCertificateStatusResponseJson{
		Status:     status,
		OcspResult: ocspResp,
	}, nil
}

// SPDX-License-Identifier: Apache-2.0

package has2be

import (
	"context"
	"github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/thoughtworks/maeve-csms/manager/ocpp"
	types "github.com/thoughtworks/maeve-csms/manager/ocpp/has2be"
	"github.com/thoughtworks/maeve-csms/manager/services"
)

type GetCertificateStatusHandler struct {
	CertificateValidationService services.CertificateValidationService
}

func (g GetCertificateStatusHandler) HandleCall(ctx context.Context, chargeStationId string, request ocpp.Request) (ocpp.Response, error) {
	span := trace.SpanFromContext(ctx)

	req := request.(*types.GetCertificateStatusRequestJson)

	span.SetAttributes(attribute.String("cert_status.serial_number", req.OcspRequestData.SerialNumber))

	ocpp201RequestData := ocpp201.OCSPRequestDataType{
		HashAlgorithm:  ocpp201.HashAlgorithmEnumType(req.OcspRequestData.HashAlgorithm),
		IssuerKeyHash:  req.OcspRequestData.IssuerKeyHash,
		IssuerNameHash: req.OcspRequestData.IssuerNameHash,
		SerialNumber:   req.OcspRequestData.SerialNumber,
	}

	if req.OcspRequestData.ResponderURL != nil {
		ocpp201RequestData.ResponderURL = *req.OcspRequestData.ResponderURL
	}

	status := types.GetCertificateStatusEnumTypeAccepted
	ocspResp, err := g.CertificateValidationService.ValidateHashedCertificateChain(ctx, []ocpp201.OCSPRequestDataType{ocpp201RequestData})
	if err != nil {
		span.SetAttributes(attribute.String("cert_status.error", err.Error()))
	}
	if ocspResp == nil {
		status = types.GetCertificateStatusEnumTypeFailed
	}

	span.SetAttributes(attribute.String("request.status", string(status)))

	return &types.GetCertificateStatusResponseJson{
		Status:     status,
		OcspResult: ocspResp,
	}, nil
}

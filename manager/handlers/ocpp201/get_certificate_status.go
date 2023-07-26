// SPDX-License-Identifier: Apache-2.0

package ocpp201

import (
	"context"

	"github.com/thoughtworks/maeve-csms/manager/ocpp"
	types "github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"github.com/thoughtworks/maeve-csms/manager/services"
	"golang.org/x/exp/slog"
)

type GetCertificateStatusHandler struct {
	CertificateValidationService services.CertificateValidationService
}

func (g GetCertificateStatusHandler) HandleCall(ctx context.Context, chargeStationId string, request ocpp.Request) (ocpp.Response, error) {
	req := request.(*types.GetCertificateStatusRequestJson)

	slog.Info("Get certificate status", slog.String("serialNumber", req.OcspRequestData.SerialNumber))

	status := types.GetCertificateStatusEnumTypeAccepted
	ocspResp, err := g.CertificateValidationService.ValidateHashedCertificateChain([]types.OCSPRequestDataType{req.OcspRequestData})
	if err != nil {
		slog.Error("validating hashed certificate chain", "err", err)
	}
	if ocspResp == nil {
		status = types.GetCertificateStatusEnumTypeFailed
	}

	return &types.GetCertificateStatusResponseJson{
		Status:     status,
		OcspResult: ocspResp,
	}, nil
}

// SPDX-License-Identifier: Apache-2.0

package ocpp201

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/thoughtworks/maeve-csms/manager/ocpp"
	types "github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"github.com/thoughtworks/maeve-csms/manager/services"
	"golang.org/x/exp/slog"
)

type GetCertificateStatusHandler struct {
	CertificateValidationService services.CertificateValidationService
}

func (g GetCertificateStatusHandler) HandleCall(ctx context.Context, chargeStationId string, request ocpp.Request) (ocpp.Response, error) {
	span := trace.SpanFromContext(ctx)

	req := request.(*types.GetCertificateStatusRequestJson)

	slog.Info("Get certificate status", slog.String("serialNumber", req.OcspRequestData.SerialNumber))

	status := types.GetCertificateStatusEnumTypeAccepted
	ocspResp, err := g.CertificateValidationService.ValidateHashedCertificateChain(ctx, []types.OCSPRequestDataType{req.OcspRequestData})
	if err != nil {
		slog.Error("validating hashed certificate chain", "err", err)
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

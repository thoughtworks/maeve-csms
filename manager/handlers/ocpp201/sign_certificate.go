// SPDX-License-Identifier: Apache-2.0

package ocpp201

import (
	"context"
	"github.com/thoughtworks/maeve-csms/manager/store"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/thoughtworks/maeve-csms/manager/ocpp"
	types "github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"github.com/thoughtworks/maeve-csms/manager/services"
	"golang.org/x/exp/slog"
)

type SignCertificateHandler struct {
	ChargeStationCertificateProvider services.ChargeStationCertificateProvider
	Store                            store.Engine
}

func (s SignCertificateHandler) HandleCall(ctx context.Context, chargeStationId string, request ocpp.Request) (ocpp.Response, error) {
	span := trace.SpanFromContext(ctx)

	req := request.(*types.SignCertificateRequestJson)

	certificateType := types.CertificateSigningUseEnumTypeV2GCertificate
	if req.CertificateType != nil {
		certificateType = *req.CertificateType
	}

	span.SetAttributes(attribute.String("sign_cert.cert_type", string(certificateType)))

	status := types.GenericStatusEnumTypeRejected

	if s.ChargeStationCertificateProvider != nil {
		var certType services.CertificateType
		var storeType store.CertificateType
		if certificateType == types.CertificateSigningUseEnumTypeChargingStationCertificate {
			certType = services.CertificateTypeCSO
			storeType = store.CertificateTypeChargeStation
		} else {
			certType = services.CertificateTypeV2G
			storeType = store.CertificateTypeEVCC
		}

		pemChain, err := s.ChargeStationCertificateProvider.ProvideCertificate(ctx, certType, req.Csr)
		if err != nil {
			slog.Error("failed to sign certificate", "err", err)
			span.AddEvent("failed to sign certificate", trace.WithAttributes(attribute.String("err", err.Error())))
		} else {
			certId, err := GetCertificateId(pemChain)
			if err != nil {
				slog.Error("failed to get certificate id", "err", err)
				span.AddEvent("failed to get certificate id", trace.WithAttributes(attribute.String("err", err.Error())))
			} else {
				err = s.Store.UpdateChargeStationInstallCertificates(ctx, chargeStationId, &store.ChargeStationInstallCertificates{
					Certificates: []*store.ChargeStationInstallCertificate{
						{
							CertificateType:               storeType,
							CertificateId:                 certId,
							CertificateData:               pemChain,
							CertificateInstallationStatus: store.CertificateInstallationPending,
						},
					},
				})
				if err != nil {
					slog.Error("failed to update charge station install certificates", "err", err)
					span.AddEvent("failed to update charge station install certificates", trace.WithAttributes(attribute.String("err", err.Error())))
				} else {
					status = types.GenericStatusEnumTypeAccepted
				}
			}
		}
	}

	span.SetAttributes(attribute.String("request.status", string(status)))

	return &types.SignCertificateResponseJson{
		Status: status,
	}, nil
}

// SPDX-License-Identifier: Apache-2.0

package ocpp201

import (
	"context"
	"github.com/thoughtworks/maeve-csms/manager/store"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/thoughtworks/maeve-csms/manager/ocpp"
	"github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
)

type CertificateSignedResultHandler struct {
	Store store.Engine
}

func (c CertificateSignedResultHandler) HandleCallResult(ctx context.Context, chargeStationId string, request ocpp.Request, response ocpp.Response, state any) error {
	req := request.(*ocpp201.CertificateSignedRequestJson)
	resp := response.(*ocpp201.CertificateSignedResponseJson)

	span := trace.SpanFromContext(ctx)

	if req.CertificateType == nil {
		typ := ocpp201.CertificateSigningUseEnumTypeV2GCertificate
		req.CertificateType = &typ
	}

	span.SetAttributes(
		attribute.String("certificate_signed.type", string(*req.CertificateType)),
		attribute.String("certificate_signed.status", string(resp.Status)),
	)

	certId, err := GetCertificateId(req.CertificateChain)
	if err != nil {
		return err
	}
	span.SetAttributes(attribute.String("certificate_signed.id", certId))

	var storeType store.CertificateType
	switch *req.CertificateType {
	case ocpp201.CertificateSigningUseEnumTypeV2GCertificate:
		storeType = store.CertificateTypeEVCC
	case ocpp201.CertificateSigningUseEnumTypeChargingStationCertificate:
		storeType = store.CertificateTypeChargeStation
	}

	var installStatus store.CertificateInstallationStatus
	switch resp.Status {
	case ocpp201.CertificateSignedStatusEnumTypeAccepted:
		installStatus = store.CertificateInstallationAccepted
	case ocpp201.CertificateSignedStatusEnumTypeRejected:
		installStatus = store.CertificateInstallationRejected
	}

	err = c.Store.UpdateChargeStationInstallCertificates(ctx, chargeStationId, &store.ChargeStationInstallCertificates{
		Certificates: []*store.ChargeStationInstallCertificate{
			{
				CertificateType:               storeType,
				CertificateId:                 certId,
				CertificateData:               req.CertificateChain,
				CertificateInstallationStatus: installStatus,
			},
		},
	})
	if err != nil {
		return err
	}

	return nil
}

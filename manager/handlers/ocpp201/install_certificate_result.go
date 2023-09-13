package ocpp201

import (
	"context"
	"github.com/thoughtworks/maeve-csms/manager/ocpp"
	"github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"github.com/thoughtworks/maeve-csms/manager/store"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type InstallCertificateResultHandler struct {
	Store store.Engine
}

func (i InstallCertificateResultHandler) HandleCallResult(ctx context.Context, chargeStationId string, request ocpp.Request, response ocpp.Response, state any) error {
	req := request.(*ocpp201.InstallCertificateRequestJson)
	resp := response.(*ocpp201.InstallCertificateResponseJson)

	span := trace.SpanFromContext(ctx)

	span.SetAttributes(
		attribute.String("install_certificate.type", string(req.CertificateType)),
		attribute.String("install_certificate.status", string(resp.Status)))

	certId, err := GetCertificateId(req.Certificate)
	if err != nil {
		return err
	}
	span.SetAttributes(attribute.String("install_certificate.id", certId))

	var storeType store.CertificateType
	switch req.CertificateType {
	case ocpp201.InstallCertificateUseEnumTypeV2GRootCertificate:
		storeType = store.CertificateTypeV2G
	case ocpp201.InstallCertificateUseEnumTypeMORootCertificate:
		storeType = store.CertificateTypeMO
	case ocpp201.InstallCertificateUseEnumTypeCSMSRootCertificate:
		storeType = store.CertificateTypeCSMS
	case ocpp201.InstallCertificateUseEnumTypeManufacturerRootCertificate:
		storeType = store.CertificateTypeMF
	}

	var installStatus store.CertificateInstallationStatus
	switch resp.Status {
	case ocpp201.InstallCertificateStatusEnumTypeAccepted:
		installStatus = store.CertificateInstallationAccepted
	case ocpp201.InstallCertificateStatusEnumTypeRejected:
		installStatus = store.CertificateInstallationRejected
	case ocpp201.InstallCertificateStatusEnumTypeFailed:
		installStatus = store.CertificateInstallationPending
	}

	err = i.Store.UpdateChargeStationInstallCertificates(ctx, chargeStationId, &store.ChargeStationInstallCertificates{
		Certificates: []*store.ChargeStationInstallCertificate{
			{
				CertificateType:               storeType,
				CertificateId:                 certId,
				CertificateData:               req.Certificate,
				CertificateInstallationStatus: installStatus,
			},
		},
	})
	if err != nil {
		return err
	}

	return nil
}

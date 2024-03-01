// SPDX-License-Identifier: Apache-2.0

package ocpp201_test

import (
	"context"
	"encoding/pem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thoughtworks/maeve-csms/manager/handlers"
	handlers201 "github.com/thoughtworks/maeve-csms/manager/handlers/ocpp201"
	"github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"github.com/thoughtworks/maeve-csms/manager/store"
	"github.com/thoughtworks/maeve-csms/manager/store/inmemory"
	"k8s.io/utils/clock"
	"testing"
)

func TestInstallCertificateResultHandler(t *testing.T) {
	engine := inmemory.NewStore(clock.RealClock{})
	handler := handlers201.InstallCertificateResultHandler{Store: engine}

	pemBlock := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: []byte("test"),
	}
	pemBytes := pem.EncodeToMemory(pemBlock)
	id, err := handlers201.GetCertificateId(string(pemBytes))
	require.NoError(t, err)

	respAccepted := &ocpp201.InstallCertificateResponseJson{
		Status: ocpp201.InstallCertificateStatusEnumTypeAccepted,
	}
	respRejected := &ocpp201.InstallCertificateResponseJson{
		Status: ocpp201.InstallCertificateStatusEnumTypeRejected,
	}
	respFailed := &ocpp201.InstallCertificateResponseJson{
		Status: ocpp201.InstallCertificateStatusEnumTypeFailed,
	}

	testCases := []struct {
		name        string
		req         *ocpp201.InstallCertificateRequestJson
		resp        *ocpp201.InstallCertificateResponseJson
		storeType   store.CertificateType
		storeStatus store.CertificateInstallationStatus
		attrs       map[string]any
	}{
		{
			name: "V2G accepted",
			req: &ocpp201.InstallCertificateRequestJson{
				Certificate:     string(pemBytes),
				CertificateType: ocpp201.InstallCertificateUseEnumTypeV2GRootCertificate,
			},
			resp:        respAccepted,
			storeType:   store.CertificateTypeV2G,
			storeStatus: store.CertificateInstallationAccepted,
			attrs: map[string]any{
				"install_certificate.type":   "V2GRootCertificate",
				"install_certificate.status": "Accepted",
				"install_certificate.id":     id,
			},
		},
		{
			name: "V2G rejected",
			req: &ocpp201.InstallCertificateRequestJson{
				Certificate:     string(pemBytes),
				CertificateType: ocpp201.InstallCertificateUseEnumTypeV2GRootCertificate,
			},
			resp:        respRejected,
			storeType:   store.CertificateTypeV2G,
			storeStatus: store.CertificateInstallationRejected,
			attrs: map[string]any{
				"install_certificate.type":   "V2GRootCertificate",
				"install_certificate.status": "Rejected",
				"install_certificate.id":     id,
			},
		},
		{
			name: "V2G failed",
			req: &ocpp201.InstallCertificateRequestJson{
				Certificate:     string(pemBytes),
				CertificateType: ocpp201.InstallCertificateUseEnumTypeV2GRootCertificate,
			},
			resp:        respFailed,
			storeType:   store.CertificateTypeV2G,
			storeStatus: store.CertificateInstallationPending,
			attrs: map[string]any{
				"install_certificate.type":   "V2GRootCertificate",
				"install_certificate.status": "Failed",
				"install_certificate.id":     id,
			},
		},
		{
			name: "MO accepted",
			req: &ocpp201.InstallCertificateRequestJson{
				Certificate:     string(pemBytes),
				CertificateType: ocpp201.InstallCertificateUseEnumTypeMORootCertificate,
			},
			resp:        respAccepted,
			storeType:   store.CertificateTypeMO,
			storeStatus: store.CertificateInstallationAccepted,
			attrs: map[string]any{
				"install_certificate.type":   "MORootCertificate",
				"install_certificate.status": "Accepted",
				"install_certificate.id":     id,
			},
		},
		{
			name: "MF accepted",
			req: &ocpp201.InstallCertificateRequestJson{
				Certificate:     string(pemBytes),
				CertificateType: ocpp201.InstallCertificateUseEnumTypeManufacturerRootCertificate,
			},
			resp:        respAccepted,
			storeType:   store.CertificateTypeMF,
			storeStatus: store.CertificateInstallationAccepted,
			attrs: map[string]any{
				"install_certificate.type":   "ManufacturerRootCertificate",
				"install_certificate.status": "Accepted",
				"install_certificate.id":     id,
			},
		},
		{
			name: "CSMS accepted",
			req: &ocpp201.InstallCertificateRequestJson{
				Certificate:     string(pemBytes),
				CertificateType: ocpp201.InstallCertificateUseEnumTypeCSMSRootCertificate,
			},
			resp:        respAccepted,
			storeType:   store.CertificateTypeCSMS,
			storeStatus: store.CertificateInstallationAccepted,
			attrs: map[string]any{
				"install_certificate.type":   "CSMSRootCertificate",
				"install_certificate.status": "Accepted",
				"install_certificate.id":     id,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tracer, exporter := handlers.GetTracer()

			ctx := context.Background()

			func() {
				ctx, span := tracer.Start(ctx, tc.name)
				defer span.End()

				err := handler.HandleCallResult(ctx, "test", tc.req, tc.resp, nil)
				require.NoError(t, err)
			}()

			handlers.AssertSpan(t, &exporter.GetSpans()[0], tc.name, tc.attrs)

			certs, err := engine.LookupChargeStationInstallCertificates(ctx, "test")
			require.NoError(t, err)

			require.Len(t, certs.Certificates, 1)
			assert.Equal(t, id, certs.Certificates[0].CertificateId)
			assert.Equal(t, string(pemBytes), certs.Certificates[0].CertificateData)
			assert.Equal(t, tc.storeType, certs.Certificates[0].CertificateType)
			assert.Equal(t, tc.storeStatus, certs.Certificates[0].CertificateInstallationStatus)
		})
	}
}

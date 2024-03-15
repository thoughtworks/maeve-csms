// SPDX-License-Identifier: Apache-2.0

package ocpp201_test

import (
	"context"
	"encoding/pem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	handlers201 "github.com/thoughtworks/maeve-csms/manager/handlers/ocpp201"
	"github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"github.com/thoughtworks/maeve-csms/manager/store"
	"github.com/thoughtworks/maeve-csms/manager/store/inmemory"
	"github.com/thoughtworks/maeve-csms/manager/testutil"
	"k8s.io/utils/clock"
	"testing"
)

func TestCertificateSignedResultHandler(t *testing.T) {
	engine := inmemory.NewStore(clock.RealClock{})
	handler := handlers201.CertificateSignedResultHandler{Store: engine}

	pemBlock := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: []byte("test"),
	}
	pemBytes := pem.EncodeToMemory(pemBlock)
	id, err := handlers201.GetCertificateId(string(pemBytes))
	require.NoError(t, err)

	typV2G := ocpp201.CertificateSigningUseEnumTypeV2GCertificate
	typChargingStation := ocpp201.CertificateSigningUseEnumTypeChargingStationCertificate

	respAccepted := &ocpp201.CertificateSignedResponseJson{
		Status: ocpp201.CertificateSignedStatusEnumTypeAccepted,
	}
	respRejected := &ocpp201.CertificateSignedResponseJson{
		Status: ocpp201.CertificateSignedStatusEnumTypeRejected,
	}

	testCases := []struct {
		name        string
		req         *ocpp201.CertificateSignedRequestJson
		resp        *ocpp201.CertificateSignedResponseJson
		storeType   store.CertificateType
		storeStatus store.CertificateInstallationStatus
		attrs       map[string]any
	}{
		{
			name: "ChargingStation accepted",
			req: &ocpp201.CertificateSignedRequestJson{
				CertificateChain: string(pemBytes),
				CertificateType:  &typChargingStation,
			},
			resp:        respAccepted,
			storeType:   store.CertificateTypeChargeStation,
			storeStatus: store.CertificateInstallationAccepted,
			attrs: map[string]any{
				"certificate_signed.type":   "ChargingStationCertificate",
				"certificate_signed.status": "Accepted",
				"certificate_signed.id":     id,
			},
		},
		{
			name: "ChargingStation rejected",
			req: &ocpp201.CertificateSignedRequestJson{
				CertificateChain: string(pemBytes),
				CertificateType:  &typChargingStation,
			},
			resp:        respRejected,
			storeType:   store.CertificateTypeChargeStation,
			storeStatus: store.CertificateInstallationRejected,
			attrs: map[string]any{
				"certificate_signed.type":   "ChargingStationCertificate",
				"certificate_signed.status": "Rejected",
				"certificate_signed.id":     id,
			},
		},
		{
			name: "V2G accepted",
			req: &ocpp201.CertificateSignedRequestJson{
				CertificateChain: string(pemBytes),
				CertificateType:  &typV2G,
			},
			resp:        respAccepted,
			storeType:   store.CertificateTypeEVCC,
			storeStatus: store.CertificateInstallationAccepted,
			attrs: map[string]any{
				"certificate_signed.type":   "V2GCertificate",
				"certificate_signed.status": "Accepted",
				"certificate_signed.id":     id,
			},
		},
		{
			name: "V2G rejected",
			req: &ocpp201.CertificateSignedRequestJson{
				CertificateChain: string(pemBytes),
				CertificateType:  &typV2G,
			},
			resp:        respRejected,
			storeType:   store.CertificateTypeEVCC,
			storeStatus: store.CertificateInstallationRejected,
			attrs: map[string]any{
				"certificate_signed.type":   "V2GCertificate",
				"certificate_signed.status": "Rejected",
				"certificate_signed.id":     id,
			},
		},
		{
			name: "nil type",
			req: &ocpp201.CertificateSignedRequestJson{
				CertificateChain: string(pemBytes),
			},
			resp:        respAccepted,
			storeType:   store.CertificateTypeEVCC,
			storeStatus: store.CertificateInstallationAccepted,
			attrs: map[string]any{
				"certificate_signed.type":   "V2GCertificate",
				"certificate_signed.status": "Accepted",
				"certificate_signed.id":     id,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tracer, exporter := testutil.GetTracer()

			ctx := context.Background()

			func() {
				ctx, span := tracer.Start(ctx, tc.name)
				defer span.End()

				err := handler.HandleCallResult(ctx, "test", tc.req, tc.resp, nil)
				require.NoError(t, err)
			}()

			testutil.AssertSpan(t, &exporter.GetSpans()[0], tc.name, tc.attrs)

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

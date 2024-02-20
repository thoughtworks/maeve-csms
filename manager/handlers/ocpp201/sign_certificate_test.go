// SPDX-License-Identifier: Apache-2.0

package ocpp201_test

import (
	"context"
	"encoding/pem"
	"github.com/stretchr/testify/require"
	"github.com/thoughtworks/maeve-csms/manager/handlers"
	handlers201 "github.com/thoughtworks/maeve-csms/manager/handlers/ocpp201"
	"github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"github.com/thoughtworks/maeve-csms/manager/services"
	"github.com/thoughtworks/maeve-csms/manager/store"
	"github.com/thoughtworks/maeve-csms/manager/store/inmemory"
	"k8s.io/utils/clock"
	"testing"
)

type mockCertificateProvider struct{}

func (m mockCertificateProvider) ProvideCertificate(context.Context, services.CertificateType, string, string) (pemEncodedCertificateChain string, err error) {
	block := pem.Block{
		Type:  "CERTIFICATE",
		Bytes: []byte("test"),
	}
	return string(pem.EncodeToMemory(&block)), nil
}

func TestSignCertificateHandler(t *testing.T) {
	engine := inmemory.NewStore(clock.RealClock{})

	typ := ocpp201.CertificateSigningUseEnumTypeV2GCertificate
	req := &ocpp201.SignCertificateRequestJson{
		CertificateType: &typ,
		Csr:             "test",
	}

	handler := handlers201.SignCertificateHandler{
		ChargeStationCertificateProvider: mockCertificateProvider{},
		Store:                            engine,
	}

	tracer, exporter := handlers.GetTracer()

	func() {
		ctx, span := tracer.Start(context.Background(), t.Name())
		defer span.End()

		response, err := handler.HandleCall(ctx, "test", req)
		require.NoError(t, err)

		resp, ok := response.(*ocpp201.SignCertificateResponseJson)
		require.True(t, ok)
		require.Equal(t, ocpp201.GenericStatusEnumTypeAccepted, resp.Status)
	}()

	handlers.AssertSpan(t, &exporter.GetSpans()[0], t.Name(), map[string]any{
		"sign_cert.cert_type": "V2GCertificate",
		"request.status":      "Accepted",
	})

	certs, err := engine.LookupChargeStationInstallCertificates(context.Background(), "test")
	require.NoError(t, err)
	require.Len(t, certs.Certificates, 1)
	require.Equal(t, store.CertificateTypeEVCC, certs.Certificates[0].CertificateType)
	require.Equal(t, store.CertificateInstallationPending, certs.Certificates[0].CertificateInstallationStatus)
}

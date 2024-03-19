// SPDX-License-Identifier: Apache-2.0

package ocpp201_test

import (
	"context"
	"github.com/stretchr/testify/require"
	"github.com/thoughtworks/maeve-csms/manager/handlers/ocpp201"
	types "github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"github.com/thoughtworks/maeve-csms/manager/testutil"
	"testing"
)

func TestGetInstalledCertificateIdsResultHandler(t *testing.T) {
	handler := ocpp201.GetInstalledCertificateIdsResultHandler{}

	tracer, exporter := testutil.GetTracer()

	ctx := context.Background()

	func() {
		ctx, span := tracer.Start(ctx, `test`)
		defer span.End()

		req := &types.GetInstalledCertificateIdsRequestJson{
			CertificateType: []types.GetCertificateIdUseEnumType{
				types.GetCertificateIdUseEnumTypeCSMSRootCertificate,
				types.GetCertificateIdUseEnumTypeMORootCertificate,
			},
		}
		resp := &types.GetInstalledCertificateIdsResponseJson{
			Status: types.GetInstalledCertificateStatusEnumTypeAccepted,
			CertificateHashDataChain: []types.CertificateHashDataChainType{
				{
					CertificateHashData: types.CertificateHashDataType{
						HashAlgorithm:  types.HashAlgorithmEnumTypeSHA256,
						IssuerKeyHash:  "ABC123",
						IssuerNameHash: "ABCDEF",
						SerialNumber:   "12345678",
					},
					CertificateType:          types.GetCertificateIdUseEnumTypeCSMSRootCertificate,
					ChildCertificateHashData: nil,
				},
			},
		}

		err := handler.HandleCallResult(ctx, "cs001", req, resp, nil)
		require.NoError(t, err)
	}()

	testutil.AssertSpan(t, &exporter.GetSpans()[0], "test", map[string]any{
		"get_installed_certificate.types":  "CSMSRootCertificate,MORootCertificate",
		"get_installed_certificate.status": "Accepted",
	})
}

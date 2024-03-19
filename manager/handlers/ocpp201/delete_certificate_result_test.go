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

func TestDeleteCertificateResultHandler(t *testing.T) {
	handler := ocpp201.DeleteCertificateResultHandler{}

	tracer, exporter := testutil.GetTracer()

	ctx := context.Background()

	func() {
		ctx, span := tracer.Start(ctx, `test`)
		defer span.End()

		req := &types.DeleteCertificateRequestJson{
			CertificateHashData: types.CertificateHashDataType{
				HashAlgorithm:  types.HashAlgorithmEnumTypeSHA256,
				IssuerKeyHash:  "ABC123",
				IssuerNameHash: "ABCDEF",
				SerialNumber:   "12345678",
			},
		}
		resp := &types.DeleteCertificateResponseJson{
			Status: types.DeleteCertificateStatusEnumTypeAccepted,
		}

		err := handler.HandleCallResult(ctx, "cs001", req, resp, nil)
		require.NoError(t, err)
	}()

	testutil.AssertSpan(t, &exporter.GetSpans()[0], "test", map[string]any{
		"delete_certificate.serial_number": "12345678",
		"delete_certificate.status":        "Accepted",
	})
}

package has2be_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	handlers "github.com/thoughtworks/maeve-csms/manager/handlers/has2be"
	types "github.com/thoughtworks/maeve-csms/manager/ocpp/has2be"
	"github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"github.com/thoughtworks/maeve-csms/manager/services"
	"testing"
)

type dummyCertificateValidationService struct {
	T *testing.T
}

func (d dummyCertificateValidationService) ValidatePEMCertificateChain(ctx context.Context, pemChain []byte, eMAID string) (*string, error) {
	d.T.Fatal("not implemented")
	return nil, nil
}

func (d dummyCertificateValidationService) ValidateHashedCertificateChain(ctx context.Context, ocspRequestData []ocpp201.OCSPRequestDataType) (*string, error) {
	switch ocspRequestData[0].SerialNumber {
	case "invalid-chain":
		return nil, services.ValidationErrorCertChain
	case "revoked":
		ocspResult := "revoked"
		return &ocspResult, services.ValidationErrorCertRevoked
	default:
		ocspResult := "ocsp-result"
		return &ocspResult, nil
	}
}

func TestGetCertificateStatus(t *testing.T) {
	req := &types.GetCertificateStatusRequestJson{
		OcspRequestData: types.OCSPRequestDataType{
			HashAlgorithm:  "SHA256",
			IssuerKeyHash:  "key-hash",
			IssuerNameHash: "name-hash",
			SerialNumber:   "serial-number",
		},
	}

	h := handlers.GetCertificateStatusHandler{
		CertificateValidationService: dummyCertificateValidationService{T: t},
	}

	got, err := h.HandleCall(context.Background(), "cs001", req)
	require.NoError(t, err)

	ocspResult := "ocsp-result"
	want := &types.GetCertificateStatusResponseJson{
		Status:     types.GetCertificateStatusEnumTypeAccepted,
		OcspResult: &ocspResult,
	}

	assert.Equal(t, want, got)
}

func TestGetCertificateStatusInvalidChain(t *testing.T) {
	req := &types.GetCertificateStatusRequestJson{
		OcspRequestData: types.OCSPRequestDataType{
			HashAlgorithm:  "SHA256",
			IssuerKeyHash:  "key-hash",
			IssuerNameHash: "name-hash",
			SerialNumber:   "invalid-chain",
		},
	}

	h := handlers.GetCertificateStatusHandler{
		CertificateValidationService: dummyCertificateValidationService{T: t},
	}

	got, err := h.HandleCall(context.Background(), "cs001", req)
	require.NoError(t, err)

	want := &types.GetCertificateStatusResponseJson{
		Status: types.GetCertificateStatusEnumTypeFailed,
	}

	assert.Equal(t, want, got)
}

func TestGetCertificateStatusRevoked(t *testing.T) {
	req := &types.GetCertificateStatusRequestJson{
		OcspRequestData: types.OCSPRequestDataType{
			HashAlgorithm:  "SHA256",
			IssuerKeyHash:  "key-hash",
			IssuerNameHash: "name-hash",
			SerialNumber:   "revoked",
		},
	}

	h := handlers.GetCertificateStatusHandler{
		CertificateValidationService: dummyCertificateValidationService{T: t},
	}

	got, err := h.HandleCall(context.Background(), "cs001", req)
	require.NoError(t, err)

	ocspResult := "revoked"
	want := &types.GetCertificateStatusResponseJson{
		Status:     types.GetCertificateStatusEnumTypeAccepted,
		OcspResult: &ocspResult,
	}

	assert.Equal(t, want, got)
}

package ocpp201_test

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	handlers "github.com/thoughtworks/maeve-csms/manager/handlers/ocpp201"
	types "github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"github.com/thoughtworks/maeve-csms/manager/services"
	"testing"
)

type mockCertValidationService struct {
}

func (m mockCertValidationService) ValidatePEMCertificateChain(certificate []byte, eMAID string) (*string, error) {
	switch string(certificate) {
	case "invalidCertChain":
		return nil, services.ValidationErrorCertChain
	case "revokedCertChain":
		return nil, services.ValidationErrorCertRevoked
	case "expiredCertChain":
		return nil, services.ValidationErrorCertExpired
	case "signatureError":
		return nil, errors.New("error")
	}
	return nil, nil
}

func (m mockCertValidationService) ValidateHashedCertificateChain(ocspRequestData []types.OCSPRequestDataType) (*string, error) {
	if len(ocspRequestData) > 0 {
		switch ocspRequestData[0].SerialNumber {
		case "invalidCertChain":
			return nil, services.ValidationErrorCertChain
		case "revokedCertChain":
			return nil, services.ValidationErrorCertRevoked
		case "expiredCertChain":
			return nil, services.ValidationErrorCertExpired
		case "signatureError":
			return nil, errors.New("error")
		}
	}

	return nil, nil
}

func TestAuthorizeKnownRfidCard(t *testing.T) {
	ah := handlers.AuthorizeHandler{
		TokenStore: services.InMemoryTokenStore{
			Tokens: map[string]*services.Token{
				"ISO14443:MYRFIDCARD": {
					Type: string(types.IdTokenEnumTypeISO14443),
					Uid:  "MYRFIDCARD",
				},
			},
		},
	}

	req := &types.AuthorizeRequestJson{
		IdToken: types.IdTokenType{
			Type:    types.IdTokenEnumTypeISO14443,
			IdToken: "MYRFIDCARD",
		},
	}

	got, err := ah.HandleCall(context.Background(), "cs001", req)
	assert.NoError(t, err)

	want := &types.AuthorizeResponseJson{
		IdTokenInfo: types.IdTokenInfoType{
			Status: types.AuthorizationStatusEnumTypeAccepted,
		},
	}

	assert.Equal(t, want, got)
}

func TestAuthorizeWithUnknownRfidCard(t *testing.T) {
	ah := handlers.AuthorizeHandler{
		TokenStore: services.InMemoryTokenStore{
			Tokens: map[string]*services.Token{
				"ISO14443:MYRFIDCARD": {
					Type: string(types.IdTokenEnumTypeISO14443),
					Uid:  "MYRFIDCARD",
				},
			},
		},
	}

	req := &types.AuthorizeRequestJson{
		IdToken: types.IdTokenType{
			Type:    types.IdTokenEnumTypeISO14443,
			IdToken: "MYBADRFID",
		},
	}

	got, err := ah.HandleCall(context.Background(), "cs001", req)
	assert.NoError(t, err)

	want := &types.AuthorizeResponseJson{
		IdTokenInfo: types.IdTokenInfoType{
			Status: types.AuthorizationStatusEnumTypeUnknown,
		},
	}

	assert.Equal(t, want, got)
}

func TestAuthorizeWithEmaidAndCertificateChain(t *testing.T) {
	ah := handlers.AuthorizeHandler{
		TokenStore: services.InMemoryTokenStore{
			Tokens: map[string]*services.Token{
				"eMAID:MYEMAID": {
					Type: string(types.IdTokenEnumTypeEMAID),
					Uid:  "MYEMAID",
				},
			},
		},
		CertificateValidationService: mockCertValidationService{},
	}

	mockCertificate := "mockcert"
	req := &types.AuthorizeRequestJson{
		IdToken: types.IdTokenType{
			Type:    types.IdTokenEnumTypeEMAID,
			IdToken: "MYEMAID",
		},
		Certificate: &mockCertificate,
	}

	got, err := ah.HandleCall(context.Background(), "cs001", req)
	assert.NoError(t, err)

	certStatus := types.AuthorizeCertificateStatusEnumTypeAccepted
	want := &types.AuthorizeResponseJson{
		IdTokenInfo: types.IdTokenInfoType{
			Status: types.AuthorizationStatusEnumTypeAccepted,
		},
		CertificateStatus: &certStatus,
	}

	assert.Equal(t, want, got)
}

func TestAuthorizeWithEmaidAndInvalidCertificateChain(t *testing.T) {
	ah := handlers.AuthorizeHandler{
		TokenStore: services.InMemoryTokenStore{
			Tokens: map[string]*services.Token{
				"eMAID:MYEMAID": {
					Type: string(types.IdTokenEnumTypeEMAID),
					Uid:  "MYEMAID",
				},
			},
		},
		CertificateValidationService: mockCertValidationService{},
	}

	testCases := []string{"invalidCertChain", "revokedCertChain", "expiredCertChain", "signatureError"}
	expectedErrors := []types.AuthorizeCertificateStatusEnumType{
		types.AuthorizeCertificateStatusEnumTypeCertChainError,
		types.AuthorizeCertificateStatusEnumTypeCertificateRevoked,
		types.AuthorizeCertificateStatusEnumTypeCertificateExpired,
		types.AuthorizeCertificateStatusEnumTypeSignatureError,
	}

	for index, testCase := range testCases {
		t.Run(testCase, func(t *testing.T) {
			req := &types.AuthorizeRequestJson{
				IdToken: types.IdTokenType{
					Type:    types.IdTokenEnumTypeEMAID,
					IdToken: "MYEMAID",
				},
				Certificate: &testCase,
			}

			got, err := ah.HandleCall(context.Background(), "cs001", req)
			assert.NoError(t, err)

			want := &types.AuthorizeResponseJson{
				IdTokenInfo: types.IdTokenInfoType{
					Status: types.AuthorizationStatusEnumTypeBlocked,
				},
				CertificateStatus: &expectedErrors[index],
			}

			assert.Equal(t, want, got)
		})
	}

}

func TestAuthorizeWithEmaidAndCertificateHashes(t *testing.T) {
	ah := handlers.AuthorizeHandler{
		TokenStore: services.InMemoryTokenStore{
			Tokens: map[string]*services.Token{
				"eMAID:MYEMAID": {
					Type: string(types.IdTokenEnumTypeEMAID),
					Uid:  "MYEMAID",
				},
			},
		},
		CertificateValidationService: mockCertValidationService{},
	}

	req := &types.AuthorizeRequestJson{
		IdToken: types.IdTokenType{
			Type:    types.IdTokenEnumTypeEMAID,
			IdToken: "MYEMAID",
		},
		Iso15118CertificateHashData: &[]types.OCSPRequestDataType{
			{
				SerialNumber: "mockCertificate",
			},
		},
	}

	got, err := ah.HandleCall(context.Background(), "cs001", req)
	assert.NoError(t, err)

	certStatus := types.AuthorizeCertificateStatusEnumTypeAccepted
	want := &types.AuthorizeResponseJson{
		IdTokenInfo: types.IdTokenInfoType{
			Status: types.AuthorizationStatusEnumTypeAccepted,
		},
		CertificateStatus: &certStatus,
	}

	assert.Equal(t, want, got)
}

func TestAuthorizeWithEmaidAndInvalidCertificateHashes(t *testing.T) {
	ah := handlers.AuthorizeHandler{
		TokenStore: services.InMemoryTokenStore{
			Tokens: map[string]*services.Token{
				"eMAID:MYEMAID": {
					Type: string(types.IdTokenEnumTypeEMAID),
					Uid:  "MYEMAID",
				},
			},
		},
		CertificateValidationService: mockCertValidationService{},
	}

	testCases := []string{"invalidCertChain", "revokedCertChain", "expiredCertChain", "signatureError"}
	expectedErrors := []types.AuthorizeCertificateStatusEnumType{
		types.AuthorizeCertificateStatusEnumTypeCertChainError,
		types.AuthorizeCertificateStatusEnumTypeCertificateRevoked,
		types.AuthorizeCertificateStatusEnumTypeCertificateExpired,
		types.AuthorizeCertificateStatusEnumTypeSignatureError,
	}

	for index, testCase := range testCases {
		t.Run(testCase, func(t *testing.T) {
			req := &types.AuthorizeRequestJson{
				IdToken: types.IdTokenType{
					Type:    types.IdTokenEnumTypeEMAID,
					IdToken: "MYEMAID",
				},
				Iso15118CertificateHashData: &[]types.OCSPRequestDataType{
					{
						SerialNumber: testCase,
					},
				},
			}

			got, err := ah.HandleCall(context.Background(), "cs001", req)
			assert.NoError(t, err)

			want := &types.AuthorizeResponseJson{
				IdTokenInfo: types.IdTokenInfoType{
					Status: types.AuthorizationStatusEnumTypeBlocked,
				},
				CertificateStatus: &expectedErrors[index],
			}

			assert.Equal(t, want, got)
		})
	}

}

func TestAuthorizeWithEmaidAndNoCertificateData(t *testing.T) {
	ah := handlers.AuthorizeHandler{
		TokenStore: services.InMemoryTokenStore{
			Tokens: map[string]*services.Token{
				"eMAID:MYEMAID": {
					Type: string(types.IdTokenEnumTypeEMAID),
					Uid:  "MYEMAID",
				},
			},
		},
		CertificateValidationService: mockCertValidationService{},
	}

	req := &types.AuthorizeRequestJson{
		IdToken: types.IdTokenType{
			Type:    types.IdTokenEnumTypeEMAID,
			IdToken: "MYEMAID",
		},
	}

	got, err := ah.HandleCall(context.Background(), "cs001", req)
	assert.NoError(t, err)

	want := &types.AuthorizeResponseJson{
		IdTokenInfo: types.IdTokenInfoType{
			Status: types.AuthorizationStatusEnumTypeAccepted,
		},
	}

	assert.Equal(t, want, got)
}

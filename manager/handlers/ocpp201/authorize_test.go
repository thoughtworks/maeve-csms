// SPDX-License-Identifier: Apache-2.0

package ocpp201_test

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	handlers "github.com/thoughtworks/maeve-csms/manager/handlers/ocpp201"
	types "github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"github.com/thoughtworks/maeve-csms/manager/services"
	"github.com/thoughtworks/maeve-csms/manager/store"
	"github.com/thoughtworks/maeve-csms/manager/store/inmemory"
	clockutil "k8s.io/utils/clock"
	"testing"
	"time"
)

type mockCertValidationService struct {
}

func (m mockCertValidationService) ValidatePEMCertificateChain(_ context.Context, certificate []byte, _ string) (*string, error) {
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

func (m mockCertValidationService) ValidateHashedCertificateChain(_ context.Context, ocspRequestData []types.OCSPRequestDataType) (*string, error) {
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

func setupTokenStore(tokenStore store.TokenStore) error {
	err := tokenStore.SetToken(context.Background(), &store.Token{
		CountryCode: "GB",
		PartyId:     "TWK",
		Type:        "RFID",
		Uid:         "MYRFIDCARD",
		ContractId:  "GBTWK012345678V",
		Issuer:      "Thoughtworks",
		Valid:       true,
		CacheMode:   "ALWAYS",
		LastUpdated: time.Now().Format(time.RFC3339),
	})
	if err != nil {
		return err
	}
	err = tokenStore.SetToken(context.Background(), &store.Token{
		CountryCode: "GB",
		PartyId:     "TWK",
		Type:        "OTHER",
		Uid:         "MYEMAID",
		ContractId:  "GBTWK123456789B",
		Issuer:      "Thoughtworks",
		Valid:       true,
		CacheMode:   "ALWAYS",
		LastUpdated: time.Now().Format(time.RFC3339),
	})
	if err != nil {
		return err
	}
	return nil
}

func TestAuthorizeKnownRfidCard(t *testing.T) {
	clock := clockutil.RealClock{}
	engine := inmemory.NewStore(clock)
	err := setupTokenStore(engine)
	require.NoError(t, err)
	tokenAuthService := &services.OcppTokenAuthService{
		Clock:      clock,
		TokenStore: engine,
	}

	ah := handlers.AuthorizeHandler{
		TokenAuthService:             tokenAuthService,
		CertificateValidationService: mockCertValidationService{},
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
	clock := clockutil.RealClock{}
	engine := inmemory.NewStore(clock)
	err := setupTokenStore(engine)
	require.NoError(t, err)
	tokenAuthService := &services.OcppTokenAuthService{
		Clock:      clock,
		TokenStore: engine,
	}

	ah := handlers.AuthorizeHandler{
		TokenAuthService:             tokenAuthService,
		CertificateValidationService: mockCertValidationService{},
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
	clock := clockutil.RealClock{}
	engine := inmemory.NewStore(clock)
	err := setupTokenStore(engine)
	require.NoError(t, err)
	tokenAuthService := &services.OcppTokenAuthService{
		Clock:      clock,
		TokenStore: engine,
	}

	ah := handlers.AuthorizeHandler{
		TokenAuthService:             tokenAuthService,
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
	clock := clockutil.RealClock{}
	engine := inmemory.NewStore(clock)
	err := setupTokenStore(engine)
	require.NoError(t, err)
	tokenAuthService := &services.OcppTokenAuthService{
		Clock:      clock,
		TokenStore: engine,
	}

	ah := handlers.AuthorizeHandler{
		TokenAuthService:             tokenAuthService,
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
	clock := clockutil.RealClock{}
	engine := inmemory.NewStore(clock)
	err := setupTokenStore(engine)
	require.NoError(t, err)
	tokenAuthService := &services.OcppTokenAuthService{
		Clock:      clock,
		TokenStore: engine,
	}

	ah := handlers.AuthorizeHandler{
		TokenAuthService:             tokenAuthService,
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
	clock := clockutil.RealClock{}
	engine := inmemory.NewStore(clock)
	err := setupTokenStore(engine)
	require.NoError(t, err)
	tokenAuthService := &services.OcppTokenAuthService{
		Clock:      clock,
		TokenStore: engine,
	}

	ah := handlers.AuthorizeHandler{
		TokenAuthService:             tokenAuthService,
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
	clock := clockutil.RealClock{}
	engine := inmemory.NewStore(clock)
	err := setupTokenStore(engine)
	require.NoError(t, err)
	tokenAuthService := &services.OcppTokenAuthService{
		Clock:      clock,
		TokenStore: engine,
	}

	ah := handlers.AuthorizeHandler{
		TokenAuthService:             tokenAuthService,
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

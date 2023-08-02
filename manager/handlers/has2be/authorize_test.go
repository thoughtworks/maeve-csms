package has2be_test

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	handlers "github.com/thoughtworks/maeve-csms/manager/handlers/has2be"
	types "github.com/thoughtworks/maeve-csms/manager/ocpp/has2be"
	"github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"github.com/thoughtworks/maeve-csms/manager/services"
	"github.com/thoughtworks/maeve-csms/manager/store"
	"github.com/thoughtworks/maeve-csms/manager/store/inmemory"
	"testing"
	"time"
)

type mockCertValidationService struct {
}

func (m mockCertValidationService) ValidatePEMCertificateChain(ctx context.Context, certificate []byte, eMAID string) (*string, error) {
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

func (m mockCertValidationService) ValidateHashedCertificateChain(ctx context.Context, ocspRequestData []ocpp201.OCSPRequestDataType) (*string, error) {
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
		CacheMode:   "NEVER",
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
		CacheMode:   "NEVER",
		LastUpdated: time.Now().Format(time.RFC3339),
	})
	if err != nil {
		return err
	}
	return nil
}

func TestAuthorizeWithEmaidAndCertificateHashes(t *testing.T) {
	engine := inmemory.NewStore()
	err := setupTokenStore(engine)
	require.NoError(t, err)

	ah := handlers.AuthorizeHandler{
		TokenStore:                   engine,
		CertificateValidationService: mockCertValidationService{},
	}

	req := &types.AuthorizeRequestJson{
		IdToken: types.IdTokenType{
			Type:    types.IdTokenEnumTypeEMAID,
			IdToken: "MYEMAID",
		},
		ISO15118CertificateHashData: []types.OCSPRequestDataType{
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
		CertificateStatus: certStatus,
	}

	assert.Equal(t, want, got)
}

func TestAuthorizeWithEmaidAndInvalidCertificateHashes(t *testing.T) {
	engine := inmemory.NewStore()
	err := setupTokenStore(engine)
	require.NoError(t, err)

	ah := handlers.AuthorizeHandler{
		TokenStore:                   engine,
		CertificateValidationService: mockCertValidationService{},
	}

	testCases := []string{"invalidCertChain", "revokedCertChain", "expiredCertChain", "signatureError"}
	expectedErrors := []types.AuthorizeCertificateStatusEnumType{
		types.AuthorizeCertificateStatusEnumTypeCertificateRevoked,
		types.AuthorizeCertificateStatusEnumTypeCertificateRevoked,
		types.AuthorizeCertificateStatusEnumTypeCertificateRevoked,
		types.AuthorizeCertificateStatusEnumTypeCertificateRevoked,
	}

	for index, testCase := range testCases {
		t.Run(testCase, func(t *testing.T) {
			req := &types.AuthorizeRequestJson{
				IdToken: types.IdTokenType{
					Type:    types.IdTokenEnumTypeEMAID,
					IdToken: "MYEMAID",
				},
				ISO15118CertificateHashData: []types.OCSPRequestDataType{
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
				CertificateStatus: expectedErrors[index],
			}

			assert.Equal(t, want, got)
		})
	}
}

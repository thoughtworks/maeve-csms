// SPDX-License-Identifier: Apache-2.0

//go:build integration

package services_test

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	types "github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"github.com/thoughtworks/maeve-csms/manager/services"
	"net/http"
	"testing"
)

func TestCertificateValidationServiceWithHubjectCertificate(t *testing.T) {
	ctx := context.Background()
	httpClient := http.DefaultClient

	certSignerService := services.OpcpChargeStationCertificateProvider{
		HttpTokenService: services.NewHubjectTestHttpTokenService(
			"https://hubject.stoplight.io/api/v1/projects/cHJqOjk0NTg5/nodes/6bb8b3bc79c2e-authorization-token",
			http.DefaultClient),
		BaseURL:    "https://open.plugncharge-test.hubject.com",
		ISOVersion: services.ISO15118V2,
		HttpClient: httpClient,
	}

	csr := createCertificateSigningRequest(t)
	chain, err := certSignerService.ProvideCertificate(ctx, services.CertificateTypeV2G, string(csr), "")
	require.NoError(t, err)

	certificateValidationService := services.OnlineCertificateValidationService{
		RootCertificateProvider: services.OpcpRootCertificateProviderService{
			BaseURL: "https://open.plugncharge-test.hubject.com",
			TokenService: services.NewHubjectTestHttpTokenService(
				"https://hubject.stoplight.io/api/v1/projects/cHJqOjk0NTg5/nodes/6bb8b3bc79c2e-authorization-token",
				httpClient),
			HttpClient: httpClient,
		},
		MaxOCSPAttempts: 3,
		HttpClient:      httpClient,
	}

	ocspData, err := certificateValidationService.ValidatePEMCertificateChain(ctx, []byte(chain), "cs001")
	assert.NoError(t, err)
	assert.NotNil(t, ocspData)
}

func TestCertificateValidationServiceWithHubjectCertificateHashes(t *testing.T) {
	ctx := context.Background()
	httpClient := http.DefaultClient

	certSignerService := services.OpcpChargeStationCertificateProvider{
		HttpTokenService: services.NewHubjectTestHttpTokenService(
			"https://hubject.stoplight.io/api/v1/projects/cHJqOjk0NTg5/nodes/6bb8b3bc79c2e-authorization-token",
			httpClient),
		BaseURL:    "https://open.plugncharge-test.hubject.com",
		ISOVersion: services.ISO15118V2,
		HttpClient: httpClient,
	}

	csr := createCertificateSigningRequest(t)
	chain, err := certSignerService.ProvideCertificate(ctx, services.CertificateTypeV2G, string(csr), "")
	require.NoError(t, err)

	certificateValidationService := services.OnlineCertificateValidationService{
		RootCertificateProvider: services.OpcpRootCertificateProviderService{
			TokenService: services.NewHubjectTestHttpTokenService(
				"https://hubject.stoplight.io/api/v1/projects/cHJqOjk0NTg5/nodes/6bb8b3bc79c2e-authorization-token",
				httpClient),
			HttpClient: httpClient,
		},
		MaxOCSPAttempts: 3,
		HttpClient:      httpClient,
	}

	block, next := pem.Decode([]byte(chain))
	leaf, err := x509.ParseCertificate(block.Bytes)
	require.NoError(t, err)

	block, next = pem.Decode(next)
	issuer, err := x509.ParseCertificate(block.Bytes)
	require.NoError(t, err)

	pkey, err := getPublicKeyBytes(issuer.RawSubjectPublicKeyInfo)
	require.NoError(t, err)

	issuerKeyHash := hashBytes(pkey)
	issuerNameHash := hashBytes(issuer.RawSubject)
	serialNumber := leaf.SerialNumber.Text(16)

	ocspResp, err := certificateValidationService.ValidateHashedCertificateChain(ctx, []types.OCSPRequestDataType{
		{
			HashAlgorithm:  "SHA256",
			IssuerKeyHash:  issuerKeyHash,
			IssuerNameHash: issuerNameHash,
			ResponderURL:   leaf.OCSPServer[0],
			SerialNumber:   serialNumber,
		},
	})

	require.NoError(t, err)
	require.NotNil(t, ocspResp)
}

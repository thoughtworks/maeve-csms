//go:build integration

package services_test

import (
	"crypto/x509"
	"encoding/pem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	types "github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"github.com/thoughtworks/maeve-csms/manager/services"
	"os"
	"testing"
)

const v2gRootCertificate = `
subject=/C=DE/O=Hubject GmbH/DC=V2G/CN=V2G Root CA QA G1
issuer=/C=DE/O=Hubject GmbH/DC=V2G/CN=V2G Root CA QA G1
-----BEGIN CERTIFICATE-----
MIICUzCCAfmgAwIBAgIQaasA0lm730LOgFKa0wzl7TAKBggqhkjOPQQDAjBVMQsw
CQYDVQQGEwJERTEVMBMGA1UEChMMSHViamVjdCBHbWJIMRMwEQYKCZImiZPyLGQB
GRYDVjJHMRowGAYDVQQDExFWMkcgUm9vdCBDQSBRQSBHMTAgFw0xOTA0MjYwODM3
MTVaGA8yMDU5MDQyNjA4MzcxNVowVTELMAkGA1UEBhMCREUxFTATBgNVBAoTDEh1
YmplY3QgR21iSDETMBEGCgmSJomT8ixkARkWA1YyRzEaMBgGA1UEAxMRVjJHIFJv
b3QgQ0EgUUEgRzEwWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAAShT8kSNcC+74TN
D82On2Y2TOf8mYfxw73lKZ7t9cmEXHpMdAgsWBQ4LI+pOMhe6NOHzJbzP38kQTg4
zLfw3kU0o4GoMIGlMBMGA1UdJQQMMAoGCCsGAQUFBwMJMA8GA1UdEwEB/wQFMAMB
Af8wEQYDVR0OBAoECEtF/4Il/BCWMEUGA1UdIAQ+MDwwOgYMKwYBBAGCxDUBAgEA
MCowKAYIKwYBBQUHAgEWHGh0dHBzOi8vd3d3Lmh1YmplY3QuY29tL3BraS8wEwYD
VR0jBAwwCoAIS0X/giX8EJYwDgYDVR0PAQH/BAQDAgEGMAoGCCqGSM49BAMCA0gA
MEUCIQCq3Qx2BLYVFb7Lt5XXpSlUViYv4cIUOQE1Ce9o2Jyy1QIgZRmVzMVjHZA+
toiM000PCUrLppqbLpcRN4MP8kE0OhU=
-----END CERTIFICATE-----`

func TestCertificateValidationServiceWithHubjectCertificate(t *testing.T) {
	bearerToken, ok := os.LookupEnv("HUBJECT_TOKEN")
	if !ok {
		t.Fatal("no bearer token for Hubject API - set the HUBJECT_TOKEN environment variable")
	}

	certSignerService := services.HubjectCertificateSignerService{
		BearerToken: bearerToken,
		BaseURL:     "https://open.plugncharge-test.hubject.com",
		ISOVersion:  services.ISO15118V2,
	}

	csr := createCertificateSigningRequest(t)
	chain, err := certSignerService.SignCertificate(services.CertificateTypeV2G, string(csr))
	require.NoError(t, err)

	block, _ := pem.Decode([]byte(v2gRootCertificate))
	v2gRoot, err := x509.ParseCertificate(block.Bytes)
	require.NoError(t, err)

	certificateValidationService := services.OnlineCertificateValidationService{
		RootCertificates: []*x509.Certificate{v2gRoot},
		MaxOCSPAttempts:  3,
	}

	ocspData, err := certificateValidationService.ValidatePEMCertificateChain([]byte(chain), "cs001")
	assert.NoError(t, err)
	assert.NotNil(t, ocspData)
}

func TestCertificateValidationServiceWithHubjectCertificateHashes(t *testing.T) {
	bearerToken, ok := os.LookupEnv("HUBJECT_TOKEN")
	if !ok {
		t.Fatal("no bearer token for Hubject API - set the HUBJECT_TOKEN environment variable")
	}

	certSignerService := services.HubjectCertificateSignerService{
		BearerToken: bearerToken,
		BaseURL:     "https://open.plugncharge-test.hubject.com",
		ISOVersion:  services.ISO15118V2,
	}

	csr := createCertificateSigningRequest(t)
	chain, err := certSignerService.SignCertificate(services.CertificateTypeV2G, string(csr))
	require.NoError(t, err)

	block, _ := pem.Decode([]byte(v2gRootCertificate))
	v2gRoot, err := x509.ParseCertificate(block.Bytes)
	require.NoError(t, err)

	certificateValidationService := services.OnlineCertificateValidationService{
		RootCertificates: []*x509.Certificate{v2gRoot},
		MaxOCSPAttempts:  3,
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

	ocspResp, err := certificateValidationService.ValidateHashedCertificateChain([]types.OCSPRequestDataType{
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

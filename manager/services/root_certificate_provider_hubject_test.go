// SPDX-License-Identifier: Apache-2.0

//go:build integration

package services_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thoughtworks/maeve-csms/manager/services"
	"net/http"
	"os"
	"testing"
)

func TestProvideCertificatesWithHubjectRootCertificatePool(t *testing.T) {
	bearerToken, ok := os.LookupEnv("HUBJECT_TOKEN")
	if !ok {
		t.Fatal("no bearer token for Hubject API - set the HUBJECT_TOKEN environment variable")
	}

	rcp := &services.OpcpRootCertificateRetrieverService{
		BaseURL:         "https://open.plugncharge-test.hubject.com",
		HttpAuthService: services.NewFixedTokenHttpAuthService(bearerToken),
		HttpClient:      http.DefaultClient,
	}

	certificates, err := rcp.ProvideCertificates(context.Background())
	require.NoError(t, err)

	assert.Greater(t, len(certificates), 0)

	var certCommonNames []string
	for _, cert := range certificates {
		certCommonNames = append(certCommonNames, cert.Subject.CommonName)
	}
	assert.Contains(t, certCommonNames, "V2G Root CA QA G1")
}

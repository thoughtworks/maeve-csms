// SPDX-License-Identifier: Apache-2.0

//go:build integration

package services_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thoughtworks/maeve-csms/manager/services"
	"net/http"
	"testing"
)

func TestProvideCertificatesWithHubjectRootCertificatePool(t *testing.T) {
	rcp := &services.OpcpRootCertificateProviderService{
		BaseURL: "https://open.plugncharge-test.hubject.com",
		TokenService: services.NewHubjectTestHttpTokenService(
			"https://hubject.stoplight.io/api/v1/projects/cHJqOjk0NTg5/nodes/6bb8b3bc79c2e-authorization-token",
			http.DefaultClient),
		HttpClient: http.DefaultClient,
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

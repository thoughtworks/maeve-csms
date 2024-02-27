// SPDX-License-Identifier: Apache-2.0

package services_test

import (
	"context"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thoughtworks/maeve-csms/manager/services"
	"k8s.io/utils/clock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

const rcpUrl = "/v1/root/rootCerts"

type moRootCertificatePoolHandler struct {
	caCert *x509.Certificate
}

func (h moRootCertificatePoolHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(r.URL.String(), rcpUrl) {
		if r.Header.Get("authorization") != "Bearer Token" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		b64crt := base64.StdEncoding.EncodeToString(h.caCert.Raw)

		enc, err := json.Marshal(services.OpcpRootCertificateReturnType{
			RootCertificateCollection: services.OpcpRootCertificateCollectionType{
				RootCertificates: []services.OpcpRootCertificateType{
					{
						RootCertificateId: h.caCert.SerialNumber.String(),
						CommonName:        h.caCert.Subject.CommonName,
						OrganizationName:  h.caCert.Subject.Organization[0],
						DN:                h.caCert.Subject.String(),
						ValidFrom:         h.caCert.NotBefore.Format(time.RFC3339),
						ValidTo:           h.caCert.NotAfter.Format(time.RFC3339),
						CACertificate:     b64crt,
					},
				},
			},
		})
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		_, _ = w.Write(enc)
		return
	} else {
		http.NotFound(w, r)
	}
}

func TestFileMoRootCertificateProviderServiceRetrieveCertificates(t *testing.T) {
	fileRetrievalService := services.FileRootCertificateProviderService{
		FilePaths: []string{"testdata/root_ca.pem"},
	}

	certs, e := fileRetrievalService.ProvideCertificates(context.TODO())
	assert.NoError(t, e)
	assert.Equal(t, "V2G Root CA QA G1", certs[0].Issuer.CommonName)
}

func TestOpcpMoRootCertificateProviderServiceRetrieveCertificatesNotAuthorised(t *testing.T) {
	rootCA, _ := createRootCACertificate(t, "V2G Root CA QA G1")
	handler := moRootCertificatePoolHandler{caCert: rootCA}
	server := httptest.NewServer(handler)
	service := services.OpcpRootCertificateProviderService{
		TokenService: services.NewFixedHttpTokenService("BadToken"),
		BaseURL:      server.URL,
		HttpClient:   http.DefaultClient,
	}

	defer server.Close()

	_, err := service.ProvideCertificates(context.TODO())

	var httpError services.HttpError
	assert.ErrorAs(t, err, &httpError)
	assert.Equal(t, http.StatusUnauthorized, int(httpError))
}

func TestOpcpMoRootCertificateProviderServiceRetrieveCertificates(t *testing.T) {
	rootCA, _ := createRootCACertificate(t, "V2G Root CA QA G1")
	handler := moRootCertificatePoolHandler{caCert: rootCA}
	server := httptest.NewServer(handler)
	service := services.OpcpRootCertificateProviderService{
		TokenService: services.NewFixedHttpTokenService("Token"),
		BaseURL:      server.URL,
		HttpClient:   http.DefaultClient,
	}
	defer server.Close()

	result, err := service.ProvideCertificates(context.TODO())

	require.NoError(t, err)

	assert.Equal(t, "V2G Root CA QA G1", result[0].Issuer.CommonName)
}

type CountingRootCertificateProviderService struct {
	Count        int
	Certificates []*x509.Certificate
}

func (c *CountingRootCertificateProviderService) ProvideCertificates(context.Context) ([]*x509.Certificate, error) {
	c.Count++
	return c.Certificates, nil
}

func TestCachingRootCertificateProviderService(t *testing.T) {
	rootCA, _ := createRootCACertificate(t, "V2G Root CA QA G1")
	svc := &CountingRootCertificateProviderService{
		Certificates: []*x509.Certificate{rootCA},
	}
	cachingService := services.NewCachingRootCertificateProviderService(svc, 100*time.Millisecond, &clock.RealClock{})

	certs, err := cachingService.ProvideCertificates(context.TODO())
	require.NoError(t, err)
	assert.Equal(t, 1, svc.Count)
	assert.Equal(t, "V2G Root CA QA G1", certs[0].Issuer.CommonName)

	certs, err = cachingService.ProvideCertificates(context.TODO())
	require.NoError(t, err)
	assert.Equal(t, 1, svc.Count)
	assert.Equal(t, "V2G Root CA QA G1", certs[0].Issuer.CommonName)

	time.Sleep(100 * time.Millisecond)

	certs, err = cachingService.ProvideCertificates(context.TODO())
	require.NoError(t, err)
	assert.Equal(t, 2, svc.Count)
	assert.Equal(t, "V2G Root CA QA G1", certs[0].Issuer.CommonName)
}

func TestCompositeRootCertificateProviderService(t *testing.T) {
	fileRetrievalService := services.FileRootCertificateProviderService{
		FilePaths: []string{"testdata/root_ca.pem"},
	}
	rootCA, _ := createRootCACertificate(t, "V2G Root CA QA G2")
	handler := moRootCertificatePoolHandler{caCert: rootCA}
	server := httptest.NewServer(handler)
	opcpService := services.OpcpRootCertificateProviderService{
		TokenService: services.NewFixedHttpTokenService("Token"),
		BaseURL:      server.URL,
		HttpClient:   http.DefaultClient,
	}

	defer server.Close()

	compositeService := services.CompositeRootCertificateProviderService{
		Providers: []services.RootCertificateProviderService{fileRetrievalService, opcpService},
	}

	certs, err := compositeService.ProvideCertificates(context.TODO())
	require.NoError(t, err)
	assert.Equal(t, 2, len(certs))
	assert.Equal(t, "V2G Root CA QA G1", certs[0].Issuer.CommonName)
	assert.Equal(t, "V2G Root CA QA G2", certs[1].Issuer.CommonName)
}

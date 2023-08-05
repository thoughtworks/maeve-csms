package services

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

const rcpUrl = "/v1/root/rootCerts"

type moRootCertificatePoolHandler struct{}

func newHandler() moRootCertificatePoolHandler {
	return moRootCertificatePoolHandler{}
}

func (h moRootCertificatePoolHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(r.URL.String(), rcpUrl) {
		if r.Header.Get("authorization") != "Bearer Token" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		crt := createRootCACertificate("V2G Root CA QA G1")
		x5c, err := x509.ParseCertificate(crt)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		b64crt := base64.StdEncoding.EncodeToString(crt)

		enc, err := json.Marshal(OpcpRootCertificateReturnType{
			RootCertificateCollection: OpcpRootCertificateCollectionType{
				RootCertificates: []OpcpRootCertificateType{
					{
						RootCertificateId: x5c.SerialNumber.String(),
						CommonName:        x5c.Subject.CommonName,
						OrganizationName:  x5c.Subject.Organization[0],
						DN:                x5c.Subject.String(),
						ValidFrom:         x5c.NotBefore.Format(time.RFC3339),
						ValidTo:           x5c.NotAfter.Format(time.RFC3339),
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

func TestFileMoRootCertificateRetrievalServiceRetrieveCertificates(t *testing.T) {
	fileRetrievalService := FileRootCertificateRetrieverService{
		FilePaths: []string{"testdata/root_ca.pem"},
	}

	certs, e := fileRetrievalService.ProvideCertificates(context.TODO())
	assert.NoError(t, e)
	assert.Equal(t, "V2G Root CA QA G1", certs[0].Issuer.CommonName)
}

func TestOpcpMoRootCertificateRetrievalServiceRetrieveCertificatesNotAuthorised(t *testing.T) {
	handler := moRootCertificatePoolHandler{}
	server := httptest.NewServer(handler)
	service := OpcpRootCertificateRetrieverService{
		HttpAuthService: NewFixedTokenHttpAuthService("BadToken"),
		BaseURL:         server.URL,
		HttpClient:      http.DefaultClient,
	}

	defer server.Close()

	_, err := service.ProvideCertificates(context.TODO())

	var httpError HttpError
	assert.ErrorAs(t, err, &httpError)
	assert.Equal(t, http.StatusUnauthorized, int(httpError))
}

func TestOpcpMoRootCertificateRetrievalServiceRetrieveCertificates(t *testing.T) {
	handler := newHandler()
	server := httptest.NewServer(handler)
	service := OpcpRootCertificateRetrieverService{
		HttpAuthService: NewFixedTokenHttpAuthService("Token"),
		BaseURL:         server.URL,
		HttpClient:      http.DefaultClient,
	}
	defer server.Close()

	result, err := service.ProvideCertificates(context.TODO())

	require.NoError(t, err)

	assert.Equal(t, "V2G Root CA QA G1", result[0].Issuer.CommonName)
}

func createRootCACertificate(commonName string) []byte {
	caKeyPair, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	caCertTemplate := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName:   commonName,
			Organization: []string{"Thoughtworks"},
		},
	}

	crt, _ := x509.CreateCertificate(rand.Reader, &caCertTemplate, &caCertTemplate, &caKeyPair.PublicKey, caKeyPair)
	return crt
}

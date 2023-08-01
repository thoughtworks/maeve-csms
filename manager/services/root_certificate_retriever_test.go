package services

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mozilla.org/pkcs7"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const baseUrl = "/mo/cacerts/ISO15118-2"

type DummyFileReader struct{}

func (reader DummyFileReader) ReadFile(_ string) (content []byte, e error) {
	derBytes := createRootCACertificate("Thoughtworks")
	pemData := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: derBytes,
	})
	return pemData, nil
}

type moRootCertificatePoolHandler struct {
	baseUrl string
}

func newHandler() moRootCertificatePoolHandler {
	return moRootCertificatePoolHandler{
		baseUrl: baseUrl,
	}
}

func (h moRootCertificatePoolHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(r.URL.String(), baseUrl) {
		if r.Header.Get("authorization") != "Bearer Token" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		crt := createRootCACertificate("V2G Root CA QA G1")
		p7, _ := pkcs7.DegenerateCertificate(crt)

		enc := make([]byte, base64.StdEncoding.EncodedLen(len(p7)))
		base64.StdEncoding.Encode(enc, p7)

		w.Write(enc)
		return
	}
}

func TestFileMoRootCertificateRetrievalServiceRetrieveCertificates(t *testing.T) {
	fileRetrievalService := FileRootCertificateRetrieverService{
		FilePaths:  []string{"test"},
		FileReader: DummyFileReader{},
	}

	certs, e := fileRetrievalService.RetrieveCertificates(context.TODO())
	assert.NoError(t, e)
	assert.Equal(t, "Thoughtworks", certs[0].Issuer.CommonName)
}

func TestOpcpMoRootCertificateRetrievalServiceRetrieveCertificatesNotAuthorised(t *testing.T) {
	handler := moRootCertificatePoolHandler{baseUrl: baseUrl}
	server := httptest.NewServer(handler)
	service := OpcpRootCertificateRetrieverService{
		MoOPCPToken:    "",
		MoRootCertPool: server.URL + baseUrl,
		HttpClient:     http.DefaultClient,
	}

	defer server.Close()

	_, err := service.RetrieveCertificates(context.TODO())

	var httpError HttpError
	assert.ErrorAs(t, err, &httpError)
	assert.Equal(t, http.StatusUnauthorized, int(httpError))
}

func TestOpcpMoRootCertificateRetrievalServiceRetrieveCertificates(t *testing.T) {
	handler := newHandler()
	server := httptest.NewServer(handler)
	service := OpcpRootCertificateRetrieverService{
		MoOPCPToken:    "Token",
		MoRootCertPool: server.URL + baseUrl,
		HttpClient:     http.DefaultClient,
	}
	defer server.Close()

	result, err := service.RetrieveCertificates(context.TODO())

	require.NoError(t, err)

	assert.Equal(t, "V2G Root CA QA G1", result[0].Issuer.CommonName)
}

func createRootCACertificate(commonName string) []byte {
	caKeyPair, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	caCertTemplate := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName: commonName,
		},
	}

	crt, _ := x509.CreateCertificate(rand.Reader, &caCertTemplate, &caCertTemplate, &caKeyPair.PublicKey, caKeyPair)
	return crt
}

package cmd

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
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

func TestPullMoTrustAnchorsNotAuthorised(t *testing.T) {
	handler := moRootCertificatePoolHandler{baseUrl: baseUrl}
	server := httptest.NewServer(handler)
	defer server.Close()

	_, err := pullMoTrustAnchors(server.URL + baseUrl)

	assert.Equal(t, "received code 401 in response", err.Error())
}

func TestPullMoTrustAnchors(t *testing.T) {
	handler := newHandler()
	moOPCPToken = "Token"
	server := httptest.NewServer(handler)
	defer server.Close()

	result, err := pullMoTrustAnchors(server.URL + baseUrl)

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

	caCertBytes, _ := x509.CreateCertificate(rand.Reader, &caCertTemplate, &caCertTemplate, &caKeyPair.PublicKey, caKeyPair)
	return caCertBytes
}

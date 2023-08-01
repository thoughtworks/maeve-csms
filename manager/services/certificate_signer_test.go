// SPDX-License-Identifier: Apache-2.0

package services_test

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thoughtworks/maeve-csms/manager/services"
	"go.mozilla.org/pkcs7"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type opcpHttpHandler struct {
	caCert  *x509.Certificate
	caKey   *ecdsa.PrivateKey
	intCert *x509.Certificate
	intKey  *ecdsa.PrivateKey
}

func newOPCPHttpHandler(t *testing.T) opcpHttpHandler {
	caCert, caKey := createRootCACertificate(t, "test")
	intCert, intKey := createIntermediateCACertificate(t, "int", "", caCert, caKey)

	return opcpHttpHandler{
		caCert:  caCert,
		caKey:   caKey,
		intCert: intCert,
		intKey:  intKey,
	}
}

func (h opcpHttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.String() == "/cpo/simpleenroll/ISO15118-2" {
		if r.Header.Get("authorization") != "Bearer TestToken" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		if r.Header.Get("content-type") != "application/pkcs10" {
			http.Error(w, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
			return
		}

		enc, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Printf("reading body: %v", err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		defer func() {
			_ = r.Body.Close()
		}()

		csr := make([]byte, base64.StdEncoding.DecodedLen(len(enc)))
		n, err := base64.StdEncoding.Decode(csr, enc)
		if err != nil {
			fmt.Printf("decoding body: %v", err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		cert, err := signCertificateRequest(csr[:n], h.intCert, h.intKey)
		if err != nil {
			fmt.Printf("signing certificate request: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		p7, err := pkcs7.DegenerateCertificate(cert)
		if err != nil {
			fmt.Printf("degenerate certificate: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		enc = make([]byte, base64.StdEncoding.EncodedLen(len(p7)))
		base64.StdEncoding.Encode(enc, p7)

		_, _ = w.Write(enc)

		return
	} else if r.URL.String() == "/cpo/cacerts/ISO15118-2" {
		if r.Header.Get("authorization") != "Bearer TestToken" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		var certBytes []byte

		certBytes = append(certBytes, h.intCert.Raw...)
		certBytes = append(certBytes, h.caCert.Raw...)

		p7, err := pkcs7.DegenerateCertificate(certBytes)
		if err != nil {
			fmt.Printf("degenerate certificate: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		enc := make([]byte, base64.StdEncoding.EncodedLen(len(p7)))
		base64.StdEncoding.Encode(enc, p7)

		_, _ = w.Write(enc)

		return
	}

	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

func TestOPCPCertificateSignerService(t *testing.T) {
	opcp := newOPCPHttpHandler(t)

	server := httptest.NewServer(opcp)
	defer server.Close()

	opcpSigner := services.OpcpCpoCertificateSignerService{
		BaseURL:     server.URL,
		BearerToken: "TestToken",
		ISOVersion:  services.ISO15118V2,
	}

	csr := createCertificateSigningRequest(t)

	pemChain, err := opcpSigner.SignCertificate(context.TODO(), services.CertificateTypeV2G, string(csr))
	require.NoError(t, err)

	pemBytes := []byte(pemChain)
	block, pemBytes := pem.Decode(pemBytes)
	leafCert, err := x509.ParseCertificate(block.Bytes)
	require.NoError(t, err)
	assert.Equal(t, "CN=cs001,O=Thoughtworks", leafCert.Subject.String())

	require.NotNil(t, pemBytes)
	block, pemBytes = pem.Decode(pemBytes)
	intCert, err := x509.ParseCertificate(block.Bytes)
	require.NoError(t, err)
	assert.Equal(t, "CN=int,O=Thoughtworks", intCert.Subject.String())

	require.NotNil(t, pemBytes)
	block, _ = pem.Decode(pemBytes)
	require.Nil(t, block)
}

func createCertificateSigningRequest(t *testing.T) []byte {
	keyPair, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)

	csrTemplate := x509.CertificateRequest{
		Subject: pkix.Name{
			CommonName:   "cs001",
			Organization: []string{"Thoughtworks"},
		},
		SignatureAlgorithm: x509.ECDSAWithSHA256,
	}

	csrCertificate, err := x509.CreateCertificateRequest(rand.Reader, &csrTemplate, keyPair)
	require.NoError(t, err)

	return pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: csrCertificate,
	})
}

func signCertificateRequest(csrBytes []byte, caCert *x509.Certificate, caKey *ecdsa.PrivateKey) ([]byte, error) {
	serial, err := rand.Int(rand.Reader, (&big.Int{}).Exp(big.NewInt(2), big.NewInt(159), nil))
	if err != nil {
		return nil, fmt.Errorf("creating serial number: %w", err)
	}

	csr, err := x509.ParseCertificateRequest(csrBytes)
	if err != nil {
		return nil, fmt.Errorf("parsing certificate request: %w", err)
	}

	now := time.Now()
	template := x509.Certificate{
		SerialNumber:          serial,
		Subject:               csr.Subject,
		NotBefore:             now,
		NotAfter:              now.Add(time.Minute),
		BasicConstraintsValid: true,
		IsCA:                  false,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
	}

	return x509.CreateCertificate(rand.Reader, &template, caCert, csr.PublicKey, caKey)
}

// SPDX-License-Identifier: Apache-2.0

package services_test

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"github.com/thoughtworks/maeve-csms/manager/services"
	"golang.org/x/crypto/ocsp"
	"io"
	"k8s.io/utils/strings/slices"
	"math"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type OCSPResponder struct {
	T                    *testing.T
	Certificates         []*x509.Certificate
	Keys                 []*ecdsa.PrivateKey
	RevokedSerialNumbers []string
}

func (o *OCSPResponder) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reqBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	defer func() {
		_ = r.Body.Close()
	}()
	ocspReq, err := ocsp.ParseRequest(reqBytes)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	o.T.Logf("checking certificate with serial number %v", ocspReq.SerialNumber)

	var subjectCert *x509.Certificate
	for _, cert := range o.Certificates {
		if cert.SerialNumber.Text(16) == ocspReq.SerialNumber.Text(16) {
			subjectCert = cert
		}
	}
	if subjectCert == nil {
		o.T.Logf("can't find subject certificate")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	o.T.Logf("subject certificate: %s", subjectCert.Subject.String())

	var issuerCert *x509.Certificate
	var issuerKey *ecdsa.PrivateKey
	for index, cert := range o.Certificates {
		if cert.Subject.String() == subjectCert.Issuer.String() {
			issuerCert = cert
			issuerKey = o.Keys[index]
		}
	}
	if issuerCert == nil {
		o.T.Logf("can't find issuer certificate")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	o.T.Logf("issuer certificate: %s", issuerCert.Subject.String())

	now := time.Now()
	var ocspResp ocsp.Response
	if slices.Contains(o.RevokedSerialNumbers, subjectCert.SerialNumber.Text(16)) {
		ocspResp = ocsp.Response{
			Status:           ocsp.Revoked,
			SerialNumber:     ocspReq.SerialNumber,
			ThisUpdate:       now,
			NextUpdate:       now,
			RevokedAt:        now,
			RevocationReason: ocsp.PrivilegeWithdrawn,
		}
	} else {
		ocspResp = ocsp.Response{
			Status:       ocsp.Good,
			SerialNumber: ocspReq.SerialNumber,
			ThisUpdate:   now,
			NextUpdate:   now,
		}
	}

	respBytes, err := ocsp.CreateResponse(issuerCert, issuerCert, ocspResp, issuerKey)
	if err != nil {
		o.T.Logf("creating ocsp response: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(respBytes)
}

func setupOCSPResponder(t *testing.T, serverUrl string, ocspResponder *OCSPResponder) ([]*x509.Certificate, *x509.Certificate, *x509.Certificate) {
	ca1Cert, ca1Key := createRootCACertificate(t, "ca1")
	ca2Cert, ca2Key := createRootCACertificate(t, "ca2")
	intCACert, intCAKey := createIntermediateCACertificate(t, "int1", serverUrl, ca2Cert, ca2Key)
	leafCert, _ := createLeafCertificate(t, "MYEMAID", serverUrl, intCACert, intCAKey)

	ocspResponder.Certificates = append(ocspResponder.Certificates, ca1Cert, ca2Cert, intCACert, leafCert)
	ocspResponder.Keys = append(ocspResponder.Keys, ca1Key, ca2Key, intCAKey, nil)

	return []*x509.Certificate{ca1Cert, ca2Cert}, intCACert, leafCert
}

func TestValidatingPEMCertificateChain(t *testing.T) {
	ocspResponder := &OCSPResponder{
		T: t,
	}

	server := httptest.NewServer(ocspResponder)
	defer server.Close()

	rootCACerts, intCACert, leafCert := setupOCSPResponder(t, server.URL, ocspResponder)

	validationService := services.OnlineCertificateValidationService{
		RootCertificates: rootCACerts,
		MaxOCSPAttempts:  3,
		HttpClient:       http.DefaultClient,
	}

	pemChain := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: leafCert.Raw,
	})
	pemChain = append(pemChain, pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: intCACert.Raw,
	})...)

	ocspResp, err := validationService.ValidatePEMCertificateChain(context.TODO(), pemChain, "MYEMAID")
	require.NoError(t, err)

	validateOCSPResponse(t, ocspResp)
}

func validateOCSPResponse(t *testing.T, ocspResp *string) {
	assert.NotNil(t, ocspResp)
	respBytes, err := base64.StdEncoding.DecodeString(*ocspResp)
	require.NoError(t, err)
	_, err = ocsp.ParseResponse(respBytes, nil)
	require.NoError(t, err)
}

func TestValidatingPEMCertificateChainWithRevokedCertificate(t *testing.T) {
	ocspResponder := &OCSPResponder{
		T: t,
	}

	server := httptest.NewServer(ocspResponder)
	defer server.Close()

	rootCACerts, intCACert, leafCert := setupOCSPResponder(t, server.URL, ocspResponder)
	ocspResponder.RevokedSerialNumbers = []string{intCACert.SerialNumber.Text(16)}

	validationService := services.OnlineCertificateValidationService{
		RootCertificates: rootCACerts,
		MaxOCSPAttempts:  3,
		HttpClient:       http.DefaultClient,
	}

	pemChain := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: leafCert.Raw,
	})
	pemChain = append(pemChain, pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: intCACert.Raw,
	})...)

	ocspResp, err := validationService.ValidatePEMCertificateChain(context.TODO(), pemChain, "MYEMAID")
	require.Error(t, err)
	assert.Nil(t, ocspResp)
}

func TestValidatingPEMCertificateChainWithWrongEmaid(t *testing.T) {
	ocspResponder := &OCSPResponder{
		T: t,
	}

	server := httptest.NewServer(ocspResponder)
	defer server.Close()

	rootCACerts, intCACert, leafCert := setupOCSPResponder(t, server.URL, ocspResponder)

	validationService := services.OnlineCertificateValidationService{
		RootCertificates: rootCACerts,
		MaxOCSPAttempts:  3,
		HttpClient:       http.DefaultClient,
	}

	pemChain := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: leafCert.Raw,
	})
	pemChain = append(pemChain, pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: intCACert.Raw,
	})...)

	ocspResp, err := validationService.ValidatePEMCertificateChain(context.TODO(), pemChain, "MYWRONGEMAID")
	require.Error(t, err)
	assert.Nil(t, ocspResp)
}

func TestValidatingPEMCertificateChainInvalidChain(t *testing.T) {
	ocspResponder := &OCSPResponder{
		T: t,
	}

	server := httptest.NewServer(ocspResponder)
	defer server.Close()

	rootCACerts, _, leafCert := setupOCSPResponder(t, server.URL, ocspResponder)

	validationService := services.OnlineCertificateValidationService{
		RootCertificates: rootCACerts,
		MaxOCSPAttempts:  3,
		HttpClient:       http.DefaultClient,
	}

	pemChain := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: leafCert.Raw,
	})

	ocspResp, err := validationService.ValidatePEMCertificateChain(context.TODO(), pemChain, "MYEMAID")
	require.Error(t, err)
	assert.Nil(t, ocspResp)
}

func TestValidatingPEMCertificateChainIncludingRootCertificate(t *testing.T) {
	ocspResponder := &OCSPResponder{
		T: t,
	}

	server := httptest.NewServer(ocspResponder)
	defer server.Close()

	rootCACerts, intCACert, leafCert := setupOCSPResponder(t, server.URL, ocspResponder)

	validationService := services.OnlineCertificateValidationService{
		RootCertificates: rootCACerts,
		MaxOCSPAttempts:  3,
		HttpClient:       http.DefaultClient,
	}

	pemChain := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: leafCert.Raw,
	})
	pemChain = append(pemChain, pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: intCACert.Raw,
	})...)
	pemChain = append(pemChain, pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: rootCACerts[1].Raw,
	})...)

	ocspResp, err := validationService.ValidatePEMCertificateChain(context.TODO(), pemChain, "MYEMAID")
	require.NoError(t, err)

	validateOCSPResponse(t, ocspResp)
}

func TestValidatingPEMCertificateChainIncludingUntrustedRootCertificate(t *testing.T) {
	ocspResponder := &OCSPResponder{
		T: t,
	}

	server := httptest.NewServer(ocspResponder)
	defer server.Close()

	rootCACerts, intCACert, leafCert := setupOCSPResponder(t, server.URL, ocspResponder)

	validationService := services.OnlineCertificateValidationService{
		RootCertificates: []*x509.Certificate{rootCACerts[0]},
		MaxOCSPAttempts:  3,
		HttpClient:       http.DefaultClient,
	}

	pemChain := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: leafCert.Raw,
	})
	pemChain = append(pemChain, pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: intCACert.Raw,
	})...)
	pemChain = append(pemChain, pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: rootCACerts[1].Raw,
	})...)

	ocspResp, err := validationService.ValidatePEMCertificateChain(context.TODO(), pemChain, "MYEMAID")
	var valErr services.ValidationError
	assert.Nil(t, ocspResp)
	require.ErrorAs(t, err, &valErr)
	assert.Equal(t, services.ValidationErrorCertChain, valErr)
}

func TestValidatingPEMCertificateChainWithNoOCSPOnLeafCertificate(t *testing.T) {
	ocspResponder := &OCSPResponder{
		T: t,
	}

	server := httptest.NewServer(ocspResponder)
	defer server.Close()

	rootCACerts, intCACert, leafCert := setupOCSPResponder(t, "", ocspResponder)

	validationService := services.OnlineCertificateValidationService{
		RootCertificates: rootCACerts,
		MaxOCSPAttempts:  3,
		HttpClient:       http.DefaultClient,
	}

	pemChain := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: leafCert.Raw,
	})
	pemChain = append(pemChain, pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: intCACert.Raw,
	})...)

	ocspResp, err := validationService.ValidatePEMCertificateChain(context.TODO(), pemChain, "MYEMAID")
	var valErr services.ValidationError
	assert.Nil(t, ocspResp)
	t.Log(err)
	require.ErrorAs(t, err, &valErr)
	assert.Equal(t, services.ValidationErrorCertChain, valErr)
}

func TestValidatingHashedCertificateChain(t *testing.T) {
	ocspResponder := &OCSPResponder{
		T: t,
	}

	server := httptest.NewServer(ocspResponder)
	defer server.Close()

	rootCACerts, intCACert, leafCert := setupOCSPResponder(t, server.URL, ocspResponder)

	validationService := services.OnlineCertificateValidationService{
		RootCertificates: rootCACerts,
		MaxOCSPAttempts:  3,
		HttpClient:       http.DefaultClient,
	}

	intCAPublicKeyBytes, err := getPublicKeyBytes(intCACert.RawSubjectPublicKeyInfo)
	require.NoError(t, err)
	ca2PublicKeyBytes, err := getPublicKeyBytes(rootCACerts[1].RawSubjectPublicKeyInfo)
	require.NoError(t, err)

	ocspResp, err := validationService.ValidateHashedCertificateChain(context.TODO(), []ocpp201.OCSPRequestDataType{
		{
			HashAlgorithm:  "SHA256",
			IssuerNameHash: hashBytes(leafCert.RawIssuer),
			IssuerKeyHash:  hashBytes(intCAPublicKeyBytes),
			ResponderURL:   leafCert.OCSPServer[0],
			SerialNumber:   leafCert.SerialNumber.Text(16),
		},
		{
			HashAlgorithm:  "SHA256",
			IssuerNameHash: hashBytes(intCACert.RawIssuer),
			IssuerKeyHash:  hashBytes(ca2PublicKeyBytes),
			ResponderURL:   intCACert.OCSPServer[0],
			SerialNumber:   intCACert.SerialNumber.Text(16),
		},
	})
	require.NoError(t, err)

	validateOCSPResponse(t, ocspResp)
}

func TestValidatingHashedCertificateChainWithRevokedCertificate(t *testing.T) {
	ocspResponder := &OCSPResponder{
		T: t,
	}

	server := httptest.NewServer(ocspResponder)
	defer server.Close()

	rootCACerts, intCACert, leafCert := setupOCSPResponder(t, server.URL, ocspResponder)
	ocspResponder.RevokedSerialNumbers = []string{intCACert.SerialNumber.String()}

	validationService := services.OnlineCertificateValidationService{
		RootCertificates: rootCACerts,
		MaxOCSPAttempts:  3,
		HttpClient:       http.DefaultClient,
	}

	intCAPublicKeyBytes, err := getPublicKeyBytes(intCACert.RawSubjectPublicKeyInfo)
	require.NoError(t, err)
	ca2PublicKeyBytes, err := getPublicKeyBytes(rootCACerts[1].RawSubjectPublicKeyInfo)
	require.NoError(t, err)

	ocspResp, err := validationService.ValidateHashedCertificateChain(context.TODO(), []ocpp201.OCSPRequestDataType{
		{
			HashAlgorithm:  "SHA256",
			IssuerNameHash: hashBytes(leafCert.RawIssuer),
			IssuerKeyHash:  hashBytes(intCAPublicKeyBytes),
			ResponderURL:   leafCert.OCSPServer[0],
			SerialNumber:   leafCert.SerialNumber.String(),
		},
		{
			HashAlgorithm:  "SHA256",
			IssuerNameHash: hashBytes(intCACert.RawIssuer),
			IssuerKeyHash:  hashBytes(ca2PublicKeyBytes),
			ResponderURL:   intCACert.OCSPServer[0],
			SerialNumber:   intCACert.SerialNumber.String(),
		},
	})
	require.Error(t, err)
	assert.Nil(t, ocspResp)
}

func getPublicKeyBytes(rawSubjectPublicKeyInfo []byte) ([]byte, error) {
	var publicKeyInfo struct {
		Algorithm pkix.AlgorithmIdentifier
		PublicKey asn1.BitString
	}
	if _, err := asn1.Unmarshal(rawSubjectPublicKeyInfo, &publicKeyInfo); err != nil {
		return nil, err
	}

	return publicKeyInfo.PublicKey.RightAlign(), nil
}

func hashBytes(b []byte) string {
	shaBytes := sha256.Sum256(b)
	return hex.EncodeToString(shaBytes[:])
}

func createRootCACertificate(t *testing.T, commonName string) (*x509.Certificate, *ecdsa.PrivateKey) {
	caKeyPair, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generating root CA key pair: %v", err)
	}

	caSerialNumber, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		t.Fatalf("generating root CA serial number: %v", err)
	}

	caCertTemplate := x509.Certificate{
		SerialNumber: caSerialNumber,
		Subject: pkix.Name{
			CommonName:   commonName,
			Organization: []string{"Thoughtworks"},
		},
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Minute),
	}

	caCertBytes, err :=
		x509.CreateCertificate(rand.Reader, &caCertTemplate, &caCertTemplate, &caKeyPair.PublicKey, caKeyPair)
	if err != nil {
		t.Fatalf("creating root CA certificate: %v", err)
	}
	caCert, err := x509.ParseCertificate(caCertBytes)
	if err != nil {
		t.Fatalf("parsing root CA certificate: %v", err)
	}

	return caCert, caKeyPair
}

func createIntermediateCACertificate(t *testing.T, commonName, ocspResponderUrl string, rootCA *x509.Certificate, rootKey *ecdsa.PrivateKey) (*x509.Certificate, *ecdsa.PrivateKey) {
	caKeyPair, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generating intermediate CA key pair: %v", err)
	}

	caSerialNumber, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		t.Fatalf("generating intermediate CA serial number: %v", err)
	}

	caCertTemplate := x509.Certificate{
		SerialNumber: caSerialNumber,
		Subject: pkix.Name{
			CommonName:   commonName,
			Organization: []string{"Thoughtworks"},
		},
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Minute),
	}
	if ocspResponderUrl != "" {
		caCertTemplate.OCSPServer = []string{ocspResponderUrl}
	}

	caCertBytes, err :=
		x509.CreateCertificate(rand.Reader, &caCertTemplate, rootCA, &caKeyPair.PublicKey, rootKey)
	if err != nil {
		t.Fatalf("creating intermediate CA certificate: %v", err)
	}
	caCert, err := x509.ParseCertificate(caCertBytes)
	if err != nil {
		t.Fatalf("parsing intermediate CA certificate: %v", err)
	}

	return caCert, caKeyPair
}

func createLeafCertificate(t *testing.T, commonName, ocspResponderUrl string, caCert *x509.Certificate, caKey *ecdsa.PrivateKey) (*x509.Certificate, *ecdsa.PrivateKey) {
	leafKeyPair, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generating leaf certificate key pair: %v", err)
	}

	caSerialNumber, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		t.Fatalf("generating leaf certificate serial number: %v", err)
	}

	leafCertTemplate := x509.Certificate{
		SerialNumber: caSerialNumber,
		Subject: pkix.Name{
			CommonName:   commonName,
			Organization: []string{"Thoughtworks"},
		},
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
		IsCA:                  false,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Minute),
	}
	if ocspResponderUrl != "" {
		leafCertTemplate.OCSPServer = []string{ocspResponderUrl}
	}

	leafCertBytes, err :=
		x509.CreateCertificate(rand.Reader, &leafCertTemplate, caCert, &leafKeyPair.PublicKey, caKey)
	if err != nil {
		t.Fatalf("creating leaf certificate: %v", err)
	}
	leafCert, err := x509.ParseCertificate(leafCertBytes)
	if err != nil {
		t.Fatalf("parsing leaf certificate: %v", err)
	}

	return leafCert, leafKeyPair
}

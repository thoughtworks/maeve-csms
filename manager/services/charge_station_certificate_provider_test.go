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
	"encoding/pem"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thoughtworks/maeve-csms/manager/services"
	"github.com/thoughtworks/maeve-csms/manager/store/inmemory"
	"go.mozilla.org/pkcs7"
	"io"
	"k8s.io/utils/clock"
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

func TestOPCPChargeStationCertificateProvider(t *testing.T) {
	opcp := newOPCPHttpHandler(t)

	server := httptest.NewServer(opcp)
	defer server.Close()

	provider := services.OpcpChargeStationCertificateProvider{
		BaseURL:          server.URL,
		HttpTokenService: services.NewFixedHttpTokenService("TestToken"),
		ISOVersion:       services.ISO15118V2,
		HttpClient:       http.DefaultClient,
	}

	csr := createCertificateSigningRequest(t)

	pemChain, err := provider.ProvideCertificate(context.TODO(), services.CertificateTypeV2G, string(csr), "cs001")
	require.NoError(t, err)

	pemBytes := []byte(pemChain)
	block, pemBytes := pem.Decode(pemBytes)
	leafCert, err := x509.ParseCertificate(block.Bytes)
	require.NoError(t, err)
	assert.Equal(t, "CN=cs001,O=Thoughtworks,C=GB", leafCert.Subject.String())

	require.NotNil(t, pemBytes)
	block, pemBytes = pem.Decode(pemBytes)
	intCert, err := x509.ParseCertificate(block.Bytes)
	require.NoError(t, err)
	assert.Equal(t, "CN=int,O=Thoughtworks", intCert.Subject.String())

	require.NotNil(t, pemBytes)
	block, _ = pem.Decode(pemBytes)
	require.Nil(t, block)
}

//	func TestOPCPChargeStationCertificateProviderWithWrongId(t *testing.T) {
//		opcp := newOPCPHttpHandler(t)
//
//		server := httptest.NewServer(opcp)
//		defer server.Close()
//
//		provider := services.OpcpChargeStationCertificateProvider{
//			BaseURL:          server.URL,
//			HttpTokenService: services.NewFixedHttpTokenService("TestToken"),
//			ISOVersion:       services.ISO15118V2,
//			HttpClient:       http.DefaultClient,
//		}
//
//		csr := createCertificateSigningRequest(t)
//
//		pemChain, err := provider.ProvideCertificate(context.TODO(), services.CertificateTypeV2G, string(csr), "not-cs001")
//		assert.Error(t, err)
//		assert.Equal(t, "", pemChain)
//	}
func TestLocalChargeStationCertificateProvider(t *testing.T) {
	store := inmemory.NewStore(clock.RealClock{})

	caCert, caKey := createRootCACertificate(t, "test")
	intCert, intKey := createIntermediateCACertificate(t, "int", "", caCert, caKey)

	privateKey, err := x509.MarshalPKCS8PrivateKey(intKey)
	require.NoError(t, err)

	pemCertificate := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: intCert.Raw,
	})

	pemPrivateKey := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privateKey,
	})

	certificateProvider := &services.LocalChargeStationCertificateProvider{
		Store:             store,
		CertificateReader: services.StringSource{Data: string(pemCertificate)},
		PrivateKeyReader:  services.StringSource{Data: string(pemPrivateKey)},
	}

	pemCsr := createCertificateSigningRequest(t)

	ctx := context.TODO()
	chain, err := certificateProvider.ProvideCertificate(ctx, services.CertificateTypeCSO, string(pemCsr), "cs001")
	require.NoError(t, err)

	require.NotEmpty(t, chain)

	block, r := pem.Decode([]byte(chain))
	count := 0
	issuerDN := ""
	for block != nil {
		if block.Type == "CERTIFICATE" {
			x509Cert, err := x509.ParseCertificate(block.Bytes)
			require.NoError(t, err)
			if count != 0 {
				require.Equal(t, issuerDN, x509Cert.Subject.String())
			}
			issuerDN = x509Cert.Issuer.String()
			count++
		}
		block, r = pem.Decode(r)
	}
	require.Equal(t, 2, count)
}

//	func TestLocalChargeStationCertificateProviderWithWrongId(t *testing.T) {
//		store := inmemory.NewStore(clock.RealClock{})
//
//		caCert, caKey := createRootCACertificate(t, "test")
//		intCert, intKey := createIntermediateCACertificate(t, "int", "", caCert, caKey)
//
//		privateKey, err := x509.MarshalPKCS8PrivateKey(intKey)
//		require.NoError(t, err)
//
//		pemCertificate := pem.EncodeToMemory(&pem.Block{
//			Type:  "CERTIFICATE",
//			Bytes: intCert.Raw,
//		})
//
//		pemPrivateKey := pem.EncodeToMemory(&pem.Block{
//			Type:  "PRIVATE KEY",
//			Bytes: privateKey,
//		})
//
//		certificateProvider := &services.LocalChargeStationCertificateProvider{
//			Store:             store,
//			CertificateReader: services.StringSource{Data: string(pemCertificate)},
//			PrivateKeyReader:  services.StringSource{Data: string(pemPrivateKey)},
//		}
//
//		pemCsr := createCertificateSigningRequest(t)
//
//		ctx := context.TODO()
//		chain, err := certificateProvider.ProvideCertificate(ctx, services.CertificateTypeCSO, string(pemCsr), "not-cs001")
//		assert.Error(t, err)
//		assert.Equal(t, "", chain)
//	}
func TestLocalChargeStationCertificateProviderWithRSAKey(t *testing.T) {
	store := inmemory.NewStore(clock.RealClock{})

	caCert, caKey := createRootCACertificate(t, "test")
	intCert, intKey := createIntermediateRSACACertificate(t, "int", "", caCert, caKey)

	privateKey := x509.MarshalPKCS1PrivateKey(intKey)

	pemCertificate := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: intCert.Raw,
	})

	pemPrivateKey := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKey,
	})

	certificateProvider := &services.LocalChargeStationCertificateProvider{
		Store:             store,
		CertificateReader: services.StringSource{Data: string(pemCertificate)},
		PrivateKeyReader:  services.StringSource{Data: string(pemPrivateKey)},
	}

	pemCsr := createCertificateSigningRequest(t)

	ctx := context.TODO()
	chain, err := certificateProvider.ProvideCertificate(ctx, services.CertificateTypeCSO, string(pemCsr), "cs001")
	require.NoError(t, err)

	require.NotEmpty(t, chain)

	block, r := pem.Decode([]byte(chain))
	count := 0
	issuerDN := ""
	for block != nil {
		if block.Type == "CERTIFICATE" {
			x509Cert, err := x509.ParseCertificate(block.Bytes)
			require.NoError(t, err)
			if count != 0 {
				require.Equal(t, issuerDN, x509Cert.Subject.String())
			}
			issuerDN = x509Cert.Issuer.String()
			count++
		}
		block, r = pem.Decode(r)
	}
	require.Equal(t, 2, count)
}

func TestLocalChargeStationCertificateProviderWithECKey(t *testing.T) {
	store := inmemory.NewStore(clock.RealClock{})

	caCert, caKey := createRootCACertificate(t, "test")
	intCert, intKey := createIntermediateCACertificate(t, "int", "", caCert, caKey)

	privateKey, err := x509.MarshalECPrivateKey(intKey)
	require.NoError(t, err)

	pemCertificate := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: intCert.Raw,
	})

	pemPrivateKey := pem.EncodeToMemory(&pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: privateKey,
	})

	certificateProvider := &services.LocalChargeStationCertificateProvider{
		Store:             store,
		CertificateReader: services.StringSource{Data: string(pemCertificate)},
		PrivateKeyReader:  services.StringSource{Data: string(pemPrivateKey)},
	}

	pemCsr := createCertificateSigningRequest(t)

	ctx := context.TODO()
	chain, err := certificateProvider.ProvideCertificate(ctx, services.CertificateTypeCSO, string(pemCsr), "cs001")
	require.NoError(t, err)

	require.NotEmpty(t, chain)

	block, r := pem.Decode([]byte(chain))
	count := 0
	issuerDN := ""
	for block != nil {
		if block.Type == "CERTIFICATE" {
			x509Cert, err := x509.ParseCertificate(block.Bytes)
			require.NoError(t, err)
			if count != 0 {
				require.Equal(t, issuerDN, x509Cert.Subject.String())
			}
			issuerDN = x509Cert.Issuer.String()
			count++
		}
		block, r = pem.Decode(r)
	}
	require.Equal(t, 2, count)
}

func TestLocalChargeStationCertificateIsAddedToStore(t *testing.T) {
	store := inmemory.NewStore(clock.RealClock{})

	caCert, caKey := createRootCACertificate(t, "test")
	intCert, intKey := createIntermediateCACertificate(t, "int", "", caCert, caKey)

	privateKey, err := x509.MarshalPKCS8PrivateKey(intKey)
	require.NoError(t, err)

	pemCertificate := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: intCert.Raw,
	})

	pemPrivateKey := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privateKey,
	})

	certificateProvider := &services.LocalChargeStationCertificateProvider{
		Store:             store,
		CertificateReader: services.StringSource{Data: string(pemCertificate)},
		PrivateKeyReader:  services.StringSource{Data: string(pemPrivateKey)},
	}

	pemCsr := createCertificateSigningRequest(t)

	ctx := context.TODO()
	chain, err := certificateProvider.ProvideCertificate(ctx, services.CertificateTypeCSO, string(pemCsr), "cs001")
	require.NoError(t, err)

	var x509Certificate *x509.Certificate
	block, r := pem.Decode([]byte(chain))
	for x509Certificate == nil && block != nil {
		if block.Type == "CERTIFICATE" {
			x509Certificate, err = x509.ParseCertificate(block.Bytes)
			require.NoError(t, err)
		}
		block, r = pem.Decode(r)
	}
	require.NotNil(t, x509Certificate)

	cert, err := store.LookupCertificate(ctx, getCertificateHash(x509Certificate))
	require.NoError(t, err)
	require.NotEqual(t, "", cert)
}

func getCertificateHash(cert *x509.Certificate) string {
	hash := sha256.Sum256(cert.Raw)
	b64Hash := base64.RawURLEncoding.EncodeToString(hash[:])
	return b64Hash
}

var (
	oidExtensionKeyUsage         = asn1.ObjectIdentifier{2, 5, 29, 15}
	oidExtensionExtendedKeyUsage = asn1.ObjectIdentifier{2, 5, 29, 37}
	oidExtensionBasicConstraints = asn1.ObjectIdentifier{2, 5, 29, 19}

	oidExtKeyUsageServerAuth = asn1.ObjectIdentifier{1, 3, 6, 1, 5, 5, 7, 3, 1}
	oidExtKeyUsageClientAuth = asn1.ObjectIdentifier{1, 3, 6, 1, 5, 5, 7, 3, 2}
)

type basicConstraints struct {
	IsCA       bool `asn1:"optional"`
	MaxPathLen int  `asn1:"optional,default:-1"`
}

func reverseBitsInAByte(in byte) byte {
	b1 := in>>4 | in<<4
	b2 := b1>>2&0x33 | b1<<2&0xcc
	b3 := b2>>1&0x55 | b2<<1&0xaa
	return b3
}

// asn1BitLength returns the bit-length of bitString by considering the
// most-significant bit in a byte to be the "first" bit. This convention
// matches ASN.1, but differs from almost everything else.
func asn1BitLength(bitString []byte) int {
	bitLen := len(bitString) * 8

	for i := range bitString {
		b := bitString[len(bitString)-i-1]

		for bit := uint(0); bit < 8; bit++ {
			if (b>>bit)&1 == 1 {
				return bitLen
			}
			bitLen--
		}
	}

	return 0
}

func marshalKeyUsage() (pkix.Extension, error) {
	ext := pkix.Extension{Id: oidExtensionKeyUsage, Critical: true}
	ku := x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment

	var a [2]byte
	a[0] = reverseBitsInAByte(byte(ku))
	a[1] = reverseBitsInAByte(byte(ku >> 8))

	l := 1
	if a[1] != 0 {
		l = 2
	}

	bitString := a[:l]
	var err error
	ext.Value, err = asn1.Marshal(asn1.BitString{Bytes: bitString, BitLength: asn1BitLength(bitString)})
	return ext, err
}

func marshalExtKeyUsage(oids []asn1.ObjectIdentifier) (pkix.Extension, error) {
	ext := pkix.Extension{Id: oidExtensionExtendedKeyUsage}

	var err error
	ext.Value, err = asn1.Marshal(oids)
	return ext, err
}

func marshalBasicConstraints() (pkix.Extension, error) {
	ext := pkix.Extension{Id: oidExtensionBasicConstraints, Critical: true}
	var err error
	ext.Value, err = asn1.Marshal(basicConstraints{false, -1})
	return ext, err
}

func createCertificateSigningRequest(t *testing.T) []byte {
	keyPair, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)

	var extensions []pkix.Extension
	bcExt, err := marshalBasicConstraints()
	require.NoError(t, err)
	extensions = append(extensions, bcExt)
	kuExt, err := marshalKeyUsage()
	require.NoError(t, err)
	extensions = append(extensions, kuExt)
	ekuExt, err := marshalExtKeyUsage([]asn1.ObjectIdentifier{oidExtKeyUsageClientAuth, oidExtKeyUsageServerAuth})
	require.NoError(t, err)
	extensions = append(extensions, ekuExt)

	csrTemplate := x509.CertificateRequest{
		Subject: pkix.Name{
			CommonName:   "cs001",
			Organization: []string{"Thoughtworks"},
			Country:      []string{"GB"},
		},
		SignatureAlgorithm: x509.ECDSAWithSHA256,
		ExtraExtensions:    extensions,
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

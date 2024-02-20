// SPDX-License-Identifier: Apache-2.0

package services

import (
	"bytes"
	"context"
	"crypto"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"go.mozilla.org/pkcs7"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.19.0"
	"go.opentelemetry.io/otel/trace"
	"io"
	"math/big"
	"net/http"
	"time"
)

type CertificateType int

const (
	CertificateTypeV2G CertificateType = iota
	CertificateTypeCSO
)

type ChargeStationCertificateProvider interface {
	ProvideCertificate(ctx context.Context, typ CertificateType, pemEncodedCSR string, csId string) (pemEncodedCertificateChain string, err error)
}

type ISOVersion string

const (
	ISO15118V2  ISOVersion = "ISO15118-2"
	ISO15118V20 ISOVersion = "ISO15118-20"
)

// The OpcpChargeStationCertificateProvider issues certificates using the CPO CA
type OpcpChargeStationCertificateProvider struct {
	BaseURL          string
	HttpTokenService HttpTokenService
	ISOVersion       ISOVersion
	HttpClient       *http.Client
}

type HttpError int

func (h HttpError) Error() string {
	return fmt.Sprintf("http status: %d", h)
}

func (h OpcpChargeStationCertificateProvider) ProvideCertificate(ctx context.Context, typ CertificateType, pemEncodedCSR string, csId string) (string, error) {
	if typ == CertificateTypeV2G {
		cn, err := getCommonNameFromCSR(pemEncodedCSR)
		if err != nil {
			return "", err
		}

		if cn != csId {
			return "", fmt.Errorf("certificate common name does not match charge station id")
		}

		csr, err := convertCSR(pemEncodedCSR)
		if err != nil {
			return "", err
		}

		cert, err := h.requestCertificateWithRetry(ctx, csr)
		if err != nil {
			return "", fmt.Errorf("requesting certificate: %w", err)
		}

		chain, err := h.requestChainWithRetry(ctx)
		if err != nil {
			return "", fmt.Errorf("requesting ca certificates: %w", err)
		}

		return encodeLeafAndChain(cert, chain)
	} else {
		return "", fmt.Errorf("certificate type %d not supported", typ)
	}
}

func (h OpcpChargeStationCertificateProvider) requestCertificateWithRetry(ctx context.Context, csr []byte) ([]byte, error) {
	span := trace.SpanFromContext(ctx)
	newCtx, span := span.TracerProvider().Tracer("manager").Start(ctx, "sign_certificate")
	defer span.End()

	certBytes, err := h.requestCertificate(newCtx, csr, false)
	if err != nil {
		span.SetAttributes(semconv.HTTPResendCount(1))
		certBytes, err = h.requestCertificate(newCtx, csr, true)
		if err != nil {
			span.SetStatus(codes.Error, "retries exhausted")
			span.RecordError(err)
		}
	}
	return certBytes, err
}

func (h OpcpChargeStationCertificateProvider) requestCertificate(ctx context.Context, csr []byte, isRetry bool) ([]byte, error) {
	enrollUrl := fmt.Sprintf("%s/cpo/simpleenroll/%s", h.BaseURL, h.ISOVersion)
	req, err := http.NewRequestWithContext(ctx, "POST", enrollUrl, bytes.NewReader(csr))
	if err != nil {
		return nil, err
	}

	token, err := h.HttpTokenService.GetToken(ctx, isRetry)
	if err != nil {
		return nil, fmt.Errorf("getting token: %w", err)
	}

	req.Header.Add("x-api-key", token)
	req.Header.Add("authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Add("content-type", "application/pkcs10")

	resp, err := h.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode == http.StatusOK {
		return io.ReadAll(resp.Body)
	}

	return nil, HttpError(resp.StatusCode)
}

func (h OpcpChargeStationCertificateProvider) requestChainWithRetry(ctx context.Context) ([]byte, error) {
	span := trace.SpanFromContext(ctx)
	newCtx, span := span.TracerProvider().Tracer("manager").Start(ctx, "get_ca_certificates")
	defer span.End()

	chainBytes, err := h.requestChain(newCtx, false)
	if err != nil {
		chainBytes, err = h.requestChain(newCtx, true)
	}
	return chainBytes, err
}

func (h OpcpChargeStationCertificateProvider) requestChain(ctx context.Context, isRetry bool) ([]byte, error) {
	caUrl := fmt.Sprintf("%s/cpo/cacerts/%s", h.BaseURL, h.ISOVersion)
	req, err := http.NewRequestWithContext(ctx, "GET", caUrl, nil)
	if err != nil {
		return nil, err
	}

	token, err := h.HttpTokenService.GetToken(ctx, isRetry)
	if err != nil {
		return nil, fmt.Errorf("getting token: %w", err)
	}

	req.Header.Add("x-api-key", token)
	req.Header.Add("authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := h.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode == http.StatusOK {
		return io.ReadAll(resp.Body)
	}

	return nil, HttpError(resp.StatusCode)
}

func convertCSR(pemEncodedCSR string) ([]byte, error) {
	block, _ := pem.Decode([]byte(pemEncodedCSR))
	if block == nil {
		return nil, fmt.Errorf("no PEM data in input")
	}
	if block.Type != "CERTIFICATE REQUEST" {
		return nil, fmt.Errorf("no certificate request in PEM block")
	}

	enc := make([]byte, base64.StdEncoding.EncodedLen(len(block.Bytes)))
	base64.StdEncoding.Encode(enc, block.Bytes)
	return enc, nil
}

func encodeLeafAndChain(leaf, chain []byte) (string, error) {
	decodedLeaf := make([]byte, base64.StdEncoding.DecodedLen(len(leaf)))
	n, err := base64.StdEncoding.Decode(decodedLeaf, leaf)
	if err != nil {
		return "", err
	}

	p7Leaf, err := pkcs7.Parse(decodedLeaf[:n])
	if err != nil {
		return "", err
	}

	decodedChain := make([]byte, base64.StdEncoding.DecodedLen(len(chain)))
	n, err = base64.StdEncoding.Decode(decodedChain, chain)
	if err != nil {
		return "", err
	}

	p7Chain, err := pkcs7.Parse(decodedChain[:n])
	if err != nil {
		return "", err
	}

	var pemBytes []byte

	subject := p7Leaf.Certificates[0].Subject.String()
	issuer := p7Leaf.Certificates[0].Issuer.String()
	pemBytes = append(pemBytes, fmt.Sprintf("subject: %s\n", subject)...)
	pemBytes = append(pemBytes, fmt.Sprintf("issuer: %s\n", issuer)...)
	pemBytes = append(pemBytes, pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: p7Leaf.Certificates[0].Raw,
	})...)

	for {
		var match bool
		for _, cert := range p7Chain.Certificates {
			if cert.Subject.String() == issuer && cert.Subject.String() != cert.Issuer.String() {
				subject = cert.Subject.String()
				issuer = cert.Issuer.String()
				pemBytes = append(pemBytes, fmt.Sprintf("subject: %s\n", subject)...)
				pemBytes = append(pemBytes, fmt.Sprintf("issuer: %s\n", issuer)...)
				pemBytes = append(pemBytes, pem.EncodeToMemory(&pem.Block{
					Type:  "CERTIFICATE",
					Bytes: cert.Raw,
				})...)
				match = true
				break
			}
		}
		if !match {
			break
		}
	}

	return string(pemBytes), nil
}

type DefaultChargeStationCertificateProvider struct{}

func (n DefaultChargeStationCertificateProvider) ProvideCertificate(context.Context, CertificateType, string, string) (pemEncodedCertificateChain string, err error) {
	return "", fmt.Errorf("not implemented")
}

// LocalChargeStationCertificateProvider issues CSO certificates using local data
type LocalChargeStationCertificateProvider struct {
	CertificateReader io.Reader
	PrivateKeyReader  io.Reader
}

func (l *LocalChargeStationCertificateProvider) ProvideCertificate(ctx context.Context, typ CertificateType, pemEncodedCSR string, csId string) (pemEncodedCertificateChain string, err error) {
	if typ != CertificateTypeCSO {
		return "", fmt.Errorf("local provider cannot provide V2G certificate")
	}

	cn, err := getCommonNameFromCSR(pemEncodedCSR)
	if err != nil {
		return "", err
	}

	if cn != csId {
		return "", fmt.Errorf("certificate common name does not match charge station id")
	}

	certificate, err := readSigningCertificate(l.CertificateReader)
	if err != nil {
		return "", err
	}

	privateKey, err := readPrivateKey(l.PrivateKeyReader)
	if err != nil {
		return "", err
	}

	csr, err := readCsr(pemEncodedCSR)
	if err != nil {
		return "", err
	}

	serial, err := rand.Int(rand.Reader, (&big.Int{}).Exp(big.NewInt(2), big.NewInt(159), nil))
	if err != nil {
		return "", fmt.Errorf("creating serial number: %w", err)
	}

	now := time.Now()

	var extensions []pkix.Extension
	ekuExt, err := marshalExtKeyUsage([]asn1.ObjectIdentifier{oidExtKeyUsageClientAuth, oidExtKeyUsageServerAuth})
	if err != nil {
		return "", fmt.Errorf("marshalling ext key usage: %v", err)
	}
	extensions = append(extensions, ekuExt)

	template := x509.Certificate{
		SerialNumber:          serial,
		Subject:               csr.Subject,
		NotBefore:             now,
		NotAfter:              now.Add(time.Minute),
		BasicConstraintsValid: true,
		IsCA:                  false,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtraExtensions:       extensions,
	}

	signedCert, err := x509.CreateCertificate(rand.Reader, &template, certificate, csr.PublicKey, privateKey)
	if err != nil {
		return "", fmt.Errorf("creating certificate: %v", err)
	}

	pemSigningCert := string(pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certificate.Raw,
	}))

	pemSignedCert := string(pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: signedCert,
	}))

	return pemSignedCert + pemSigningCert, nil
}

func readSigningCertificate(certificateReader io.Reader) (*x509.Certificate, error) {
	pemCertificate, err := io.ReadAll(certificateReader)
	if err != nil {
		return nil, fmt.Errorf("reading certificate: %v", err)
	}

	var certificate *x509.Certificate
	block, r := pem.Decode(pemCertificate)
	for certificate == nil && block != nil {
		if block.Type == "CERTIFICATE" {
			certificate, err = x509.ParseCertificate(block.Bytes)
			if err != nil {
				return nil, fmt.Errorf("cannot parse certificate: %v", err)
			}
		}
		block, r = pem.Decode(r)
	}

	if certificate == nil {
		return nil, fmt.Errorf("no signing certificate")
	}

	return certificate, nil
}

func readPrivateKey(privateKeyReader io.Reader) (crypto.PrivateKey, error) {
	pemPrivateKey, err := io.ReadAll(privateKeyReader)
	if err != nil {
		return "", fmt.Errorf("reading private key: %v", err)
	}

	var privateKey crypto.PrivateKey
	block, r := pem.Decode(pemPrivateKey)
	for privateKey == nil && block != nil {
		if block.Type == "PRIVATE KEY" {
			privateKey, err = x509.ParsePKCS8PrivateKey(block.Bytes)
			if err != nil {
				return "", fmt.Errorf("cannot parse private key: %v", err)
			}
		}
		if block.Type == "RSA PRIVATE KEY" {
			privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
			if err != nil {
				return "", fmt.Errorf("cannot parse rsa private key: %v", err)
			}
		}
		if block.Type == "EC PRIVATE KEY" {
			privateKey, err = x509.ParseECPrivateKey(block.Bytes)
			if err != nil {
				return "", fmt.Errorf("cannot parse ec private key: %v", err)
			}
		}
		block, r = pem.Decode(r)
	}

	if privateKey == nil {
		return "", fmt.Errorf("no private key")
	}

	return privateKey, nil
}

func readCsr(pemEncodedCSR string) (*x509.CertificateRequest, error) {
	var csr *x509.CertificateRequest
	block, r := pem.Decode([]byte(pemEncodedCSR))
	for csr == nil && block != nil {
		var err error
		if block.Type == "CERTIFICATE REQUEST" {
			csr, err = x509.ParseCertificateRequest(block.Bytes)
			if err != nil {
				return nil, fmt.Errorf("cannot parse csr: %v", err)
			}
		}
		block, r = pem.Decode(r)
	}

	if csr == nil {
		return nil, fmt.Errorf("no csr")
	}

	return csr, nil
}

var (
	oidExtensionExtendedKeyUsage = asn1.ObjectIdentifier{2, 5, 29, 37}

	oidExtKeyUsageServerAuth = asn1.ObjectIdentifier{1, 3, 6, 1, 5, 5, 7, 3, 1}
	oidExtKeyUsageClientAuth = asn1.ObjectIdentifier{1, 3, 6, 1, 5, 5, 7, 3, 2}
)

func marshalExtKeyUsage(oids []asn1.ObjectIdentifier) (pkix.Extension, error) {
	ext := pkix.Extension{Id: oidExtensionExtendedKeyUsage}

	var err error
	ext.Value, err = asn1.Marshal(oids)
	return ext, err
}

type DelegatingChargeStationCertificateProvider struct {
	V2GChargeStationCertificateProvider ChargeStationCertificateProvider
	CSOChargeStationCertificateProvider ChargeStationCertificateProvider
}

func (d *DelegatingChargeStationCertificateProvider) ProvideCertificate(ctx context.Context, typ CertificateType, pemEncodedCSR string, csId string) (pemEncodedCertificateChain string, err error) {
	if typ == CertificateTypeV2G {
		return d.V2GChargeStationCertificateProvider.ProvideCertificate(ctx, typ, pemEncodedCSR, csId)
	} else if typ == CertificateTypeCSO {
		return d.CSOChargeStationCertificateProvider.ProvideCertificate(ctx, typ, pemEncodedCSR, csId)
	} else {
		return "", fmt.Errorf("unknown certificate type: %d", typ)
	}
}

func getCommonNameFromCSR(csr string) (string, error) {
	block, _ := pem.Decode([]byte(csr))
	if block != nil && block.Type == "CERTIFICATE REQUEST" {
		csr, err := x509.ParseCertificateRequest(block.Bytes)
		if err != nil {
			return "", fmt.Errorf("parsing certificate request: %w", err)
		}
		return csr.Subject.CommonName, nil
	}
	return "", fmt.Errorf("no certificate request in PEM block")
}

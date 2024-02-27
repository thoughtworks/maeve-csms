// SPDX-License-Identifier: Apache-2.0

package services

import (
	"context"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
	"io"
	"k8s.io/utils/clock"
	"net/http"
	"os"
	"sync"
	"time"
)

type RootCertificateProviderService interface {
	ProvideCertificates(ctx context.Context) ([]*x509.Certificate, error)
}

type CompositeRootCertificateProviderService struct {
	Providers []RootCertificateProviderService
}

func (c CompositeRootCertificateProviderService) ProvideCertificates(ctx context.Context) ([]*x509.Certificate, error) {
	var certs []*x509.Certificate

	for _, provider := range c.Providers {
		providerCerts, err := provider.ProvideCertificates(ctx)
		if err != nil {
			return nil, err
		}
		certs = append(certs, providerCerts...)
	}

	return certs, nil
}

type FileRootCertificateProviderService struct {
	FilePaths []string
}

func (s FileRootCertificateProviderService) ProvideCertificates(context.Context) (certs []*x509.Certificate, err error) {
	for _, pemFile := range s.FilePaths {
		//#nosec G304 - only files specified by the person running the application will be loaded
		bytes, e := os.ReadFile(pemFile)
		if e != nil {
			return nil, e
		}

		intermediateCerts, e := parseCertificates(bytes)
		if e != nil {
			return nil, e
		}

		certs = append(certs, intermediateCerts...)
	}
	return
}

func parseCertificates(pemData []byte) ([]*x509.Certificate, error) {
	var certs []*x509.Certificate
	for {
		cert, rest, err := parseCertificate(pemData)
		if err != nil {
			return nil, err
		}
		if cert == nil {
			break
		}
		certs = append(certs, cert)
		pemData = rest
	}
	return certs, nil
}

func parseCertificate(pemData []byte) (cert *x509.Certificate, rest []byte, err error) {
	block, rest := pem.Decode(pemData)
	if block == nil {
		return
	}
	if block.Type != "CERTIFICATE" {
		return
	}
	cert, err = x509.ParseCertificate(block.Bytes)
	if err != nil {
		cert = nil
		return
	}
	return
}

type CachingRootCertificateProviderService struct {
	sync.Mutex
	delegate RootCertificateProviderService
	ttl      time.Duration
	clock    clock.PassiveClock
	expiry   time.Time
	certs    []*x509.Certificate
}

func NewCachingRootCertificateProviderService(delegate RootCertificateProviderService, ttl time.Duration, clock clock.PassiveClock) *CachingRootCertificateProviderService {
	return &CachingRootCertificateProviderService{
		delegate: delegate,
		ttl:      ttl,
		clock:    clock,
		expiry:   clock.Now(),
	}
}

func (c *CachingRootCertificateProviderService) ProvideCertificates(ctx context.Context) ([]*x509.Certificate, error) {
	c.Lock()
	defer c.Unlock()

	if c.clock.Now().After(c.expiry) {
		var err error
		c.certs, err = c.delegate.ProvideCertificates(ctx)
		if err != nil {
			return nil, err
		}
		c.expiry = c.clock.Now().Add(c.ttl)
	}

	return c.certs, nil
}

type OpcpRootCertificateReturnType struct {
	RootCertificateCollection OpcpRootCertificateCollectionType `json:"RootCertificateCollection"`
}

type OpcpRootCertificateCollectionType struct {
	RootCertificates []OpcpRootCertificateType `json:"rootCertificates"`
}

type OpcpRootCertificateType struct {
	RootCertificateId          string `json:"rootCertificateId"`
	DN                         string `json:"distinguishedName"`
	CACertificate              string `json:"caCertificate"`
	CommonName                 string `json:"commonName"`
	RootAuthorityKeyIdentifier string `json:"rootAuthorityKeyIdentifier"`
	RootIssuerSerialNumber     string `json:"rootIssuerSerialNumber"`
	ValidFrom                  string `json:"validFrom"`
	ValidTo                    string `json:"validTo"`
	OrganizationName           string `json:"organizationName"`
	CertificateRevocationList  string `json:"certificateRevocationList"`
	CrossCertificatePair       string `json:"crossCertificatePair"`
	LabeledURI                 string `json:"labeledUri"`
	RootType                   string `json:"rootType"`
	Fingerprint                string `json:"fingerprint"`
}

type OpcpRootCertificateProviderService struct {
	BaseURL      string
	TokenService HttpTokenService
	HttpClient   *http.Client
}

func (s OpcpRootCertificateProviderService) ProvideCertificates(ctx context.Context) (certs []*x509.Certificate, err error) {
	body, err := s.retrieveCertificatesFromUrlWithRetry(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve certificates from url: %w", err)
	}

	var rootCertificates OpcpRootCertificateReturnType
	err = json.Unmarshal(body, &rootCertificates)
	if err != nil {
		return nil, err
	}

	for _, certDetails := range rootCertificates.RootCertificateCollection.RootCertificates {
		certBytes, err := base64.StdEncoding.DecodeString(certDetails.CACertificate)
		if err != nil {
			return nil, fmt.Errorf("failed to decode certificate %s: %w", certDetails.RootCertificateId, err)
		}
		cert, err := x509.ParseCertificate(certBytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse certificate %s: %w", certDetails.RootCertificateId, err)
		}
		certs = append(certs, cert)
	}

	return
}

func (s OpcpRootCertificateProviderService) retrieveCertificatesFromUrlWithRetry(ctx context.Context) ([]byte, error) {
	span := trace.SpanFromContext(ctx)
	newCtx, span := span.TracerProvider().Tracer("manager").Start(ctx, "get_certificates_from_url")
	defer span.End()

	certs, err := s.retrieveCertificatesFromUrl(newCtx, false)
	if err != nil {
		span.SetAttributes(semconv.HTTPResendCount(1))
		certs, err = s.retrieveCertificatesFromUrl(newCtx, true)
		if err != nil {
			span.SetStatus(codes.Error, "retries exhausted")
			span.RecordError(err)
		}
	}
	return certs, err
}

func (s OpcpRootCertificateProviderService) retrieveCertificatesFromUrl(ctx context.Context, retry bool) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/v1/root/rootCerts", s.BaseURL), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	token, err := s.TokenService.GetToken(ctx, retry)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("X-API-Key", token)

	resp, err := s.HttpClient.Do(req)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		return nil, HttpError(resp.StatusCode)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

type X509RootCertificateProviderService struct {
	Certificates []*x509.Certificate
}

func (x X509RootCertificateProviderService) ProvideCertificates(ctx context.Context) ([]*x509.Certificate, error) {
	return x.Certificates, nil
}

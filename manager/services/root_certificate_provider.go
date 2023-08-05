package services

import (
	"context"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"os"
)

type RootCertificateProviderService interface {
	ProvideCertificates(ctx context.Context) ([]*x509.Certificate, error)
}

type FileRootCertificateRetrieverService struct {
	FilePaths []string
}

func (s FileRootCertificateRetrieverService) ProvideCertificates(context.Context) (certs []*x509.Certificate, err error) {
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

type OpcpRootCertificateRetrieverService struct {
	BaseURL         string
	HttpAuthService HttpAuthService
	HttpClient      *http.Client
}

func (s OpcpRootCertificateRetrieverService) ProvideCertificates(ctx context.Context) (certs []*x509.Certificate, err error) {
	body, err := s.retrieveCertificatesFromUrl(ctx)
	if err != nil {
		return nil, err
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

func (s OpcpRootCertificateRetrieverService) retrieveCertificatesFromUrl(ctx context.Context) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/v1/root/rootCerts", s.BaseURL), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	err = s.HttpAuthService.Authenticate(req)
	if err != nil {
		return nil, err
	}

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

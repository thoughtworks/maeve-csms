package services

import (
	"bytes"
	"crypto"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/twlabs/maeve-csms/manager/ocpp/ocpp201"
	"golang.org/x/crypto/ocsp"
	"io"
	"log"
	"math/big"
	"net/http"
)

// OCSPError is an error returned by the OCSP server in response to a check
type OCSPError int

func (o OCSPError) Error() string {
	return fmt.Sprintf("ocsp validation failed: %d", o)
}

type ValidationError int

const (
	ValidationErrorCertChain ValidationError = iota
	ValidationErrorCertExpired
	ValidationErrorCertRevoked
	ValidationErrorWrongEmaid
)

func (v ValidationError) Error() string {
	switch v {
	case ValidationErrorCertChain:
		return "certificate chain invalid"
	case ValidationErrorCertExpired:
		return "certificate expired"
	case ValidationErrorCertRevoked:
		return "certificate revoked"
	case ValidationErrorWrongEmaid:
		return "wrong emaid"
	default:
		return fmt.Sprintf("unknown error: %d", v)
	}
}

type CertificateValidationService interface {
	ValidatePEMCertificateChain(pemChain []byte, eMAID string) (*string, error)
	ValidateHashedCertificateChain(ocspRequestData []ocpp201.OCSPRequestDataType) (*string, error)
}

type OnlineCertificateValidationService struct {
	RootCertificates []*x509.Certificate
	MaxOCSPAttempts  int
}

func (o OnlineCertificateValidationService) ValidatePEMCertificateChain(pemChain []byte, eMAID string) (*string, error) {
	certificateChain, err := ParseCertificates(pemChain)
	if err != nil {
		return nil, err
	}

	if len(certificateChain) == 0 {
		return nil, fmt.Errorf("empty certificate chain")
	}

	err = validateEMAID(certificateChain[0], eMAID)
	if err != nil {
		return nil, err
	}

	err = validatePEMCertificateChain(certificateChain, o.RootCertificates)
	if err != nil {
		return nil, err
	}

	ocspResponse, err := validatePEMCertificateChainOCSPStatus(certificateChain, o.RootCertificates, o.MaxOCSPAttempts)
	if err != nil {
		return ocspResponse, err
	}

	return ocspResponse, nil
}

func (o OnlineCertificateValidationService) ValidateHashedCertificateChain(ocspRequestData []ocpp201.OCSPRequestDataType) (*string, error) {
	var ocspResponse *string
	var err error
	for _, requestData := range ocspRequestData {
		var ocspRequest []byte
		ocspRequest, err = createOCSPRequestFromHashData(requestData.IssuerNameHash, requestData.IssuerKeyHash,
			requestData.SerialNumber, string(requestData.HashAlgorithm))
		if err != nil {
			return nil, err
		}
		var ocspResp *string
		ocspResp, err = performOCSPCheck([]string{requestData.ResponderURL}, ocspRequest, nil, o.MaxOCSPAttempts)
		if err != nil {
			return ocspResp, err
		}
		if ocspResponse == nil {
			ocspResponse = ocspResp
		}
	}

	return ocspResponse, err
}

func validatePEMCertificateChain(certificateChain, rootCertificates []*x509.Certificate) error {
	if len(certificateChain) < 1 {
		return errors.New("no certificates in chain")
	}

	intermediates := x509.NewCertPool()
	if len(certificateChain) > 1 {
		for _, cert := range certificateChain[1:] {
			intermediates.AddCert(cert)
		}
	}

	trustedCerts := x509.NewCertPool()
	for _, cert := range rootCertificates {
		trustedCerts.AddCert(cert)
	}

	opts := x509.VerifyOptions{
		Roots:         trustedCerts,
		Intermediates: intermediates,
		KeyUsages:     []x509.ExtKeyUsage{x509.ExtKeyUsageAny},
	}

	_, err := certificateChain[0].Verify(opts)
	var invalidCertErr *x509.CertificateInvalidError
	if errors.As(err, &invalidCertErr) && invalidCertErr.Reason == x509.Expired {
		return fmt.Errorf("validating certificate chain: %s: %w", invalidCertErr, ValidationErrorCertExpired)
	} else if err != nil {
		return fmt.Errorf("validating certificate chain: %s: %w", err, ValidationErrorCertChain)
	}

	return nil
}

func validateEMAID(certificate *x509.Certificate, eMAID string) error {
	if certificate.Subject.CommonName != eMAID {
		return fmt.Errorf("leaf certificate CN: %s, not %s: %w", certificate.Subject.CommonName, eMAID, ValidationErrorWrongEmaid)
	}

	return nil
}

func validatePEMCertificateChainOCSPStatus(certificateChain, rootCertificates []*x509.Certificate, maxRetries int) (*string, error) {
	if len(certificateChain) < 1 {
		return nil, fmt.Errorf("no certificates in chain: %w", ValidationErrorCertChain)
	}
	for i := 0; i < len(certificateChain)-1; i++ {
		ocspRequest, err := createOCSPRequestFromCertificate(certificateChain[i], certificateChain[i+1])
		if err != nil {
			return nil, err
		}
		ocspResponse, err := performOCSPCheck(certificateChain[i].OCSPServer, ocspRequest, certificateChain[i+1], maxRetries)
		if err != nil {
			return ocspResponse, err
		}
	}
	subjectCert := certificateChain[len(certificateChain)-1]
	for _, rootCert := range rootCertificates {
		if bytes.Equal(subjectCert.AuthorityKeyId, rootCert.SubjectKeyId) {
			ocspRequest, err := createOCSPRequestFromCertificate(subjectCert, rootCert)
			if err != nil {
				return nil, err
			}
			return performOCSPCheck(subjectCert.OCSPServer, ocspRequest, rootCert, maxRetries)
		}
	}

	return nil, fmt.Errorf("missing root cert for chain: %w", ValidationErrorCertChain)
}

func performOCSPCheck(ocspResponderUrls []string, ocspRequest []byte, issuerCert *x509.Certificate, maxAttempts int) (*string, error) {
	var ocspResponderUrlCount int
	if ocspResponderUrls != nil {
		ocspResponderUrlCount = len(ocspResponderUrls)
	} else {
		return nil, errors.New("no OCSP servers specified")
	}

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		responderUrl := ocspResponderUrls[attempt%ocspResponderUrlCount]
		ocspResponse, err := attemptOCSPCheck(responderUrl, ocspRequest, issuerCert)
		if err == nil {
			return ocspResponse, nil
		}
		var ocspError *OCSPError
		if errors.As(err, &ocspError) {
			return ocspResponse, fmt.Errorf("ocsp check status: %d: %w", ocspError, ValidationErrorCertRevoked)
		}
		log.Printf("ocsp check attempt %d/%d failed: %v", attempt, maxAttempts, err)
	}

	return nil, fmt.Errorf("failed to perform ocsp check after %d attempts", maxAttempts)
}

func createOCSPRequestFromCertificate(subjectCert, issuerCert *x509.Certificate) ([]byte, error) {
	return ocsp.CreateRequest(subjectCert, issuerCert, nil)
}

func createOCSPRequestFromHashData(issuerNameHashHex, issuerKeyHashHex, serialNumber, hashAlg string) ([]byte, error) {
	n := new(big.Int)
	serial, ok := n.SetString(serialNumber, 16)
	if !ok {
		return nil, fmt.Errorf("unable to parse serial number %s as base 16 string", serial)
	}

	issuerNameHash, err := hex.DecodeString(issuerNameHashHex)
	if err != nil {
		return nil, fmt.Errorf("parsing issuer name hash %s: %w", issuerNameHashHex, err)
	}

	issuerKeyHash, err := hex.DecodeString(issuerKeyHashHex)
	if err != nil {
		return nil, fmt.Errorf("parsing issuer key hash %s: %w", issuerKeyHash, err)
	}

	req := ocsp.Request{
		HashAlgorithm:  hashAlgorithm(hashAlg),
		IssuerNameHash: issuerNameHash,
		IssuerKeyHash:  issuerKeyHash,
		SerialNumber:   serial,
	}

	return req.Marshal()
}

func attemptOCSPCheck(ocspResponderUrl string, ocspRequest []byte, issuerCert *x509.Certificate) (*string, error) {
	//#nosec G107 - need to use OCSP URL specified in the certificate
	resp, err := http.Post(ocspResponderUrl, "application/ocsp-request", bytes.NewReader(ocspRequest))
	if err != nil {
		return nil, fmt.Errorf("post %s: %w", ocspResponderUrl, err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("post %s status %s", ocspResponderUrl, resp.Status)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading body: %w", err)
	}
	base64Response := base64.StdEncoding.EncodeToString(respBytes)
	ocspResp, err := ocsp.ParseResponse(respBytes, issuerCert)
	if err != nil {
		return &base64Response, fmt.Errorf("parsing ocsp response: %w", err)
	}
	if ocspResp.Status != ocsp.Good {
		return &base64Response, OCSPError(ocspResp.Status)
	}

	return &base64Response, nil
}

func ParseCertificates(pemData []byte) ([]*x509.Certificate, error) {
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

func hashAlgorithm(hash string) crypto.Hash {
	switch hash {
	case "SHA256":
		return crypto.SHA256
	case "SHA384":
		return crypto.SHA384
	case "SHA512":
		return crypto.SHA512
	default:
		return crypto.SHA256
	}
}

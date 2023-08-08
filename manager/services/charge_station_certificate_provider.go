// SPDX-License-Identifier: Apache-2.0

package services

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"go.mozilla.org/pkcs7"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
	"io"
	"net/http"
)

type CertificateType int

const (
	CertificateTypeV2G CertificateType = iota
	CertificateTypeCSO
)

type ChargeStationCertificateProvider interface {
	ProvideCertificate(ctx context.Context, typ CertificateType, pemEncodedCSR string) (pemEncodedCertificateChain string, err error)
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

func (h OpcpChargeStationCertificateProvider) ProvideCertificate(ctx context.Context, typ CertificateType, pemEncodedCSR string) (string, error) {
	if typ == CertificateTypeV2G {
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

func (n DefaultChargeStationCertificateProvider) ProvideCertificate(context.Context, CertificateType, string) (pemEncodedCertificateChain string, err error) {
	return "", fmt.Errorf("not implemented")
}

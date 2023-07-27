// SPDX-License-Identifier: Apache-2.0

package services

import (
	"bytes"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"go.mozilla.org/pkcs7"
	"io"
	"net/http"
)

type CertificateType int

const (
	CertificateTypeV2G CertificateType = iota
	CertificateTypeCSO
)

type CertificateSignerService interface {
	SignCertificate(typ CertificateType, pemEncodedCSR string) (pemEncodedCertificateChain string, err error)
}

type ISOVersion string

const (
	ISO15118V2  ISOVersion = "ISO15118-2"
	ISO15118V20 ISOVersion = "ISO15118-20"
)

// The OpcpCpoCertificateSignerService issues certificates using the CPO CA
type OpcpCpoCertificateSignerService struct {
	BaseURL     string
	BearerToken string
	ISOVersion  ISOVersion
	HttpClient  *http.Client
}

type HttpError int

func (h HttpError) Error() string {
	return fmt.Sprintf("http status: %d", h)
}

func (h OpcpCpoCertificateSignerService) SignCertificate(typ CertificateType, pemEncodedCSR string) (string, error) {
	csr, err := convertCSR(pemEncodedCSR)
	if err != nil {
		return "", err
	}

	client := h.HttpClient
	if client == nil {
		client = http.DefaultClient
	}

	enrollUrl := fmt.Sprintf("%s/cpo/simpleenroll/%s", h.BaseURL, h.ISOVersion)
	cert, err := requestCertificate(client, enrollUrl, h.BearerToken, csr)
	if err != nil {
		return "", fmt.Errorf("requesting certificate: %w", err)
	}

	caUrl := fmt.Sprintf("%s/cpo/cacerts/%s", h.BaseURL, h.ISOVersion)
	chain, err := requestChain(client, caUrl, h.BearerToken)
	if err != nil {
		return "", fmt.Errorf("requesting ca certificates: %w", err)
	}

	return encodeLeafAndChain(cert, chain)
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

func requestCertificate(client *http.Client, url, bearerToken string, csr []byte) ([]byte, error) {
	req, err := http.NewRequest("POST", url, bytes.NewReader(csr))
	if err != nil {
		return nil, err
	}
	req.Header.Add("x-api-key", bearerToken)
	req.Header.Add("authorization", fmt.Sprintf("Bearer %s", bearerToken))
	req.Header.Add("content-type", "application/pkcs10")

	resp, err := client.Do(req)
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

func requestChain(client *http.Client, url, bearerToken string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("authorization", fmt.Sprintf("Bearer %s", bearerToken))

	resp, err := client.Do(req)
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

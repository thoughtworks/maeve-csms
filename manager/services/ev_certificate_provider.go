// SPDX-License-Identifier: Apache-2.0

package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
)

const XsdMsgDefinition = "urn:iso:15118:2:2013:MsgDef"

type EvCertificateProvider interface {
	ProvideCertificate(exiRequest string) (EvCertificate15118Response, error)
}

type OpcpMoEvCertificateProvider struct {
	BaseURL     string
	BearerToken string
	HttpClient  *http.Client
}

type EvCertificate15118Response struct {
	Status                     ocpp201.Iso15118EVCertificateStatusEnumType
	CertificateInstallationRes string
}

type SignedContractDataRequest struct {
	CertificateInstallationReq string `json:"certificateInstallationReq"`
	XsdMsgDefNamespace         string `json:"xsdMsgDefNamespace"`
}

type SignedContractDataResponse struct {
	CcpResponse struct {
		EmaidContent []struct {
			MessageDef struct {
				CertificateInstallationRes string `json:"certificateInstallationRes"`
				Emaid                      string `json:"emaid"`
			} `json:"messageDef"`
		} `json:"emaidContent"`
	} `json:"CCPResponse"`
	XsdMsgDefNamespace string `json:"xsdMsgDefNamespace"`
}

func (h OpcpMoEvCertificateProvider) ProvideCertificate(exiRequest string) (EvCertificate15118Response, error) {
	client := h.HttpClient
	if client == nil {
		client = http.DefaultClient
	}

	requestUrl := fmt.Sprintf("%s/v1/ccp/signedContractData", h.BaseURL)
	requestBody := SignedContractDataRequest{
		CertificateInstallationReq: exiRequest,
		XsdMsgDefNamespace:         XsdMsgDefinition,
	}
	marshalledBody, err := json.Marshal(requestBody)
	if err != nil {
		return EvCertificate15118Response{}, fmt.Errorf("marshalling body: %w", err)
	}

	resp, err := withRetries(func() (*http.Response, error) {
		req, err := h.moRequest(requestUrl, marshalledBody)
		if err != nil {
			return &http.Response{}, fmt.Errorf("requesting certificate: %w", err)
		}
		return client.Do(req)
	}, 5)

	if err != nil {
		return EvCertificate15118Response{
			Status: ocpp201.Iso15118EVCertificateStatusEnumTypeFailed,
		}, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return EvCertificate15118Response{
			Status: ocpp201.Iso15118EVCertificateStatusEnumTypeFailed,
		}, HttpError(resp.StatusCode)
	}

	responseString, err := io.ReadAll(resp.Body)
	if err != nil {
		return EvCertificate15118Response{
			Status: ocpp201.Iso15118EVCertificateStatusEnumTypeFailed,
		}, err
	}

	var responseBody SignedContractDataResponse
	err = json.Unmarshal(responseString, &responseBody)
	if err != nil {
		return EvCertificate15118Response{}, err
	}

	emaidContents := responseBody.CcpResponse.EmaidContent

	if len(emaidContents) == 0 {
		return EvCertificate15118Response{
			Status: ocpp201.Iso15118EVCertificateStatusEnumTypeFailed,
		}, errors.New("empty emaidContent array")
	}
	if len(emaidContents[0].MessageDef.Emaid) == 0 {
		return EvCertificate15118Response{
			Status: ocpp201.Iso15118EVCertificateStatusEnumTypeFailed,
		}, errors.New("no emaid found")
	}

	response := EvCertificate15118Response{
		Status:                     ocpp201.Iso15118EVCertificateStatusEnumTypeAccepted,
		CertificateInstallationRes: emaidContents[0].MessageDef.CertificateInstallationRes,
	}

	return response, nil
}

func (h OpcpMoEvCertificateProvider) moRequest(requestUrl string, marshalledBody []byte) (*http.Request, error) {
	req, err := http.NewRequest("POST", requestUrl, bytes.NewReader(marshalledBody))
	if err != nil {
		return nil, err
	}
	req.Header.Add("authorization", fmt.Sprintf("Bearer %s", h.BearerToken))
	req.Header.Add("content-type", "application/json")
	return req, nil
}

type retryFunc func() (*http.Response, error)

func withRetries(action retryFunc, attempts int) (*http.Response, error) {
	var lastErr error

	for attempt := 1; attempt <= attempts; attempt++ {
		resp, err := action()
		if err == nil && resp.StatusCode == http.StatusOK {
			// Successful operation, return the response
			return resp, nil
		} else {
			if err == nil {
				err = HttpError(resp.StatusCode)
			}
		}
		if attempt == attempts {
			lastErr = err
			break
		}
	}

	return &http.Response{}, lastErr
}

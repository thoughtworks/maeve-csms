// SPDX-License-Identifier: Apache-2.0

package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
	"io"
	"net/http"
)

const XsdMsgDefinition = "urn:iso:15118:2:2013:MsgDef"

type EvCertificateProvider interface {
	ProvideCertificate(ctx context.Context, exiRequest string) (EvCertificate15118Response, error)
}

type OpcpMoEvCertificateProvider struct {
	BaseURL     string
	BearerToken string
	HttpClient  *http.Client
	Tracer      trace.Tracer
}

type EvCertificate15118Response struct {
	Status                     ocpp201.Iso15118EVCertificateStatusEnumType
	CertificateInstallationRes string
}

type SignedContractDataRequest struct {
	CertificateInstallationReq string `json:"certificateInstallationReq"`
	XsdMsgDefNamespace         string `json:"xsdMsgDefNamespace"`
}

type CcpResponse struct {
	EmaidContent []struct {
		MessageDef struct {
			CertificateInstallationRes string `json:"certificateInstallationRes"`
			Emaid                      string `json:"emaid"`
		} `json:"messageDef"`
	} `json:"emaidContent"`
}

type SignedContractDataResponse struct {
	CcpResponse        CcpResponse `json:"CCPResponse"`
	XsdMsgDefNamespace string      `json:"xsdMsgDefNamespace"`
}

func (h OpcpMoEvCertificateProvider) ProvideCertificate(ctx context.Context, exiRequest string) (EvCertificate15118Response, error) {
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

	resp, err := withRetries(ctx, h.Tracer, func(fnCtx context.Context) (*http.Response, error) {
		req, err := h.moRequest(fnCtx, requestUrl, marshalledBody)
		if err != nil {
			return &http.Response{}, fmt.Errorf("requesting certificate: %w", err)
		}

		resp, err := client.Do(req)

		return resp, err
	}, 3)

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

func (h OpcpMoEvCertificateProvider) moRequest(ctx context.Context, requestUrl string, marshalledBody []byte) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", requestUrl, bytes.NewReader(marshalledBody))
	if err != nil {
		return nil, err
	}
	req.Header.Add("x-api-key", h.BearerToken)
	req.Header.Add("authorization", fmt.Sprintf("Bearer %s", h.BearerToken))
	req.Header.Add("content-type", "application/json")

	return req, nil
}

type retryFunc func(context.Context) (*http.Response, error)

func withRetries(ctx context.Context, tracer trace.Tracer, action retryFunc, attempts int) (*http.Response, error) {
	newCtx, span := tracer.Start(ctx, "get_signed_contract_data")
	defer span.End()

	var lastErr error

	for attempt := 1; attempt <= attempts; attempt++ {
		resp, err := action(newCtx)
		if err == nil && resp.StatusCode == http.StatusOK {
			if attempt > 1 {
				span.SetAttributes(semconv.HTTPResendCount(attempt - 1))
			}
			// Successful operation, return the response
			return resp, nil
		} else {
			_ = resp.Body.Close()
			if err == nil {
				err = HttpError(resp.StatusCode)
			}
		}
		if attempt == attempts {
			lastErr = err
			span.SetAttributes(semconv.HTTPResendCount(attempts - 1))
			span.SetStatus(codes.Error, "retries exhausted")
			span.RecordError(err)
			break
		}
	}

	return &http.Response{}, lastErr
}

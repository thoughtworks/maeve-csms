// SPDX-License-Identifier: Apache-2.0

package services_test

import (
	"context"
	"encoding/json"
	"fmt"
	"go.opentelemetry.io/otel/trace"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"

	"github.com/stretchr/testify/assert"
	"github.com/thoughtworks/maeve-csms/manager/services"
)

const dummyExiResponse = "dummy exi response"
const maxRetryAttempts = 3

type otherHubjectHttpHandler struct {
	flakyCount int
}

func newOtherHubjectHttpHandler() otherHubjectHttpHandler {
	return otherHubjectHttpHandler{flakyCount: 0}
}

func (h *otherHubjectHttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost || r.URL.String() != "/v1/ccp/signedContractData" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	if r.Header.Get("authorization") != "Bearer TestToken" {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	if r.Header.Get("content-type") != "application/json" {
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

	var request services.SignedContractDataRequest
	err = json.Unmarshal(enc, &request)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if request.CertificateInstallationReq == "invalid" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	h.flakyCount++
	if request.CertificateInstallationReq == "flaky" && h.flakyCount < maxRetryAttempts {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	responseString := certificateInstallationResponse(request.CertificateInstallationReq, dummyExiResponse)
	w.Write([]byte(responseString))
}

func certificateInstallationResponse(request, exiResponse string) string {
	switch request {
	case "empty response":
		return `{"CCPResponse":{"emaidContent":[]}}`
	case "invalid response":
		return fmt.Sprintf(`{"CCPResponse":{"emaidContent":[{"messageDef":{"certificateInstallationRes":"%s"}}]}}`, exiResponse)
	default:
		return fmt.Sprintf(`{"CCPResponse":{"emaidContent":[{"messageDef":{"certificateInstallationRes":"%s","emaid":"EMP77TWTW99999"}}]}}`, exiResponse)
	}
}

func TestEvCertificateProvider(t *testing.T) {
	tracer := trace.NewNoopTracerProvider().Tracer("test")

	hubject := newOtherHubjectHttpHandler()

	server := httptest.NewServer(&hubject)
	defer server.Close()

	provider := services.OpcpMoEvCertificateProvider{
		BaseURL:     server.URL,
		BearerToken: "TestToken",
		Tracer:      tracer,
	}

	response, err := provider.ProvideCertificate(context.Background(), "valid")

	assert.NoError(t, err)
	assert.Equal(t, ocpp201.Iso15118EVCertificateStatusEnumTypeAccepted, response.Status)
	assert.Equal(t, dummyExiResponse, response.CertificateInstallationRes)
}

func TestEvCertificateProviderWithFlakyResponses(t *testing.T) {
	tracer := trace.NewNoopTracerProvider().Tracer("test")

	hubject := newOtherHubjectHttpHandler()

	server := httptest.NewServer(&hubject)
	defer server.Close()

	provider := services.OpcpMoEvCertificateProvider{
		BaseURL:     server.URL,
		BearerToken: "TestToken",
		Tracer:      tracer,
	}

	response, err := provider.ProvideCertificate(context.Background(), "flaky")

	assert.NoError(t, err)
	assert.Equal(t, ocpp201.Iso15118EVCertificateStatusEnumTypeAccepted, response.Status)
	assert.Equal(t, dummyExiResponse, response.CertificateInstallationRes)
}

func TestEvCertificateProviderFailureCases(t *testing.T) {
	tracer := trace.NewNoopTracerProvider().Tracer("test")

	tests := map[string]struct {
		token      string
		exiRequest string
	}{
		"wrong token":            {token: "WrongToken", exiRequest: "valid"},
		"invalid exi request":    {token: "TestToken", exiRequest: "invalid"},
		"empty response":         {token: "TestToken", exiRequest: "empty response"},
		"invalid emaid response": {token: "TestToken", exiRequest: "invalid response"},
	}

	hubject := newOtherHubjectHttpHandler()
	server := httptest.NewServer(&hubject)
	defer server.Close()

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			provider := services.OpcpMoEvCertificateProvider{
				BaseURL:     server.URL,
				BearerToken: tc.token,
				Tracer:      tracer,
			}

			response, err := provider.ProvideCertificate(context.Background(), tc.exiRequest)

			assert.Error(t, err)
			assert.Equal(t, ocpp201.Iso15118EVCertificateStatusEnumTypeFailed, response.Status)
		})
	}
}

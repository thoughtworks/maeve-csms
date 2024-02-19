package server_test

import (
	"crypto/x509"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/thoughtworks/maeve-csms/gateway/registry"
	"github.com/thoughtworks/maeve-csms/gateway/server"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTraceMatchedRequest(t *testing.T) {
	tracer, traceExporter := server.GetTracer()

	r := chi.NewRouter()
	r.Use(server.TraceRequest(tracer))
	r.HandleFunc("/id/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/id/1234", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	server.AssertSpan(t, &traceExporter.GetSpans()[0], "GET /id/{id}", map[string]any{
		"http.scheme": "ws",
		"http.method": "GET",
		"http.url":    "/id/1234",
		"http.route":  "/id/{id}",
	})
}

func TestTraceUnmatchedRequest(t *testing.T) {
	tracer, traceExporter := server.GetTracer()

	r := chi.NewRouter()
	r.Use(server.TraceRequest(tracer))
	r.HandleFunc("/something", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/other", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	server.AssertSpan(t, &traceExporter.GetSpans()[0], "GET /other", map[string]any{
		"http.scheme":      "ws",
		"http.method":      "GET",
		"http.url":         "/other",
		"http.status_code": http.StatusNotFound,
	})
}

func TestTLSOffloadWithNoClientCert(t *testing.T) {
	r := chi.NewRouter()
	r.Use(server.TLSOffload(registry.NewMockRegistry()))
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.TLS != nil {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Forwarded-Proto", "https")

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Result().StatusCode)

}

func TestTLSOffloadWithValidClientCertificate(t *testing.T) {
	r := chi.NewRouter()
	reg := registry.NewMockRegistry()
	reg.Certificates["certificate-hash"] = &x509.Certificate{}
	r.Use(server.TLSOffload(reg))
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.TLS != nil && r.TLS.PeerCertificates != nil && len(r.TLS.PeerCertificates) > 0 {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Forwarded-Proto", "https")
	req.Header.Set("X-Client-Cert-Present", "true")
	req.Header.Set("X-Client-Cert-Chain-Verified", "true")
	req.Header.Set("X-Client-Cert-Hash", "certificate-hash")

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Result().StatusCode)
}

func TestTLSOffloadWithInvalidClientCertificate(t *testing.T) {
	r := chi.NewRouter()
	reg := registry.NewMockRegistry()
	reg.Certificates["certificate-hash"] = &x509.Certificate{}
	r.Use(server.TLSOffload(reg))
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.TLS != nil && r.TLS.PeerCertificates != nil && len(r.TLS.PeerCertificates) > 0 {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Forwarded-Proto", "https")
	req.Header.Set("X-Client-Cert-Present", "true")
	req.Header.Set("X-Client-Cert-Chain-Verified", "false")
	req.Header.Set("X-Client-Cert-Hash", "certificate-hash")

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Result().StatusCode)
}

func TestTLSOffloadWithUnknownClientCertificate(t *testing.T) {
	r := chi.NewRouter()
	reg := registry.NewMockRegistry()
	r.Use(server.TLSOffload(reg))
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.TLS != nil && r.TLS.PeerCertificates != nil && len(r.TLS.PeerCertificates) > 0 {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Forwarded-Proto", "https")
	req.Header.Set("X-Client-Cert-Present", "true")
	req.Header.Set("X-Client-Cert-Chain-Verified", "true")
	req.Header.Set("X-Client-Cert-Hash", "certificate-hash")

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Result().StatusCode)
}

type errorRegistry struct{}

func (e errorRegistry) LookupChargeStation(clientId string) (*registry.ChargeStation, error) {
	return nil, fmt.Errorf("expected error")
}

func (e errorRegistry) LookupCertificate(certHash string) (*x509.Certificate, error) {
	return nil, fmt.Errorf("expected error")
}

func TestTLSOffloadWithClientCertificateRetrievalError(t *testing.T) {
	r := chi.NewRouter()
	reg := errorRegistry{}
	r.Use(server.TLSOffload(reg))
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.TLS != nil && r.TLS.PeerCertificates != nil && len(r.TLS.PeerCertificates) > 0 {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Forwarded-Proto", "https")
	req.Header.Set("X-Client-Cert-Present", "true")
	req.Header.Set("X-Client-Cert-Chain-Verified", "true")
	req.Header.Set("X-Client-Cert-Hash", "certificate-hash")

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Result().StatusCode)
}

package server

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"strconv"

	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"

	"github.com/thoughtworks/maeve-csms/gateway/registry"
	"golang.org/x/exp/slog"
)

func TraceRequest(tracer trace.Tracer) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			slog.Info("websocket connection received", "path", r.URL.Path, "method", r.Method)
			slog.Info("processing connection", "uri", r.RequestURI)

			newCtx, span := tracer.Start(r.Context(), fmt.Sprintf("%s %s", r.Method, r.URL.String()), trace.WithSpanKind(trace.SpanKindServer),
				trace.WithAttributes(
					semconv.HTTPScheme(getScheme(r)),
					semconv.HTTPMethod(r.Method),
					semconv.HTTPURL(r.URL.String())))
			defer span.End()

			h.ServeHTTP(w, r.WithContext(newCtx))

			routePattern := chi.RouteContext(r.Context()).RoutePattern()
			if routePattern != "" {
				span.SetName(fmt.Sprintf("%s %s", r.Method, routePattern))
			} else {
				span.SetStatus(codes.Error, "not found")
				span.SetAttributes(semconv.HTTPStatusCode(http.StatusNotFound))
			}
			span.SetAttributes(semconv.HTTPRoute(chi.RouteContext(r.Context()).RoutePattern()))
		})
	}
}

func TLSOffload(registry registry.DeviceRegistry) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			span := trace.SpanFromContext(r.Context())

			forwardedProtoHeader := r.Header.Get("X-Forwarded-Proto")
			span.SetAttributes(attribute.String("http.proto", forwardedProtoHeader))

			if forwardedProtoHeader == "https" {
				r.TLS = &tls.ConnectionState{
					HandshakeComplete: true,
				}

				clientCertPresentHeader := r.Header.Get("X-Client-Cert-Present")
				clientCertPresent, err := strconv.ParseBool(clientCertPresentHeader)
				span.SetAttributes(attribute.Bool("cert.present", clientCertPresent))
				if err == nil && clientCertPresent {
					clientCertChainValidHeader := r.Header.Get("X-Client-Cert-Chain-Verified")
					clientCertChainValid, err := strconv.ParseBool(clientCertChainValidHeader)
					span.SetAttributes(attribute.Bool("cert.valid", clientCertChainValid))
					if err == nil && clientCertChainValid {
						clientCertHashHeader := r.Header.Get("X-Client-Cert-Hash")
						span.SetAttributes(attribute.String("cert.hash", clientCertHashHeader))
						certificate, err := registry.LookupCertificate(clientCertHashHeader)
						if err == nil && certificate != nil {
							r.TLS.PeerCertificates = []*x509.Certificate{certificate}
						} else if err != nil {
							span.SetAttributes(attribute.String("cert.lookup.error", err.Error()))
							slog.Error("lookup certificate", "clientCertHashHeader", clientCertHashHeader, "err", err)
						} else {
							span.SetAttributes(attribute.String("cert.lookup.error", "NotFound"))
							slog.Warn("certificate not found", "clientCertHashHeader", clientCertHashHeader)
						}
					}
				}
			}

			h.ServeHTTP(w, r)
		})
	}
}

package server

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/thoughtworks/maeve-csms/gateway/registry"
	"log"
	"net/http"
	"strconv"
)

func TLSOffload(registry registry.DeviceRegistry) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			forwardedProtoHeader := r.Header.Get("X-Forwarded-Proto")
			if forwardedProtoHeader == "https" {
				r.TLS = &tls.ConnectionState{
					HandshakeComplete: true,
				}

				clientCertPresentHeader := r.Header.Get("X-Client-Cert-Present")
				clientCertPresent, err := strconv.ParseBool(clientCertPresentHeader)
				if err == nil && clientCertPresent {
					clientCertChainValidHeader := r.Header.Get("X-Client-Cert-Chain-Verified")
					clientCertChainValid, err := strconv.ParseBool(clientCertChainValidHeader)
					if err == nil && clientCertChainValid {
						clientCertHashHeader := r.Header.Get("X-Client-Cert-Hash")
						certificate, err := registry.LookupCertificate(clientCertHashHeader)
						if err == nil && certificate != nil {
							r.TLS.PeerCertificates = []*x509.Certificate{certificate}
						} else if err != nil {
							log.Printf("lookup certificate %s: %v", clientCertHashHeader, err)
						} else {
							log.Printf("certificate for %s not found", clientCertHashHeader)
						}
					}
				}
			}

			h.ServeHTTP(w, r)
		})
	}
}

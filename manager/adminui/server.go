package adminui

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"embed"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/thoughtworks/maeve-csms/manager/services"
	"github.com/thoughtworks/maeve-csms/manager/store"
	"golang.org/x/exp/slog"
	"html/template"
	"net/http"
)

var (
	//go:embed templates
	res embed.FS
)

func NewServer(externalAddr, orgName string, engine store.Engine, certificateProvider services.ChargeStationCertificateProvider) chi.Router {
	r := chi.NewRouter()

	templates := template.Must(template.ParseFS(res, "templates/*.gohtml"))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		err := templates.ExecuteTemplate(w, "index.gohtml", nil)
		if err != nil {
			slog.Error("rendering template", "err", err)
			_ = templates.ExecuteTemplate(w, "error.gohtml", nil)
		}
	})

	r.Get("/connect", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		err := templates.ExecuteTemplate(w, "connect.gohtml", nil)
		if err != nil {
			slog.Error("rendering template", "err", err)
			_ = templates.ExecuteTemplate(w, "error.gohtml", nil)
		}
	})

	r.Post("/connect", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			slog.Error("parsing form", "err", err)
			_ = templates.ExecuteTemplate(w, "error.gohtml", nil)
			return
		}

		csId := r.PostFormValue("csid")
		auth := r.PostFormValue("auth")

		if csId == "" || auth == "" {
			slog.Error("missing form parameters")
			_ = templates.ExecuteTemplate(w, "error.gohtml", nil)
			return
		}

		var data map[string]string

		switch auth {
		case "unsecured":
			password, err := createPassword()
			if err != nil {
				slog.Error("creating password", "err", err)
				_ = templates.ExecuteTemplate(w, "error.gohtml", nil)
				return
			}
			err = registerChargeStation(r.Context(), engine, csId, 0, password)
			if err != nil {
				slog.Error("registering charge station", "err", err)
				_ = templates.ExecuteTemplate(w, "error.gohtml", nil)
				return
			}
			data = map[string]string{
				"csid":     csId,
				"auth":     auth,
				"url":      fmt.Sprintf("ws://%s/ws/%s", externalAddr, csId),
				"password": password,
			}
		case "basic":
			password, err := createPassword()
			if err != nil {
				slog.Error("creating password", "err", err)
				_ = templates.ExecuteTemplate(w, "error.gohtml", nil)
				return
			}
			err = registerChargeStation(r.Context(), engine, csId, 1, password)
			if err != nil {
				slog.Error("registering charge station", "err", err)
				_ = templates.ExecuteTemplate(w, "error.gohtml", nil)
				return
			}
			data = map[string]string{
				"csid":     csId,
				"auth":     auth,
				"url":      fmt.Sprintf("wss://%s/ws/%s", externalAddr, csId),
				"password": password,
			}
		case "mtls":
			clientKey, clientCert, err := createSignedKeyPair(r.Context(), csId, orgName, certificateProvider)
			if err != nil {
				slog.Error("creating signed key pair", "err", err)
				_ = templates.ExecuteTemplate(w, "error.gohtml", nil)
				return
			}
			err = registerChargeStation(r.Context(), engine, csId, 2, "")
			if err != nil {
				slog.Error("registering charge station", "err", err)
				_ = templates.ExecuteTemplate(w, "error.gohtml", nil)
				return
			}
			data = map[string]string{
				"csid":       csId,
				"auth":       auth,
				"url":        fmt.Sprintf("wss://%s/ws/%s", externalAddr, csId),
				"clientCert": clientCert,
				"clientKey":  clientKey,
			}
		}

		w.Header().Set("Content-Type", "text/html")
		err = templates.ExecuteTemplate(w, "post-connect.gohtml", data)
		if err != nil {
			slog.Error("rendering template", "err", err)
			_ = templates.ExecuteTemplate(w, "error.gohtml", nil)
		}
	})

	return r
}

func createPassword() (string, error) {
	const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	var randData = make([]byte, 10)

	_, err := rand.Reader.Read(randData)
	if err != nil {
		return "", err
	} else {
		for i := 0; i < len(randData); i++ {
			randData[i] = alphabet[int(randData[i])%len(alphabet)]
		}
	}

	return string(randData), nil
}

func createSignedKeyPair(ctx context.Context, csId string, orgName string, certificateProvider services.ChargeStationCertificateProvider) (string, string, error) {
	keyPair, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", "", fmt.Errorf("generating rsa key: %v", err)
	}

	csrTemplate := x509.CertificateRequest{
		Subject: pkix.Name{
			CommonName:   csId,
			Organization: []string{orgName},
		},
		SignatureAlgorithm: x509.SHA256WithRSA,
	}

	csr, err := x509.CreateCertificateRequest(rand.Reader, &csrTemplate, keyPair)
	if err != nil {
		return "", "", fmt.Errorf("creating certificate request: %v", err)
	}

	pemCsr := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: csr,
	})

	chain, err := certificateProvider.ProvideCertificate(ctx, services.CertificateTypeCSO, string(pemCsr), csId)
	if err != nil {
		return "", "", fmt.Errorf("providing certificate: %v", err)
	}

	pemKey := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(keyPair),
	})

	return string(pemKey), chain, nil
}

func registerChargeStation(ctx context.Context, engine store.Engine, csId string, scheme int, password string) error {
	var profile store.SecurityProfile

	switch scheme {
	case 0:
		profile = store.UnsecuredTransportWithBasicAuth
	case 1:
		profile = store.TLSWithBasicAuth
	case 2:
		profile = store.TLSWithClientSideCertificates
	default:
		return fmt.Errorf("unknown security profile: %d", scheme)
	}

	var b64sha256 = ""
	if password != "" {
		sha256pw := sha256.Sum256([]byte(password))
		b64sha256 = base64.StdEncoding.EncodeToString(sha256pw[:])
	}

	err := engine.SetChargeStationAuth(ctx, csId, &store.ChargeStationAuth{
		SecurityProfile:      profile,
		Base64SHA256Password: b64sha256,
	})
	if err != nil {
		return err
	}

	return nil
}

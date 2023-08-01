// SPDX-License-Identifier: Apache-2.0

package server

import (
	"fmt"
	"golang.org/x/exp/slog"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/rs/cors"
	"github.com/thoughtworks/maeve-csms/manager/api"
	"github.com/thoughtworks/maeve-csms/manager/store"
	"github.com/unrolled/secure"
	"k8s.io/utils/clock"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/thoughtworks/maeve-csms/manager/templates"
)

type logEntry struct {
	method     string
	path       string
	remoteAddr string
}

func (l logEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	slog.Info(fmt.Sprintf("%s %s", l.method, l.path), slog.String("remote_addr", l.remoteAddr), slog.Int("status", status), slog.Int("bytes", bytes), slog.Duration("duration", elapsed))
}

func (l logEntry) Panic(v interface{}, stack []byte) {
	slog.Info(fmt.Sprintf("%s %s", l.method, l.path), slog.String("remote_addr", l.remoteAddr), slog.Any("panic", v), slog.String("stack", string(stack)))
}

type logFormatter struct{}

func (l logFormatter) NewLogEntry(r *http.Request) middleware.LogEntry {
	return logEntry{
		method:     r.Method,
		path:       r.URL.Path,
		remoteAddr: r.RemoteAddr,
	}
}

func NewApiHandler(engine store.Engine) http.Handler {
	apiServer, err := api.NewServer(engine, clock.RealClock{})
	if err != nil {
		panic(err)
	}

	var isDevelopment bool
	if os.Getenv("ENVIRONMENT") == "dev" {
		isDevelopment = true
	}
	secureMiddleware := secure.New(secure.Options{
		IsDevelopment:         isDevelopment,
		BrowserXssFilter:      true,
		ContentTypeNosniff:    true,
		FrameDeny:             true,
		ContentSecurityPolicy: "frame-ancestors: 'none'",
	})

	r := chi.NewRouter()

	logger := middleware.RequestLogger(logFormatter{})

	r.Use(middleware.Recoverer, secureMiddleware.Handler, cors.Default().Handler, api.ValidationMiddleware)
	r.Get("/health", health)
	r.Get("/transactions", transactions(engine))
	r.Handle("/metrics", promhttp.Handler())
	r.Get("/api/openapi.json", getSwaggerJson)
	r.With(logger).Mount("/api/v0", api.Handler(apiServer))
	return r
}

func getSwaggerJson(w http.ResponseWriter, r *http.Request) {
	swagger, err := api.GetSwagger()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	json, err := swagger.MarshalJSON()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(json)
}

func health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"OK"}`))
}

func transactions(transactionStore store.Engine) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ts, err := transactionStore.Transactions(r.Context())
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		t := template.Must(template.New("").Parse(templates.Transactions))
		err = t.Execute(w, ts)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

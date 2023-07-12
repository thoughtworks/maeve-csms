// SPDX-License-Identifier: Apache-2.0

package server

import (
	"github.com/rs/cors"
	"github.com/thoughtworks/maeve-csms/manager/api"
	"github.com/thoughtworks/maeve-csms/manager/store"
	"github.com/unrolled/secure"
	"k8s.io/utils/clock"
	"net/http"
	"os"
	"text/template"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/thoughtworks/maeve-csms/manager/services"
	"github.com/thoughtworks/maeve-csms/manager/templates"
)

func NewApiHandler(engine store.Engine, transactionStore services.TransactionStore) http.Handler {
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
	r.Use(middleware.Logger, middleware.Recoverer, secureMiddleware.Handler, cors.Default().Handler, api.ValidationMiddleware)
	r.Get("/health", health)
	r.Get("/transactions", transactions(transactionStore))
	r.Handle("/metrics", promhttp.Handler())
	r.Get("/api/openapi.json", getSwaggerJson)
	r.Mount("/api/v0", api.Handler(apiServer))
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

func transactions(transactionStore services.TransactionStore) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ts, err := transactionStore.Transactions()
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

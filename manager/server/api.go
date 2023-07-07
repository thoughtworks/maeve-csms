// SPDX-License-Identifier: Apache-2.0

package server

import (
	"github.com/twlabs/maeve-csms/manager/api"
	"github.com/twlabs/maeve-csms/manager/store"
	"net/http"
	"text/template"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/thoughtworks/maeve-csms/manager/services"
	"github.com/thoughtworks/maeve-csms/manager/templates"
)

func NewApiHandler(engine store.Engine, transactionStore services.TransactionStore) http.Handler {
	csm := api.ChargeStationManager{
		ChargeStationAuthStore: engine,
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger, middleware.Recoverer)
	r.Get("/health", health)
	r.Get("/transactions", transactions(transactionStore))
	r.Handle("/metrics", promhttp.Handler())
	r.Route("/api/v0", func(r chi.Router) {
		// TODO: add secure middleware
		r.Route("/cs/{csId}", func(r chi.Router) {
			r.Post("/", csm.CreateChargeStation)
			r.Get("/auth", csm.RetrieveChargeStationAuthDetails)
		})
	})
	return r
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

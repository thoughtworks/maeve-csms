package server

import (
	"net/http"
	"text/template"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/twlabs/ocpp2-broker-core/manager/services"
	"github.com/twlabs/ocpp2-broker-core/manager/templates"
)

func NewApiHandler(transactionStore services.TransactionStore) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger, middleware.Recoverer)
	r.Get("/health", health)
	r.Get("/transactions", transactions(transactionStore))
	r.Handle("/metrics", promhttp.Handler())
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

// SPDX-License-Identifier: Apache-2.0

package server

import (
	oapimiddleware "github.com/deepmap/oapi-codegen/pkg/chi-middleware"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
	"github.com/thoughtworks/maeve-csms/manager/handlers/ocpp16"
	"github.com/thoughtworks/maeve-csms/manager/ocpi"
	"github.com/thoughtworks/maeve-csms/manager/store"
	"github.com/thoughtworks/maeve-csms/manager/transport"
	"github.com/unrolled/secure"
	"k8s.io/utils/clock"
	"net/http"
	"os"
)

func NewOcpiHandler(engine store.Engine, clock clock.PassiveClock, ocpiApi ocpi.Api, emitter transport.Emitter) http.Handler {
	v16CallMaker := ocpp16.NewCallMaker(emitter)
	ocpiServer, err := ocpi.NewServer(ocpiApi, clock, v16CallMaker)
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

	swagger, err := ocpi.GetSwagger()
	if err != nil {
		panic(err)
	}
	swagger.Servers = nil
	r.Use(middleware.Recoverer, secureMiddleware.Handler, cors.Default().Handler, logger)
	r.Get("/openapi.json", getOcpiSwaggerJson)
	r.With(oapimiddleware.OapiRequestValidatorWithOptions(swagger, &oapimiddleware.Options{
		Options: openapi3filter.Options{
			AuthenticationFunc: ocpi.NewTokenAuthenticationFunc(engine),
		},
	})).Mount("/", ocpi.Handler(ocpiServer))

	return r
}

func getOcpiSwaggerJson(w http.ResponseWriter, r *http.Request) {
	swagger, err := ocpi.GetSwagger()
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

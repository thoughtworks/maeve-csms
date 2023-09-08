// SPDX-License-Identifier: Apache-2.0

package ocpi_test

import (
	"context"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thoughtworks/maeve-csms/manager/ocpi"
	"github.com/thoughtworks/maeve-csms/manager/store"
	"github.com/thoughtworks/maeve-csms/manager/store/inmemory"
	"k8s.io/utils/clock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthenticationWithPendingToken(t *testing.T) {
	token := "abcdef123456"

	engine := inmemory.NewStore(clock.RealClock{})
	err := engine.SetRegistrationDetails(context.Background(), token, &store.OcpiRegistration{Status: store.OcpiRegistrationStatusPending})
	require.NoError(t, err)

	authFn := ocpi.NewTokenAuthenticationFunc(engine)

	endpoints := []struct {
		Path    string
		Method  string
		Success bool
	}{
		{Path: "/ocpi/versions", Method: http.MethodGet, Success: true},
		{Path: "/ocpi/2.2", Method: http.MethodGet, Success: true},
		{Path: "/ocpi/2.2/credentials", Method: http.MethodPost, Success: true},
		{Path: "/ocpi/2.2/credentials", Method: http.MethodGet, Success: false},
		{Path: "/ocpi/2.2/credentials", Method: http.MethodPut, Success: false},
		{Path: "/ocpi/2.2/credentials", Method: http.MethodDelete, Success: false},
	}

	for _, endpoint := range endpoints {
		t.Run(fmt.Sprintf("%s %s", endpoint.Method, endpoint.Path), func(t *testing.T) {
			req := httptest.NewRequest(endpoint.Method, endpoint.Path, nil)
			req.Header.Add("Authorization", fmt.Sprintf("Token %s", token))
			err = authFn(context.Background(), &openapi3filter.AuthenticationInput{
				RequestValidationInput: &openapi3filter.RequestValidationInput{
					Request: req,
				},
				SecuritySchemeName: "token",
				SecurityScheme: &openapi3.SecurityScheme{
					Type:   "http",
					Scheme: "bearer",
				},
			})
			if endpoint.Success {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, "authorization failed: unregistered token")
			}
		})
	}
}

func TestAuthenticationWithRegisteredToken(t *testing.T) {
	token := "abcdef123456"

	engine := inmemory.NewStore(clock.RealClock{})
	err := engine.SetRegistrationDetails(context.Background(), token, &store.OcpiRegistration{Status: store.OcpiRegistrationStatusRegistered})
	require.NoError(t, err)

	authFn := ocpi.NewTokenAuthenticationFunc(engine)

	endpoints := []struct {
		Path    string
		Method  string
		Success bool
	}{
		{Path: "/ocpi/versions", Method: http.MethodGet, Success: true},
		{Path: "/ocpi/2.2", Method: http.MethodGet, Success: true},
		{Path: "/ocpi/2.2/credentials", Method: http.MethodPost, Success: true},
		{Path: "/ocpi/2.2/credentials", Method: http.MethodGet, Success: true},
		{Path: "/ocpi/2.2/credentials", Method: http.MethodPut, Success: true},
		{Path: "/ocpi/2.2/credentials", Method: http.MethodDelete, Success: true},
	}

	for _, endpoint := range endpoints {
		t.Run(fmt.Sprintf("%s %s", endpoint.Method, endpoint.Path), func(t *testing.T) {
			req := httptest.NewRequest(endpoint.Method, endpoint.Path, nil)
			req.Header.Add("Authorization", fmt.Sprintf("Token %s", token))
			err = authFn(context.Background(), &openapi3filter.AuthenticationInput{
				RequestValidationInput: &openapi3filter.RequestValidationInput{
					Request: req,
				},
				SecuritySchemeName: "token",
				SecurityScheme: &openapi3.SecurityScheme{
					Type:   "http",
					Scheme: "bearer",
				},
			})
			if endpoint.Success {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, "authorization failed: unregistered token")
			}
		})
	}
}

func TestAuthenticationWithUnknownToken(t *testing.T) {
	token := "abcdef123456"

	engine := inmemory.NewStore(clock.RealClock{})

	authFn := ocpi.NewTokenAuthenticationFunc(engine)

	req := httptest.NewRequest(http.MethodGet, "/ocpi/versions", nil)
	req.Header.Add("Authorization", fmt.Sprintf("Token %s", token))

	err := authFn(context.Background(), &openapi3filter.AuthenticationInput{
		RequestValidationInput: &openapi3filter.RequestValidationInput{
			Request: req,
		},
		SecuritySchemeName: "token",
		SecurityScheme: &openapi3.SecurityScheme{
			Type:   "http",
			Scheme: "bearer",
		},
	})

	assert.ErrorContains(t, err, "authorization failed: unknown token")
}

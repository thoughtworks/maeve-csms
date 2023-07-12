package api_test

import (
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thoughtworks/maeve-csms/manager/api"
	"github.com/thoughtworks/maeve-csms/manager/store/inmemory"
	"io"
	clockTest "k8s.io/utils/clock/testing"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestValidationMiddlewareWithInvalidRequest(t *testing.T) {
	engine := inmemory.NewStore()

	now := time.Now()
	srv, err := api.NewServer(engine, clockTest.NewFakePassiveClock(now))
	require.NoError(t, err)

	r := chi.NewRouter()
	r.Use(api.ValidationMiddleware)
	r.Mount("/", api.Handler(srv))
	server := httptest.NewServer(r)
	defer server.Close()

	req := httptest.NewRequest(http.MethodPost, "/cs/cs001", strings.NewReader(""))
	req.Header.Set("content-type", "application/json")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Result().StatusCode)
	b, err := io.ReadAll(rr.Result().Body)
	require.NoError(t, err)
	assert.JSONEq(t, `{"status":"Bad Request","error":"request body has an error: value is required but missing"}`, string(b))
}

func TestValidationMiddlewareWithValidRequest(t *testing.T) {
	engine := inmemory.NewStore()

	now := time.Now()
	srv, err := api.NewServer(engine, clockTest.NewFakePassiveClock(now))
	require.NoError(t, err)

	r := chi.NewRouter()
	r.Use(api.ValidationMiddleware)
	r.Mount("/", api.Handler(srv))
	server := httptest.NewServer(r)
	defer server.Close()

	req := httptest.NewRequest(http.MethodPost, "/cs/cs001", strings.NewReader(`{"securityProfile":0}`))
	req.Header.Set("content-type", "application/json")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Result().StatusCode)
	b, err := io.ReadAll(rr.Result().Body)
	require.NoError(t, err)
	assert.Equal(t, "", string(b))
}

func TestValidationMiddlewareWithNonApiPath(t *testing.T) {
	r := chi.NewRouter()
	r.Use(api.ValidationMiddleware)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"status":"OK"}`))
	})
	server := httptest.NewServer(r)
	defer server.Close()

	req := httptest.NewRequest(http.MethodGet, "/health", strings.NewReader(""))
	req.Header.Set("content-type", "application/json")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Result().StatusCode)
}

func TestValidationMiddlewareWithUnknownMethod(t *testing.T) {
	engine := inmemory.NewStore()

	now := time.Now()
	srv, err := api.NewServer(engine, clockTest.NewFakePassiveClock(now))
	require.NoError(t, err)

	r := chi.NewRouter()
	r.Use(api.ValidationMiddleware)
	r.Mount("/", api.Handler(srv))
	server := httptest.NewServer(r)
	defer server.Close()

	req := httptest.NewRequest(http.MethodGet, "/cs/cs001", strings.NewReader(""))
	req.Header.Set("content-type", "application/json")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Result().StatusCode)
}

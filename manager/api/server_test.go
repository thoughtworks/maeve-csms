package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thoughtworks/maeve-csms/manager/api"
	"github.com/thoughtworks/maeve-csms/manager/store"
	"github.com/thoughtworks/maeve-csms/manager/store/inmemory"
	"io"
	clockTest "k8s.io/utils/clock/testing"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestRegisterChargeStation(t *testing.T) {
	server, r, _, _ := setupServer(t)
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

func TestLookupChargeStationAuth(t *testing.T) {
	server, r, engine, _ := setupServer(t)
	defer server.Close()

	err := engine.SetChargeStationAuth(context.Background(), "cs001", &store.ChargeStationAuth{
		SecurityProfile: 1,
	})
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/cs/cs001/auth", strings.NewReader("{}"))
	req.Header.Set("accept", "application/json")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Result().StatusCode)
	b, err := io.ReadAll(rr.Result().Body)
	require.NoError(t, err)

	want := &api.ChargeStationAuth{
		SecurityProfile: 1,
	}

	got := new(api.ChargeStationAuth)
	err = json.Unmarshal(b, &got)
	require.NoError(t, err)

	assert.Equal(t, want, got)
}

func TestLookupChargeStationAuthThatDoesNotExist(t *testing.T) {
	server, r, _, _ := setupServer(t)
	defer server.Close()

	req := httptest.NewRequest(http.MethodGet, "/cs/unknown/auth", strings.NewReader("{}"))
	req.Header.Set("accept", "application/json")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Result().StatusCode)
}

func TestSetToken(t *testing.T) {
	server, r, engine, now := setupServer(t)
	defer server.Close()

	token := api.Token{
		CacheMode:   "ALWAYS",
		ContractId:  "GB-TWK-012345678-V",
		CountryCode: "GB",
		Issuer:      "Thoughtworks",
		PartyId:     "TWK",
		Type:        "RFID",
		Uid:         "012345678",
		Valid:       true,
	}
	tokenPayload, err := json.Marshal(token)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/token", bytes.NewReader(tokenPayload))
	req.Header.Set("content-type", "application/json")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Result().StatusCode)
	b, err := io.ReadAll(rr.Result().Body)
	require.NoError(t, err)
	assert.Equal(t, "", string(b))

	want := &store.Token{
		CountryCode: "GB",
		PartyId:     "TWK",
		Type:        "RFID",
		Uid:         "012345678",
		ContractId:  "GBTWK012345678V",
		Issuer:      "Thoughtworks",
		Valid:       true,
		CacheMode:   "ALWAYS",
		LastUpdated: now.Format(time.RFC3339),
	}

	got, err := engine.LookupToken(context.Background(), "012345678")
	require.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestLookupToken(t *testing.T) {
	server, r, engine, now := setupServer(t)
	defer server.Close()

	err := engine.SetToken(context.Background(), &store.Token{
		CountryCode: "GB",
		PartyId:     "TWK",
		Type:        "RFID",
		Uid:         "012345678",
		ContractId:  "GBTWK012345678V",
		Issuer:      "Thoughtworks",
		Valid:       true,
		CacheMode:   "ALWAYS",
		LastUpdated: now.Format(time.RFC3339),
	})
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/token/012345678", nil)
	req.Header.Set("accept", "application/json")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Result().StatusCode)
	decoder := json.NewDecoder(rr.Result().Body)
	var got api.Token
	err = decoder.Decode(&got)
	require.NoError(t, err)

	lastUpdatedStr := now.Format(time.RFC3339)
	lastUpdated, err := time.Parse(time.RFC3339, lastUpdatedStr)
	require.NoError(t, err)
	want := api.Token{
		CacheMode:   "ALWAYS",
		ContractId:  "GBTWK012345678V",
		CountryCode: "GB",
		Issuer:      "Thoughtworks",
		LastUpdated: &lastUpdated,
		PartyId:     "TWK",
		Type:        "RFID",
		Uid:         "012345678",
		Valid:       true,
	}

	assert.Equal(t, want, got)
}

func setupServer(t *testing.T) (*httptest.Server, *chi.Mux, store.Engine, time.Time) {
	engine := inmemory.NewStore()

	now := time.Now()
	srv, err := api.NewServer(engine, clockTest.NewFakePassiveClock(now))
	require.NoError(t, err)

	r := chi.NewRouter()
	r.Use(api.ValidationMiddleware)
	r.Mount("/", api.Handler(srv))
	server := httptest.NewServer(r)

	return server, r, engine, now
}

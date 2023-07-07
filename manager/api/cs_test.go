package api_test

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twlabs/maeve-csms/manager/api"
	"github.com/twlabs/maeve-csms/manager/store"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type InMemoryChargeStationAuthStore struct {
	chargeStations map[string]*store.ChargeStationAuth
}

func (m *InMemoryChargeStationAuthStore) SetChargeStationAuth(ctx context.Context, chargeStationId string, auth *store.ChargeStationAuth) error {
	m.chargeStations[chargeStationId] = auth
	return nil
}

func (m *InMemoryChargeStationAuthStore) LookupChargeStationAuth(ctx context.Context, chargeStationId string) (*store.ChargeStationAuth, error) {
	return m.chargeStations[chargeStationId], nil
}

func TestCreateChargeStation(t *testing.T) {
	authStore := new(InMemoryChargeStationAuthStore)
	authStore.chargeStations = make(map[string]*store.ChargeStationAuth)

	csm := api.ChargeStationManager{
		ChargeStationAuthStore: authStore,
	}

	r := chi.NewRouter()
	r.Post("/{csId}", csm.CreateChargeStation)
	server := httptest.NewServer(r)
	defer server.Close()

	req := httptest.NewRequest(http.MethodPost, "/cs001", strings.NewReader("{}"))
	req.Header.Set("content-type", "application/json")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Result().StatusCode)
	b, err := io.ReadAll(rr.Result().Body)
	require.NoError(t, err)
	assert.Equal(t, "", string(b))
}

func TestLookupChargeStation(t *testing.T) {
	authStore := new(InMemoryChargeStationAuthStore)
	authStore.chargeStations = make(map[string]*store.ChargeStationAuth)
	authStore.chargeStations["cs001"] = &store.ChargeStationAuth{
		SecurityProfile: 1,
	}

	csm := api.ChargeStationManager{
		ChargeStationAuthStore: authStore,
	}

	r := chi.NewRouter()
	r.Get("/{csId}/auth", csm.RetrieveChargeStationAuthDetails)
	server := httptest.NewServer(r)
	defer server.Close()

	req := httptest.NewRequest(http.MethodGet, "/cs001/auth", strings.NewReader("{}"))
	req.Header.Set("accept", "application/json")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Result().StatusCode)
	b, err := io.ReadAll(rr.Result().Body)
	require.NoError(t, err)

	want := &api.ChargeStationAuthDetailsResponse{
		SecurityProfile: 1,
	}

	got := new(api.ChargeStationAuthDetailsResponse)
	err = json.Unmarshal(b, &got)
	require.NoError(t, err)

	assert.Equal(t, want, got)
}

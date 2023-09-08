package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thoughtworks/maeve-csms/manager/ocpi"
	"github.com/thoughtworks/maeve-csms/manager/store"
	"github.com/thoughtworks/maeve-csms/manager/store/inmemory"
	"io"
	"k8s.io/utils/clock"
	clockTest "k8s.io/utils/clock/testing"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSwaggerHandler(t *testing.T) {
	engine := inmemory.NewStore(clock.RealClock{})
	ocpiApi := ocpi.NewOCPI(engine, http.DefaultClient, "GB", "TWK")
	handler := NewOcpiHandler(engine, clock.RealClock{}, ocpiApi)

	req := httptest.NewRequest(http.MethodGet, "/openapi.json", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	res := w.Result()
	defer func() {
		err := res.Body.Close()
		if err != nil {
			t.Errorf("closing body: %v", err)
		}
	}()

	b, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	if res.StatusCode != http.StatusOK {
		t.Errorf("status code: want %d, got %d", http.StatusOK, res.StatusCode)
	}

	var jsonData map[string]interface{}
	err = json.Unmarshal(b, &jsonData)
	require.NoError(t, err)
	require.Equal(t, jsonData["info"].(map[string]any)["title"], "Open Charge Point Interface (OCPI) 2.2")

}

func TestAPIRequestWithValidToken(t *testing.T) {
	token := "abcdef123456"
	engine := inmemory.NewStore(clock.RealClock{})
	err := engine.SetRegistrationDetails(context.Background(), token, &store.OcpiRegistration{Status: store.OcpiRegistrationStatusPending})
	require.NoError(t, err)
	now, err := time.Parse(time.RFC3339, "2023-06-15T15:05:00Z")
	require.NoError(t, err)
	clock := clockTest.NewFakePassiveClock(now)
	ocpiApi := ocpi.NewOCPI(engine, http.DefaultClient, "GB", "TWK")
	handler := NewOcpiHandler(engine, clock, ocpiApi)

	req := httptest.NewRequest(http.MethodGet, "/ocpi/versions", nil)
	req.Header.Add("Authorization", fmt.Sprintf("Token %s", token))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	res := w.Result()
	defer func() {
		err := res.Body.Close()
		if err != nil {
			t.Errorf("closing body: %v", err)
		}
	}()

	b, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	if res.StatusCode != http.StatusOK {
		t.Errorf("status code: want %d, got %d", http.StatusOK, res.StatusCode)
	}

	assert.JSONEq(t, `{"status_code":1000,"status_message":"Success","timestamp":"2023-06-15T15:05:00Z","data":[{"url":"/ocpi/2.2","version":"2.2"}]}`, string(b))
}

func TestAPIRequestWithInvalidToken(t *testing.T) {
	token := "abcdef123456"
	engine := inmemory.NewStore(clock.RealClock{})
	ocpiApi := ocpi.NewOCPI(engine, http.DefaultClient, "GB", "TWK")
	handler := NewOcpiHandler(engine, clock.RealClock{}, ocpiApi)

	req := httptest.NewRequest(http.MethodGet, "/ocpi/versions", nil)
	req.Header.Add("Authorization", fmt.Sprintf("Token %s", token))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	res := w.Result()
	defer func() {
		err := res.Body.Close()
		if err != nil {
			t.Errorf("closing body: %v", err)
		}
	}()

	if res.StatusCode != http.StatusUnauthorized {
		t.Errorf("status code: want %d, got %d", http.StatusOK, res.StatusCode)
	}
}

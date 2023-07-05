package server_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/twlabs/maeve-csms/manager/server"
)

func TestHealthHandler(t *testing.T) {
	handler := server.NewApiHandler(nil)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	res := w.Result()
	defer func() {
		err := res.Body.Close()
		if err != nil {
			t.Errorf("closing body: %v", err)
		}
	}()

	if res.StatusCode != http.StatusOK {
		t.Errorf("status code: want %d, got %d", http.StatusOK, res.StatusCode)
	}
}

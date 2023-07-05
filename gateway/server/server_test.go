package server_test

import (
	"context"
	"fmt"
	"github.com/twlabs/maeve-csms/gateway/server"
	"net/http"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	s := server.New("status server", "127.0.0.1:0", nil,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	errCh := make(chan error, 1)
	s.Start(errCh)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	defer func() {
		err := s.Stop(ctx)
		if err != nil {
			t.Errorf("stopping server: %v", err)
		}
	}()

	// give the server a moment to start
	time.Sleep(100 * time.Millisecond)

	select {
	case err := <-errCh:
		t.Fatalf("starting server: %v", err)
	default:
		// do nothing
	}

	if s.Addr() == "" {
		t.Fatal("expected server to have addr")
	}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s", s.Addr()), nil)
	if err != nil {
		t.Fatalf("creating http request: %v", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("making http request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("status code: want %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

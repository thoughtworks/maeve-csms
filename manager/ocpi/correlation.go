package ocpi

import (
	"github.com/google/uuid"
	"golang.org/x/net/context"
	"net/http"
)

type ContextKey string

const (
	ContextKeyCorrelationId ContextKey = "correlation_id"
)

func CorrelationIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id := r.Header.Get("X-Correlation-Id")
		if id == "" {
			id = uuid.New().String()
		}
		ctx = context.WithValue(ctx, ContextKeyCorrelationId, id)
		r = r.WithContext(ctx)
		w.Header().Set("X-Correlation-Id", id)
		next.ServeHTTP(w, r)
	})
}

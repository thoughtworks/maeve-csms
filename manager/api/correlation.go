package api

import (
	"github.com/google/uuid"
	"github.com/thoughtworks/maeve-csms/manager/ocpi"
	"golang.org/x/net/context"
	"net/http"
)

func CorrelationIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id := r.Header.Get("X-Correlation-Id")
		if id == "" {
			id = uuid.New().String()
		}
		ctx = context.WithValue(ctx, ocpi.ContextKeyCorrelationId, id)
		r = r.WithContext(ctx)
		w.Header().Set("X-Correlation-Id", id)
		next.ServeHTTP(w, r)
	})
}

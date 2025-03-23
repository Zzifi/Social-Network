package dbmiddleware

import (
	"context"
	"net/http"
	"user-service/internal/storage/postgre"
)

func DBMiddleware(storage *postgre.Storage) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), "db", storage)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

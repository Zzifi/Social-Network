package jwtauth

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

var (
	UserServiceURL = os.Getenv("USER_SERVICE_URL")
)

func JwtAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/user_service")

		if path == "/register" || path == "/login" {
			next.ServeHTTP(w, r)
			return
		}

		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Токен отсутствует", http.StatusUnauthorized)
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		validateURL := fmt.Sprintf("%s/session/validate?token=%s", UserServiceURL, tokenString)
		resp, err := http.Get(validateURL)
		resp.Body.Close()
		if err != nil || resp.StatusCode != http.StatusOK {
			http.Error(w, "Токен невалидный", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

package jwtauth_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	jwtauth "api-gateway/internal/http-server/middleware"
)

func mockUserService(t *testing.T, expectedToken string, statusCode int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "/session/validate") {
			t.Fatalf("Ожидался запрос на /session/validate, но пришел на %s", r.URL.Path)
		}

		token := r.URL.Query().Get("token")
		if token != expectedToken {
			http.Error(w, "Неверный токен", http.StatusUnauthorized)
			return
		}

		w.WriteHeader(statusCode)
	}))
}

func testJwtMiddleware(t *testing.T, token string, mockStatus, expectedCode int) {
	mockServer := mockUserService(t, "valid_token", mockStatus)
	defer mockServer.Close()

	originalUserServiceURL := jwtauth.UserServiceURL
	jwtauth.UserServiceURL = mockServer.URL
	defer func() { jwtauth.UserServiceURL = originalUserServiceURL }()

	middleware := jwtauth.JwtAuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/protected", nil)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	rec := httptest.NewRecorder()

	middleware.ServeHTTP(rec, req)

	if rec.Code != expectedCode {
		body, _ := io.ReadAll(rec.Body)
		t.Errorf("Ожидался статус %d, получен %d. Ответ: %s", expectedCode, rec.Code, string(body))
	}
}

func TestJwtAuthMiddleware_NoToken(t *testing.T) {
	testJwtMiddleware(t, "", http.StatusUnauthorized, http.StatusUnauthorized)
}

func TestJwtAuthMiddleware_InvalidToken(t *testing.T) {
	testJwtMiddleware(t, "invalid_token", http.StatusUnauthorized, http.StatusUnauthorized)
}

func TestJwtAuthMiddleware_ValidToken(t *testing.T) {
	testJwtMiddleware(t, "valid_token", http.StatusOK, http.StatusOK)
}

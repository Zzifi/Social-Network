package handlers_test

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"api-gateway/internal/http-server/handlers"

	"github.com/stretchr/testify/assert"
)

func TestReverseProxy(t *testing.T) {
	logger := slog.Default()

	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer backend.Close()

	pathWas := "/service/test"
	pathExpected := "/test"

	targetURL, _ := url.Parse(backend.URL)
	req := httptest.NewRequest(http.MethodGet, pathWas, nil)
	w := httptest.NewRecorder()

	handler := handlers.ReverseProxy(logger, targetURL.String())
	handler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, targetURL.Host, req.Host)
	assert.Equal(t, pathExpected, req.URL.Path)
}

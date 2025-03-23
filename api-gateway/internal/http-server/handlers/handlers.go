package handlers

import (
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func ReverseProxy(logger *slog.Logger, target string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		targetURL, err := url.Parse(target)
		if err != nil {
			http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
			logger.Error("Ошибка парсинга URL", slog.String("url", target), slog.String("error", err.Error()))
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(targetURL)
		parts := strings.SplitN(r.URL.Path, "/", 3)
		if len(parts) > 2 {
			r.URL.Path = "/" + parts[2]
		} else {
			r.URL.Path = "/"
		}
		r.Host = targetURL.Host

		logger.Info("Проксирование запроса",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("service", target),
		)

		proxy.ServeHTTP(w, r)
	}
}

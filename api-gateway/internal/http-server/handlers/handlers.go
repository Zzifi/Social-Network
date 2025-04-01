package handlers

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	postv1 "api-gateway/pkg/post-service-api"
)

func ReverseProxy(logger *slog.Logger, target string, postClient postv1.PostServiceClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		targetURL, err := url.Parse(target)
		if err != nil {
			http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
			logger.Error("Ошибка парсинга URL", slog.String("url", target), slog.String("error", err.Error()))
			return
		}

		parts := strings.SplitN(r.URL.Path, "/", 3)
		serviceName = parts[1]

		if serviceName == "post_service" {
			if len(parts) < 3 {
				http.Error(w, "Некорректный путь", http.StatusBadRequest)
				return
			}
			method := parts[2]

			var (
				reqBody []byte
				err     error
			)
			if r.Body != nil {
				reqBody, err = io.ReadAll(r.Body)
				if err != nil {
					http.Error(w, "Ошибка чтения тела запроса", http.StatusBadRequest)
					return
				}
			}

			ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
			defer cancel()

			var resp any
			switch method {
			case "create":
				var req postv1.CreatePostRequest
				if err := json.Unmarshal(reqBody, &req); err != nil {
					http.Error(w, "Ошибка разбора json", http.StatusBadRequest)
					return
				}
				resp, err = postClient.CreatePost(ctx, &req)
			case "delete":
				var req postv1.DeletePostRequest
				if err := json.Unmarshal(reqBody, &req); err != nil {
					http.Error(w, "Ошибка разбора json", http.StatusBadRequest)
					return
				}
				resp, err = postClient.DeletePost(ctx, &req)
			case "update":
				var req postv1.UpdatePostRequest
				if err := json.Unmarshal(reqBody, &req); err != nil {
					http.Error(w, "Ошибка разбора json", http.StatusBadRequest)
					return
				}
				resp, err = postClient.UpdatePost(ctx, &req)
			case "get":
				var req postv1.GetPostRequest
				if err := json.Unmarshal(reqBody, &req); err != nil {
					http.Error(w, "Ошибка разбора json", http.StatusBadRequest)
					return
				}
				resp, err = postClient.GetPost(ctx, &req)
			case "list":
				var req postv1.ListPostsRequest
				if err := json.Unmarshal(reqBody, &req); err != nil {
					http.Error(w, "Ошибка разбора json", http.StatusBadRequest)
					return
				}
				resp, err = postClient.ListPosts(ctx, &req)
			default:
				http.Error(w, "Неизвестный метод", http.StatusNotFound)
				return
			}

			if err != nil {
				logger.Error("Ошибка вызова grpc", slog.String("method", method), slog.String("error", err.Error()))
				http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(resp); err != nil {
				http.Error(w, "Ошибка кодирования json", http.StatusInternalServerError)
				return
			}
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(targetURL)
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

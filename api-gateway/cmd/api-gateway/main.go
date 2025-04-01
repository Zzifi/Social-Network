package main

import (
	"api-gateway/internal/config"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	httpHandlers "api-gateway/internal/http-server/handlers"
	mwJwtAuth "api-gateway/internal/http-server/middleware"

	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var routes = map[string]string{
	"/user_service": os.Getenv("USER_SERVICE_URL"),
	"/post_service": os.Getenv("POST_SERVICE_URL"),
}

var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

func main() {
	config := config.MustLoad()

	router := chi.NewRouter()

	conn, err := grpc.Dial(routes["/post_service"], grpc.WithInsecure())
	if err != nil:
		logger.Error("could not ", err)
	defer conn.Close()
	postClient = postv1.NewPostServiceClient(conn)

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(mwJwtAuth.JwtAuthMiddleware)

	for prefix, target := range routes {
		router.Handle(prefix+"/*", httpHandlers.ReverseProxy(logger, target, postClient))
	}

	logger.Info("starting server", slog.String("address", config.HTTPServer.Address))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	server := &http.Server{
		Addr:         config.HTTPServer.Address[strings.LastIndex(config.HTTPServer.Address, ":"):],
		Handler:      router,
		ReadTimeout:  config.HTTPServer.Timeout,
		WriteTimeout: config.HTTPServer.Timeout,
		IdleTimeout:  config.HTTPServer.IdleTimeout,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			logger.Error("failed to start server")
		}
	}()

	logger.Info("server started")

	<-done

	logger.Info("stopping server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("failed to stop server")
		return
	}

	logger.Info("server stopped")
}

package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
	"user-service/internal/config"

	"user-service/internal/storage/postgre"

	httpHandlers "user-service/internal/http-server/handlers"
	mwDB "user-service/internal/http-server/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
)

var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

func main() {
	storage, err := postgre.New(logger)
	if err != nil {
		panic(err)
	}
	defer storage.Close()

	config := config.MustLoad()

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(mwDB.DBMiddleware(storage))

	handlers := httpHandlers.NewHandlers(logger, storage)
	handlers.RegisterRoutes(router)

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

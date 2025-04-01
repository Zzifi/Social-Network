package main

import (
	"log/slog"
	"os"
	"os/signal"
	"post-service/internal/app"
	"syscall"
)

var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

func main() {
	cfg := config.MustLoad()

	application := app.New(logger, cfg.GRPC.Port, cfg.StoragePath)

	go func() {
		application.GRPCServer.MustRun()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.GRPCServer.Stop()
	logger.Info("server stopped")
}

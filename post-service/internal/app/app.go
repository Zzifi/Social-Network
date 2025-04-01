package app

import (
	"log/slog"

	grpcapp "post-service/internal/app/grpc"
	"post-service/internal/services/post"
	"post-service/storage/cassandra"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	storagePath string,
) *App {
	storage, err := cassandra.New(log)
	if err != nil {
		panic(err)
	}

	postService := post.New(log, storage)

	grpcApp := grpcapp.New(log, postService, grpcPort)

	return &App{
		GRPCServer: grpcApp,
	}
}

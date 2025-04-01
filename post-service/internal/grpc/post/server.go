package postgrpc

import (
	"context"

	postv1 "post-service/pkg/post-service-api"

	"google.golang.org/grpc"
)

type Post interface {
	CreatePost(
		ctx context.Context,
		id string,
		title string,
		description string,
		user_id string,
		is_private bool,
		tags []string,
	)
}

type serverAPI struct {
	postv1.UnimplementedPostServiceServer
	post Post
}

func Register(gRPCServer *grpc.Server, post Post) {
	postv1.RegisterPostServiceServer(gRPCServer, &serverAPI{post: post})
}

func (s *serverAPI) CreatePost(
	ctx context.Context,
	id string,
	title string,
	description string,
	user_id string,
	is_private bool,
	tags []string,
) (*postv1.LoginResponse, error) {

}

package post

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	postv1 "post-service/pkg/post-service-api"
)

type PostService struct {
	log     *slog.Logger
	storage Storage
}

var (
	ErrPostNotFound = errors.New("post not found")
)

type Storage interface {
	CreatePost(ctx context.Context, post postv1.Post) (string, error)
	UpdatePost(ctx context.Context, post postv1.Post) (postv1.Post, error)
	DeletePost(ctx context.Context, id string) error
	GetPost(ctx context.Context, id string) (postv1.Post, error)
	ListPosts(ctx context.Context, userID string, page, limit int) ([]postv1.Post, error)
}

func New(log *slog.Logger, storage Storage) *PostService {
	return &PostService{
		log:     log,
		storage: storage,
	}
}

func (s *PostService) CreatePost(ctx context.Context, req *postv1.CreatePostRequest) (*postv1.PostResponse, error) {
	s.log.Info("Creating new post", slog.String("user_id", req.UserId))

	post := postv1.Post{
		Title:       req.Title,
		Description: req.Description,
		UserID:      req.UserId,
		IsPrivate:   req.IsPrivate,
		Tags:        req.Tags,
	}

	id, err := s.storage.CreatePost(ctx, post)
	if err != nil {
		s.log.Error("Failed to create post", slog.String("error", err.Error()))
		return nil, fmt.Errorf("CreatePost: %w", err)
	}

	post.ID = id
	return &postv1.PostResponse{Post: convertToProto(post)}, nil
}

func (s *PostService) UpdatePost(ctx context.Context, req *postv1.UpdatePostRequest) (*postv1.PostResponse, error) {
	s.log.Info("Updating post", slog.String("post_id", req.Id))

	post := postv1.Post{
		ID:          req.Id,
		Title:       req.Title,
		Description: req.Description,
		IsPrivate:   req.IsPrivate,
		Tags:        req.Tags,
	}

	updatedPost, err := s.storage.UpdatePost(ctx, post)
	if err != nil {
		s.log.Error("Failed to update post", slog.String("error", err.Error()))
		return nil, fmt.Errorf("UpdatePost: %w", err)
	}

	return &postv1.PostResponse{Post: convertToProto(updatedPost)}, nil
}

func (s *PostService) DeletePost(ctx context.Context, req *postv1.DeletePostRequest) (*postv1.EmptyResponse, error) {
	s.log.Info("Deleting post", slog.String("post_id", req.Id))

	if err := s.storage.DeletePost(ctx, req.Id); err != nil {
		s.log.Error("Failed to delete post", slog.String("error", err.Error()))
		return nil, fmt.Errorf("DeletePost: %w", err)
	}

	return &postv1.EmptyResponse{}, nil
}

func (s *PostService) GetPost(ctx context.Context, req *postv1.GetPostRequest) (*postv1.PostResponse, error) {
	s.log.Info("Fetching post", slog.String("post_id", req.Id))

	post, err := s.storage.GetPost(ctx, req.Id)
	if err != nil {
		s.log.Error("Failed to fetch post", slog.String("error", err.Error()))
		return nil, fmt.Errorf("GetPost: %w", err)
	}

	return &postv1.PostResponse{Post: convertToProto(post)}, nil
}

func (s *PostService) ListPosts(ctx context.Context, req *postv1.ListPostsRequest) (*postv1.ListPostsResponse, error) {
	s.log.Info("Listing posts", slog.String("user_id", req.UserId))

	posts, err := s.storage.ListPosts(ctx, req.UserId, int(req.Page), int(req.Limit))
	if err != nil {
		s.log.Error("Failed to list posts", slog.String("error", err.Error()))
		return nil, fmt.Errorf("ListPosts: %w", err)
	}

	var protoPosts []*postv1.Post
	for _, post := range posts {
		protoPosts = append(protoPosts, convertToProto(post))
	}

	return &postv1.ListPostsResponse{Posts: protoPosts}, nil
}

func convertToProto(post postv1.Post) *postv1.Post {
	return &postv1.Post{
		Id:          post.ID,
		Title:       post.Title,
		Description: post.Description,
		UserId:      post.UserID,
		IsPrivate:   post.IsPrivate,
		Tags:        post.Tags,
	}
}

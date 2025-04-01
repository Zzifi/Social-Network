package cassandra

import (
	"context"
	"os"
	"time"

	postv1 "post-service/pkg/post-service-api"

	"github.com/gocql/gocql"
	"github.com/gookit/slog"
)

type Storage struct {
	session *gocql.Session
	logger  *slog.Logger
}

func New(logger *slog.Logger) (*Storage, error) {
	cassandraHost := os.Getenv("CASSANDRA_CONTACT_POINT")
	cassandraUsername := os.Getenv("CASSANDRA_USERNAME")
	cassandraPassword := os.Getenv("CASSANDRA_PASSWORD")

	cluster := gocql.NewCluster(cassandraHost)
	cluster.Consistency = gocql.Quorum
	cluster.ProtoVersion = 4
	cluster.ConnectTimeout = time.Second * 10
	cluster.Authenticator = gocql.PasswordAuthenticator{Username: cassandraUsername, Password: cassandraPassword}

	session, err := cluster.CreateSession()
	if err != nil {
		logger.Error("Ошибка подключения к Cassandra", err)
		return nil, err
	}

	return &Storage{session: session, logger: logger}, nil
}

func (s *Storage) CreatePost(ctx context.Context, req *postv1.CreatePostRequest) (*postv1.PostResponse, error) {
	id := gocql.TimeUUID().String()
	createdAt := time.Now().Format(time.RFC3339)
	updatedAt := createdAt

	err := s.session.Query(`INSERT INTO posts (id, title, description, user_id, created_at, updated_at, is_private, tags) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		id, req.Title, req.Description, req.UserId, createdAt, updatedAt, req.IsPrivate, req.Tags).
		WithContext(ctx).Exec()
	if err != nil {
		return nil, err
	}

	return &postv1.PostResponse{Post: &postv1.Post{
		Id:          id,
		Title:       req.Title,
		Description: req.Description,
		UserId:      req.UserId,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
		IsPrivate:   req.IsPrivate,
		Tags:        req.Tags,
	}}, nil
}

func (s *Storage) GetPost(ctx context.Context, req *postv1.GetPostRequest) (*postv1.PostResponse, error) {
	var post postv1.Post

	err := s.session.Query(`SELECT id, title, description, user_id, created_at, updated_at, is_private, tags FROM posts WHERE id = ? LIMIT 1`, req.Id).
		WithContext(ctx).Scan(
		&post.Id, &post.Title, &post.Description, &post.UserId,
		&post.CreatedAt, &post.UpdatedAt, &post.IsPrivate, &post.Tags)
	if err != nil {
		return nil, err
	}

	return &postv1.PostResponse{Post: &post}, nil
}

func (s *Storage) UpdatePost(ctx context.Context, req *postv1.UpdatePostRequest) (*postv1.PostResponse, error) {
	updatedAt := time.Now().Format(time.RFC3339)

	err := s.session.Query(`UPDATE posts SET title = ?, description = ?, is_private = ?, tags = ?, updated_at = ? WHERE id = ?`,
		req.Title, req.Description, req.IsPrivate, req.Tags, updatedAt, req.Id).
		WithContext(ctx).Exec()
	if err != nil {
		return nil, err
	}

	return s.GetPost(ctx, &postv1.GetPostRequest{Id: req.Id})
}

func (s *Storage) DeletePost(ctx context.Context, req *postv1.DeletePostRequest) (*postv1.EmptyResponse, error) {
	err := s.session.Query(`DELETE FROM posts WHERE id = ?`, req.Id).
		WithContext(ctx).Exec()
	if err != nil {
		return nil, err
	}
	return &postv1.EmptyResponse{}, nil
}

func (s *Storage) ListPosts(ctx context.Context, req *postv1.ListPostsRequest) (*postv1.ListPostsResponse, error) {
	var posts []*postv1.Post

	iter := s.session.Query(`SELECT id, title, description, user_id, created_at, updated_at, is_private, tags FROM posts WHERE user_id = ? LIMIT ?`, req.UserId, req.Limit).
		WithContext(ctx).Iter()

	var post postv1.Post
	for iter.Scan(&post.Id, &post.Title, &post.Description, &post.UserId, &post.CreatedAt, &post.UpdatedAt, &post.IsPrivate, &post.Tags) {
		posts = append(posts, &post)
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}

	return &postv1.ListPostsResponse{Posts: posts}, nil
}

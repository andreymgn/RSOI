package poststats

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net"

	pb "github.com/andreymgn/RSOI/services/poststats/proto"
	_ "github.com/lib/pq"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	ErrPostUidNotSet = errors.New("post UUID is required")
)

// Server implements poststats service
type Server struct {
	db *sql.DB
}

// NewServer returns a new server
func NewServer(connString string) (*Server, error) {
	db, err := sql.Open("postgres", connString)
	return &Server{db}, err
}

// Start starts a server
func (s *Server) Start(port int) error {
	server := grpc.NewServer()
	pb.RegisterPostStatsServer(server, s)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	return server.Serve(lis)
}

// GetPostStats returns post stats
func (s *Server) GetPostStats(ctx context.Context, req *pb.GetPostStatsRequest) (*pb.GetPostStatsResponse, error) {
	if req.PostUid == "" {
		return nil, ErrPostUidNotSet
	}

	query := "SELECT * FROM posts_stats WHERE post_uid=$1"
	row := s.db.QueryRow(query, req.PostUid)
	res := new(pb.GetPostStatsResponse)
	if err := row.Scan(&res.PostUid, &res.NumLikes, &res.NumDislikes, &res.NumViews); err != nil {
		return nil, err
	}

	return res, nil
}

// LikePost increases number of post likes
func (s *Server) LikePost(ctx context.Context, req *pb.LikePostRequest) (*pb.LikePostResponse, error) {
	if req.PostUid == "" {
		return nil, ErrPostUidNotSet
	}

	query := "UPDATE posts_stats SET num_likes = num_likes + 1 WHERE post_uid=$1"
	_, err := s.db.Exec(query, req.PostUid)
	if err != nil {
		return nil, err
	}

	res := new(pb.LikePostResponse)
	return res, nil
}

// DislikePost increases number of post dislikes
func (s *Server) DislikePost(ctx context.Context, req *pb.DislikePostRequest) (*pb.DislikePostResponse, error) {
	if req.PostUid == "" {
		return nil, ErrPostUidNotSet
	}

	query := "UPDATE posts_stats SET num_likes = num_likes - 1 WHERE post_uid=$1"
	_, err := s.db.Exec(query, req.PostUid)
	if err != nil {
		return nil, err
	}

	res := new(pb.DislikePostResponse)
	return res, nil
}

// IncreaseViews increases number of post views
func (s *Server) IncreaseViews(ctx context.Context, req *pb.IncreaseViewsRequest) (*pb.IncreaseViewsResponse, error) {
	if req.PostUid == "" {
		return nil, ErrPostUidNotSet
	}

	query := "UPDATE posts_stats SET num_views = num_views + 1 WHERE post_uid=$1"
	_, err := s.db.Exec(query, req.PostUid)
	if err != nil {
		return nil, err
	}

	res := new(pb.IncreaseViewsResponse)
	return res, nil
}

// DeletePostStats deletes stats of a post
func (s *Server) DeletePostStats(ctx context.Context, req *pb.DeletePostStatsRequest) (*pb.DeletePostStatsResponse, error) {
	if req.PostUid == "" {
		return nil, ErrPostUidNotSet
	}

	query := "DELETE FROM posts_stats WHERE post_uid=$1"
	_, err := s.db.Exec(query, req.PostUid)
	if err != nil {
		return nil, err
	}

	res := new(pb.DeletePostStatsResponse)
	return res, nil
}

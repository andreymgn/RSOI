package poststats

import (
	"errors"
	"fmt"
	"log"
	"net"

	pb "github.com/andreymgn/RSOI/services/poststats/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	ErrPostUidNotSet = errors.New("post UUID is required")
)

// Server implements poststats service
type Server struct {
	db datastore
}

// NewServer returns a new server
func NewServer(connString string) (*Server, error) {
	db, err := newDB(connString)
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

// GetPostStatsResponse converts PostStats to GetPostStatsResponse
func (ps *PostStats) GetPostStatsResponse() (*pb.GetPostStatsResponse, error) {
	res := new(pb.GetPostStatsResponse)
	res.PostUid = ps.Uid
	res.NumLikes = ps.NumLikes
	res.NumDislikes = ps.NumDislikes
	res.NumViews = ps.NumViews

	return res, nil
}

// GetPostStats returns post stats
func (s *Server) GetPostStats(ctx context.Context, req *pb.GetPostStatsRequest) (*pb.GetPostStatsResponse, error) {
	if req.PostUid == "" {
		return nil, ErrPostUidNotSet
	}

	postStats, err := s.db.get(req.PostUid)
	if err != nil {
		return nil, err
	}

	res, err := postStats.GetPostStatsResponse()
	if err != nil {
		return nil, err
	}

	return res, nil
}

// CreatePostStats creates a new post statistics record
func (s *Server) CreatePostStats(ctx context.Context, req *pb.CreatePostStatsRequest) (*pb.CreatePostStatsResponse, error) {
	if req.PostUid == "" {
		return nil, ErrPostUidNotSet
	}

	err := s.db.create(req.PostUid)
	if err != nil {
		return nil, err
	}

	res := new(pb.CreatePostStatsResponse)
	return res, nil
}

// LikePost increases number of post likes
func (s *Server) LikePost(ctx context.Context, req *pb.LikePostRequest) (*pb.LikePostResponse, error) {
	if req.PostUid == "" {
		return nil, ErrPostUidNotSet
	}

	err := s.db.like(req.PostUid)
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

	err := s.db.dislike(req.PostUid)
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

	err := s.db.view(req.PostUid)
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

	err := s.db.delete(req.PostUid)
	if err != nil {
		return nil, err
	}

	res := new(pb.DeletePostStatsResponse)
	return res, nil
}

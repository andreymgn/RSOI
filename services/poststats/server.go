package poststats

import (
	"fmt"
	"log"
	"net"

	pb "github.com/andreymgn/RSOI/services/poststats/proto"
	"github.com/google/uuid"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
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

// SinglePostStats converts PostStats to SinglePostStats
func (ps *PostStats) SinglePostStats() (*pb.SinglePostStats, error) {
	res := new(pb.SinglePostStats)
	res.PostUid = ps.UID.String()
	res.NumLikes = ps.NumLikes
	res.NumDislikes = ps.NumDislikes
	res.NumViews = ps.NumViews

	return res, nil
}

// GetPostStats returns post stats
func (s *Server) GetPostStats(ctx context.Context, req *pb.GetPostStatsRequest) (*pb.SinglePostStats, error) {
	uid, err := uuid.Parse(req.PostUid)
	if err != nil {
		return nil, err
	}

	postStats, err := s.db.get(uid)
	if err != nil {
		return nil, err
	}

	res, err := postStats.SinglePostStats()
	if err != nil {
		return nil, err
	}

	return res, nil
}

// CreatePostStats creates a new post statistics record
func (s *Server) CreatePostStats(ctx context.Context, req *pb.CreatePostStatsRequest) (*pb.SinglePostStats, error) {
	uid, err := uuid.Parse(req.PostUid)
	if err != nil {
		return nil, err
	}

	postStats, err := s.db.create(uid)
	if err != nil {
		return nil, err
	}

	return postStats.SinglePostStats()
}

// LikePost increases number of post likes
func (s *Server) LikePost(ctx context.Context, req *pb.LikePostRequest) (*pb.LikePostResponse, error) {
	uid, err := uuid.Parse(req.PostUid)
	if err != nil {
		return nil, err
	}

	err = s.db.like(uid)
	if err != nil {
		return nil, err
	}

	res := new(pb.LikePostResponse)
	return res, nil
}

// DislikePost increases number of post dislikes
func (s *Server) DislikePost(ctx context.Context, req *pb.DislikePostRequest) (*pb.DislikePostResponse, error) {
	uid, err := uuid.Parse(req.PostUid)
	if err != nil {
		return nil, err
	}

	err = s.db.dislike(uid)
	if err != nil {
		return nil, err
	}

	res := new(pb.DislikePostResponse)
	return res, nil
}

// IncreaseViews increases number of post views
func (s *Server) IncreaseViews(ctx context.Context, req *pb.IncreaseViewsRequest) (*pb.IncreaseViewsResponse, error) {
	uid, err := uuid.Parse(req.PostUid)
	if err != nil {
		return nil, err
	}

	err = s.db.view(uid)
	if err != nil {
		return nil, err
	}

	res := new(pb.IncreaseViewsResponse)
	return res, nil
}

// DeletePostStats deletes stats of a post
func (s *Server) DeletePostStats(ctx context.Context, req *pb.DeletePostStatsRequest) (*pb.DeletePostStatsResponse, error) {
	uid, err := uuid.Parse(req.PostUid)
	if err != nil {
		return nil, err
	}

	err = s.db.delete(uid)
	if err != nil {
		return nil, err
	}

	res := new(pb.DeletePostStatsResponse)
	return res, nil
}

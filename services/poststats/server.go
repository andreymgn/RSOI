package poststats

import (
	"fmt"
	"log"
	"net"

	pb "github.com/andreymgn/RSOI/services/poststats/proto"
	"github.com/google/uuid"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	opentracing "github.com/opentracing/opentracing-go"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	statusNotFound    = status.Error(codes.NotFound, "post not found")
	statusInvalidUUID = status.Error(codes.InvalidArgument, "invalid UUID")
)

func internalError(err error) error {
	return status.Error(codes.Internal, err.Error())
}

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
func (s *Server) Start(port int, tracer opentracing.Tracer) error {
	server := grpc.NewServer(grpc.UnaryInterceptor(otgrpc.OpenTracingServerInterceptor(tracer)))
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
		return nil, statusInvalidUUID
	}

	postStats, err := s.db.get(uid)
	switch err {
	case nil:
		return postStats.SinglePostStats()
	case errNotFound:
		return nil, statusNotFound
	default:
		return nil, internalError(err)
	}
}

// CreatePostStats creates a new post statistics record
func (s *Server) CreatePostStats(ctx context.Context, req *pb.CreatePostStatsRequest) (*pb.SinglePostStats, error) {
	uid, err := uuid.Parse(req.PostUid)
	if err != nil {
		return nil, statusInvalidUUID
	}

	postStats, err := s.db.create(uid)
	if err != nil {
		return nil, internalError(err)
	}

	return postStats.SinglePostStats()
}

// LikePost increases number of post likes
func (s *Server) LikePost(ctx context.Context, req *pb.LikePostRequest) (*pb.LikePostResponse, error) {
	uid, err := uuid.Parse(req.PostUid)
	if err != nil {
		return nil, statusInvalidUUID
	}

	err = s.db.like(uid)
	switch err {
	case nil:
		return new(pb.LikePostResponse), nil
	case errNotFound:
		return nil, statusNotFound
	default:
		return nil, internalError(err)
	}
}

// DislikePost increases number of post dislikes
func (s *Server) DislikePost(ctx context.Context, req *pb.DislikePostRequest) (*pb.DislikePostResponse, error) {
	uid, err := uuid.Parse(req.PostUid)
	if err != nil {
		return nil, statusInvalidUUID
	}

	err = s.db.dislike(uid)
	switch err {
	case nil:
		return new(pb.DislikePostResponse), nil
	case errNotFound:
		return nil, statusNotFound
	default:
		return nil, internalError(err)
	}
}

// IncreaseViews increases number of post views
func (s *Server) IncreaseViews(ctx context.Context, req *pb.IncreaseViewsRequest) (*pb.IncreaseViewsResponse, error) {
	uid, err := uuid.Parse(req.PostUid)
	if err != nil {
		return nil, statusInvalidUUID
	}

	err = s.db.view(uid)
	switch err {
	case nil:
		return new(pb.IncreaseViewsResponse), nil
	case errNotFound:
		return nil, statusNotFound
	default:
		return nil, internalError(err)
	}
}

// DeletePostStats deletes stats of a post
func (s *Server) DeletePostStats(ctx context.Context, req *pb.DeletePostStatsRequest) (*pb.DeletePostStatsResponse, error) {
	uid, err := uuid.Parse(req.PostUid)
	if err != nil {
		return nil, statusInvalidUUID
	}

	err = s.db.delete(uid)
	switch err {
	case nil:
		return new(pb.DeletePostStatsResponse), nil
	case errNotFound:
		return nil, statusNotFound
	default:
		return nil, internalError(err)
	}
}

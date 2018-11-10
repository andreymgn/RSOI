package poststats

import (
	"github.com/andreymgn/RSOI/services/auth"
	pb "github.com/andreymgn/RSOI/services/poststats/proto"
	"github.com/google/uuid"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	statusNotFound     = status.Error(codes.NotFound, "post not found")
	statusInvalidUUID  = status.Error(codes.InvalidArgument, "invalid UUID")
	statusInvalidToken = status.Errorf(codes.Unauthenticated, "invalid token")
)

func internalError(err error) error {
	return status.Error(codes.Internal, err.Error())
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
	valid, err := s.checkToken(req.Token)
	if err != nil {
		return nil, err
	}

	if !valid {
		return nil, statusInvalidToken
	}

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
	valid, err := s.checkToken(req.Token)
	if err != nil {
		return nil, err
	}

	if !valid {
		return nil, statusInvalidToken
	}

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
	valid, err := s.checkToken(req.Token)
	if err != nil {
		return nil, err
	}

	if !valid {
		return nil, statusInvalidToken
	}

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
	valid, err := s.checkToken(req.Token)
	if err != nil {
		return nil, err
	}

	if !valid {
		return nil, statusInvalidToken
	}

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
	valid, err := s.checkToken(req.Token)
	if err != nil {
		return nil, err
	}

	if !valid {
		return nil, statusInvalidToken
	}

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

// GetToken returns new authorization token
func (s *Server) GetToken(ctx context.Context, req *pb.GetTokenRequest) (*pb.GetTokenResponse, error) {
	appID, appSecret := req.AppId, req.AppSecret
	token, err := s.auth.Add(appID, appSecret)
	switch err {
	case nil:
		res := new(pb.GetTokenResponse)
		res.Token = token
		return res, nil
	case auth.ErrNotFound:
		return nil, statusNotFound
	case auth.ErrWrongSecret:
		return nil, status.Error(codes.Unauthenticated, "wrong secret")
	default:
		return nil, internalError(err)
	}
}

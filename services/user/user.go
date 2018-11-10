package user

import (
	"context"
	"time"

	"github.com/golang/protobuf/ptypes"

	"github.com/andreymgn/RSOI/services/auth"
	pb "github.com/andreymgn/RSOI/services/user/proto"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const DefaultExpirationTime = time.Minute * 15

var (
	statusNoUsername       = status.Error(codes.InvalidArgument, "username is empty")
	statusNotFound         = status.Error(codes.NotFound, "user not found")
	statusInvalidUUID      = status.Error(codes.InvalidArgument, "invalid UUID")
	statusInvalidToken     = status.Error(codes.Unauthenticated, "invalid grpc token")
	statusInvalidUserToken = status.Error(codes.Unauthenticated, "invalid user token")
)

func internalError(err error) error {
	return status.Error(codes.Internal, err.Error())
}

func (u *User) UserInfo() *pb.UserInfo {
	result := new(pb.UserInfo)
	result.Uid = u.UID.String()
	result.Username = u.Username
	return result
}

func (s *Server) GetUserInfo(ctx context.Context, req *pb.GetUserInfoRequest) (*pb.UserInfo, error) {
	uid, err := uuid.Parse(req.Uid)
	if err != nil {
		return nil, statusInvalidUUID
	}

	user, err := s.db.getUserInfo(uid)
	switch err {
	case nil:
		return user.UserInfo(), nil
	case errNotFound:
		return nil, statusNotFound
	default:
		return nil, internalError(err)
	}
}

func (s *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserInfo, error) {
	valid, err := s.checkServiceToken(req.Token)
	if err != nil {
		return nil, err
	}

	if !valid {
		return nil, statusInvalidToken
	}

	user, err := s.db.create(req.Username, req.Password)
	if err != nil {
		return nil, internalError(err)
	}

	return user.UserInfo(), nil
}

func (s *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	valid, err := s.checkServiceToken(req.ApiToken)
	if err != nil {
		return nil, err
	}

	if !valid {
		return nil, statusInvalidToken
	}

	uid, err := uuid.Parse(req.Uid)
	if err != nil {
		return nil, statusInvalidUUID
	}

	err = s.db.update(uid, req.Password)
	switch err {
	case nil:
		return new(pb.UpdateUserResponse), nil
	case errNotFound:
		return nil, statusNotFound
	default:
		return nil, internalError(err)
	}
}

func (s *Server) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	valid, err := s.checkServiceToken(req.ApiToken)
	if err != nil {
		return nil, err
	}

	if !valid {
		return nil, statusInvalidToken
	}

	uid, err := uuid.Parse(req.Uid)
	if err != nil {
		return nil, statusInvalidUUID
	}

	err = s.db.delete(uid)
	switch err {
	case nil:
		return new(pb.DeleteUserResponse), nil
	case errNotFound:
		return nil, statusNotFound
	default:
		return nil, internalError(err)
	}
}

func (s *Server) GetUserPosts(ctx context.Context, req *pb.GetUserPostsRequest) (*pb.GetUserPostsResponse, error) {
	var pageSize int32
	if req.PageSize == 0 {
		pageSize = 10
	} else {
		pageSize = req.PageSize
	}

	uid, err := uuid.Parse(req.Uid)
	if err != nil {
		return nil, statusInvalidUUID
	}

	uids, err := s.db.getPosts(uid, pageSize, req.PageNumber)
	res := new(pb.GetUserPostsResponse)
	res.PageSize = pageSize
	res.PageNumber = req.PageNumber
	for _, u := range uids {
		res.PostUid = append(res.PostUid, u.String())
	}

	return res, nil
}

func (s *Server) AddPost(ctx context.Context, req *pb.AddPostRequest) (*pb.AddPostResponse, error) {
	valid, err := s.checkServiceToken(req.ApiToken)
	if err != nil {
		return nil, err
	}

	if !valid {
		return nil, statusInvalidToken
	}

	uid, err := uuid.Parse(req.Uid)
	if err != nil {
		return nil, statusInvalidUUID
	}

	postUID, err := uuid.Parse(req.PostUid)
	if err != nil {
		return nil, statusInvalidUUID
	}

	createdAt, err := ptypes.Timestamp(req.CreatedAt)
	if err != nil {
		return nil, internalError(err)
	}

	err = s.db.addPost(uid, postUID, createdAt)
	if err != nil {
		return nil, internalError(err)
	}

	return new(pb.AddPostResponse), nil
}

func (s *Server) DeletePost(ctx context.Context, req *pb.DeletePostRequest) (*pb.DeletePostResponse, error) {
	valid, err := s.checkServiceToken(req.ApiToken)
	if err != nil {
		return nil, err
	}

	if !valid {
		return nil, statusInvalidToken
	}

	uid, err := uuid.Parse(req.Uid)
	if err != nil {
		return nil, statusInvalidUUID
	}

	postUID, err := uuid.Parse(req.PostUid)
	if err != nil {
		return nil, statusInvalidUUID
	}

	err = s.db.deletePost(uid, postUID)
	if err != nil {
		return nil, internalError(err)
	}

	return new(pb.DeletePostResponse), nil
}

func (s *Server) GetServiceToken(ctx context.Context, req *pb.GetServiceTokenRequest) (*pb.GetServiceTokenResponse, error) {
	appID, appSecret := req.AppId, req.AppSecret
	token, err := s.apiTokenAuth.Add(appID, appSecret)
	switch err {
	case nil:
		res := new(pb.GetServiceTokenResponse)
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

func (s *Server) GetUserToken(ctx context.Context, req *pb.GetUserTokenRequest) (*pb.GetUserTokenResponse, error) {
	uid, err := uuid.Parse(req.Uid)
	if err != nil {
		return nil, statusInvalidUUID
	}

	samePassword, err := s.db.checkPassword(uid, req.Password)
	if !samePassword {
		return nil, status.Error(codes.Unauthenticated, "wrong password")
	}

	token := uuid.New().String()
	err = s.userTokenAuth.Set(token, uid, DefaultExpirationTime).Err()
	if err != nil {
		return nil, internalError(err)
	}

	res := new(pb.GetUserTokenResponse)
	res.Token = token
	return res, nil
}

func (s *Server) GetUserByToken(ctx context.Context, req *pb.GetUserByTokenRequest) (*pb.GetUserByTokenResponse, error) {
	token := req.UserToken
	uid, err := s.userTokenAuth.Get(token).Result()
	if err != nil {
		return nil, internalError(err)
	}

	if _, err := uuid.Parse(uid); err != nil {
		return nil, statusInvalidUserToken
	}

	res := new(pb.GetUserByTokenResponse)
	res.Uid = uid
	return res, nil
}

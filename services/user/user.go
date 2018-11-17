package user

import (
	"context"
	"time"

	"github.com/go-redis/redis"

	"github.com/andreymgn/RSOI/services/auth"
	pb "github.com/andreymgn/RSOI/services/user/proto"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	AccessTokenExpirationTime  = time.Minute * 15
	RefreshTokenExpirationTime = time.Hour * 24 * 7 * 2
)

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

// UserInfo converts User to protobuf struct
func (u *User) UserInfo() *pb.UserInfo {
	result := new(pb.UserInfo)
	result.Uid = u.UID.String()
	result.Username = u.Username
	return result
}

// GetUserInfo returns User
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

// CreateUser creates a new user
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

// UpdateUser updates user
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

// DeleteUser deletes user
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

// GetServiceToken returns token for user service access
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

// GetAccessToken returns authorization token for user
func (s *Server) GetAccessToken(ctx context.Context, req *pb.GetTokenRequest) (*pb.GetAccessTokenResponse, error) {
	valid, err := s.checkServiceToken(req.ApiToken)
	if err != nil {
		return nil, err
	}

	if !valid {
		return nil, statusInvalidToken
	}

	uid, err := s.db.getUIDByUsername(req.Username)
	if err == errNotFound {
		return nil, statusNotFound
	} else if err != nil {
		return nil, internalError(err)
	}

	samePassword, err := s.db.checkPassword(uid, req.Password)
	if !samePassword {
		return nil, status.Error(codes.Unauthenticated, "wrong password")
	}

	token := uuid.New().String()
	err = s.accessTokenStorage.Set(token, uid.String(), AccessTokenExpirationTime).Err()
	if err != nil {
		return nil, internalError(err)
	}

	res := new(pb.GetAccessTokenResponse)
	res.Token = token
	res.Uid = uid.String()
	return res, nil
}

// GetUserByAccessToken checks access token existance and refreshes token expiration time
func (s *Server) GetUserByAccessToken(ctx context.Context, req *pb.GetUserByAccessTokenRequest) (*pb.GetUserByAccessTokenResponse, error) {
	valid, err := s.checkServiceToken(req.ApiToken)
	if err != nil {
		return nil, err
	}

	if !valid {
		return nil, statusInvalidToken
	}

	token := req.UserToken
	uid, err := s.accessTokenStorage.Get(token).Result()
	if err == redis.Nil {
		return nil, statusInvalidUserToken
	} else if err != nil {
		return nil, internalError(err)
	}

	if _, err := uuid.Parse(uid); err != nil {
		return nil, statusInvalidUserToken
	}

	err = s.accessTokenStorage.Expire(token, AccessTokenExpirationTime).Err()
	if err != nil {
		return nil, internalError(err)
	}

	res := new(pb.GetUserByAccessTokenResponse)
	res.Uid = uid
	return res, nil
}

// GetRefreshToken returns token which can be used to refresh access token
func (s *Server) GetRefreshToken(ctx context.Context, req *pb.GetTokenRequest) (*pb.GetRefreshTokenResponse, error) {
	valid, err := s.checkServiceToken(req.ApiToken)
	if err != nil {
		return nil, err
	}

	if !valid {
		return nil, statusInvalidToken
	}

	uid, err := s.db.getUIDByUsername(req.Username)
	if err == errNotFound {
		return nil, statusNotFound
	} else if err != nil {
		return nil, internalError(err)
	}

	samePassword, err := s.db.checkPassword(uid, req.Password)
	if !samePassword {
		return nil, status.Error(codes.Unauthenticated, "wrong password")
	}

	token := uuid.New().String()
	err = s.refreshTokenStorage.Set(token, uid.String(), RefreshTokenExpirationTime).Err()
	if err != nil {
		return nil, internalError(err)
	}

	res := new(pb.GetRefreshTokenResponse)
	res.Token = token
	return res, nil
}

func (s *Server) RefreshAccessToken(ctx context.Context, req *pb.RefreshAccessTokenRequest) (*pb.RefreshAccessTokenResponse, error) {
	valid, err := s.checkServiceToken(req.ApiToken)
	if err != nil {
		return nil, err
	}

	if !valid {
		return nil, statusInvalidToken
	}

	token := req.RefreshToken
	uid, err := s.refreshTokenStorage.Get(token).Result()
	if err == redis.Nil {
		return nil, statusInvalidUserToken
	} else if err != nil {
		return nil, internalError(err)
	}

	if _, err := uuid.Parse(uid); err != nil {
		return nil, statusInvalidUserToken
	}

	err = s.refreshTokenStorage.Del(token).Err()
	if err != nil {
		return nil, internalError(err)
	}

	refreshToken := uuid.New().String()
	err = s.refreshTokenStorage.Set(refreshToken, uid, RefreshTokenExpirationTime).Err()
	if err != nil {
		return nil, internalError(err)
	}

	accessToken := uuid.New().String()
	err = s.accessTokenStorage.Set(accessToken, uid, AccessTokenExpirationTime).Err()
	if err != nil {
		return nil, internalError(err)
	}

	res := new(pb.RefreshAccessTokenResponse)
	res.RefreshToken = refreshToken
	res.AccessToken = accessToken
	return res, nil
}

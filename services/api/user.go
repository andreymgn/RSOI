package api

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	user "github.com/andreymgn/RSOI/services/user/proto"
	"github.com/gorilla/mux"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) getUserInfo() http.HandlerFunc {
	type response struct {
		ID       string
		Username string
	}

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		uid := vars["uid"]

		ctx := r.Context()
		getUserResponse, err := s.userClient.client.GetUserInfo(ctx,
			&user.GetUserInfoRequest{Uid: uid},
		)
		if err != nil {
			handleRPCError(w, err)
			return
		}

		resp := response{getUserResponse.Uid, getUserResponse.Username}
		json, err := json.Marshal(resp)
		if err != nil {
			handleRPCError(w, err)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(json)
	}
}

func (s *Server) createUser() http.HandlerFunc {
	type request struct {
		Username string
		Password string
	}

	type response struct {
		ID       string
		Username string
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			handleRPCError(w, err)
			return
		}

		err = json.Unmarshal(b, &req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}

		ctx := r.Context()
		createUserResponse, err := s.userClient.client.CreateUser(ctx,
			&user.CreateUserRequest{Token: s.userClient.token, Username: req.Username, Password: req.Password},
		)
		if err != nil {
			if st, ok := status.FromError(err); ok && st.Code() == codes.Unauthenticated {
				err := s.updateUserToken()
				if err != nil {
					handleRPCError(w, err)
					return
				}
				createUserResponse, err = s.userClient.client.CreateUser(ctx,
					&user.CreateUserRequest{Token: s.userClient.token, Username: req.Username, Password: req.Password},
				)
				if err != nil {
					handleRPCError(w, err)
					return
				}
			} else {
				handleRPCError(w, err)
				return
			}
		}

		resp := response{createUserResponse.Uid, createUserResponse.Username}
		json, err := json.Marshal(resp)
		if err != nil {
			handleRPCError(w, err)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write(json)
	}
}

func (s *Server) getToken() http.HandlerFunc {
	type request struct {
		Username string
		Password string
		Refresh  bool `json:",omitempty"`
	}

	type response struct {
		AccessToken  string
		RefreshToken string `json:",omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			handleRPCError(w, err)
			return
		}

		err = json.Unmarshal(b, &req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}

		ctx := r.Context()
		accessTokenResponse, err := s.userClient.client.GetAccessToken(ctx,
			&user.GetTokenRequest{ApiToken: s.userClient.token, Username: req.Username, Password: req.Password},
		)
		if err != nil {
			if st, ok := status.FromError(err); ok && st.Code() == codes.Unauthenticated {
				err := s.updateUserToken()
				if err != nil {
					handleRPCError(w, err)
					return
				}
				accessTokenResponse, err = s.userClient.client.GetAccessToken(ctx,
					&user.GetTokenRequest{ApiToken: s.userClient.token, Username: req.Username, Password: req.Password},
				)
				if err != nil {
					handleRPCError(w, err)
					return
				}
			} else {
				handleRPCError(w, err)
				return
			}
		}

		resp := response{}
		resp.AccessToken = accessTokenResponse.Token

		if req.Refresh {
			refreshTokenResponse, err := s.userClient.client.GetRefreshToken(ctx,
				&user.GetTokenRequest{ApiToken: s.userClient.token, Username: req.Username, Password: req.Password},
			)
			if err != nil {
				if st, ok := status.FromError(err); ok && st.Code() == codes.Unauthenticated {
					err := s.updateUserToken()
					if err != nil {
						handleRPCError(w, err)
						return
					}
					refreshTokenResponse, err = s.userClient.client.GetRefreshToken(ctx,
						&user.GetTokenRequest{ApiToken: s.userClient.token, Username: req.Username, Password: req.Password},
					)
					if err != nil {
						handleRPCError(w, err)
						return
					}
				} else {
					handleRPCError(w, err)
					return
				}
			}

			resp.RefreshToken = refreshTokenResponse.Token
		}

		json, err := json.Marshal(resp)
		if err != nil {
			handleRPCError(w, err)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(json)
	}
}

func (s *Server) refreshToken() http.HandlerFunc {
	type request struct {
		Token string
	}

	type response struct {
		AccessToken  string
		RefreshToken string
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			handleRPCError(w, err)
			return
		}

		err = json.Unmarshal(b, &req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}

		ctx := r.Context()
		refreshTokenResponse, err := s.userClient.client.RefreshAccessToken(ctx,
			&user.RefreshAccessTokenRequest{ApiToken: s.userClient.token, RefreshToken: req.Token},
		)
		if err != nil {
			if st, ok := status.FromError(err); ok && st.Code() == codes.Unauthenticated {
				err := s.updateUserToken()
				if err != nil {
					handleRPCError(w, err)
					return
				}
				refreshTokenResponse, err = s.userClient.client.RefreshAccessToken(ctx,
					&user.RefreshAccessTokenRequest{ApiToken: s.userClient.token, RefreshToken: req.Token},
				)
				if err != nil {
					handleRPCError(w, err)
					return
				}
			} else {
				handleRPCError(w, err)
				return
			}
		}

		resp := response{refreshTokenResponse.AccessToken, refreshTokenResponse.RefreshToken}
		json, err := json.Marshal(resp)
		if err != nil {
			handleRPCError(w, err)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(json)
	}
}

func (s *Server) getUserByToken(token string) (string, error) {
	ctx := context.Background()
	uid, err := s.userClient.client.GetUserByAccessToken(ctx,
		&user.GetUserByAccessTokenRequest{ApiToken: s.userClient.token, UserToken: token},
	)
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.Unauthenticated {
			err := s.updateUserToken()
			if err != nil {
				return "", err
			}
			uid, err = s.userClient.client.GetUserByAccessToken(ctx,
				&user.GetUserByAccessTokenRequest{ApiToken: s.userClient.token, UserToken: token},
			)
			if err != nil {
				return "", err
			}
		} else {
			return "", err
		}
	}

	return uid.Uid, nil
}

func (s *Server) updateUserToken() error {
	token, err := s.userClient.client.GetServiceToken(context.Background(),
		&user.GetServiceTokenRequest{AppId: s.userClient.appID, AppSecret: s.userClient.appSecret},
	)
	if err != nil {
		return err
	}

	s.userClient.token = token.Token
	return nil
}

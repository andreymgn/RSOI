package api

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	user "github.com/andreymgn/RSOI/services/user/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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

func (s *Server) getUserToken() http.HandlerFunc {
	type request struct {
		Username string
		Password string
	}

	type response struct {
		Token string
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
		tokenResponse, err := s.userClient.client.GetUserToken(ctx,
			&user.GetUserTokenRequest{ApiToken: s.userClient.token, Username: req.Username, Password: req.Password},
		)
		if err != nil {
			if st, ok := status.FromError(err); ok && st.Code() == codes.Unauthenticated {
				err := s.updateUserToken()
				if err != nil {
					handleRPCError(w, err)
					return
				}
				tokenResponse, err = s.userClient.client.GetUserToken(ctx,
					&user.GetUserTokenRequest{ApiToken: s.userClient.token, Username: req.Username, Password: req.Password},
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

		resp := response{tokenResponse.Token}
		json, err := json.Marshal(resp)
		if err != nil {
			handleRPCError(w, err)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(json)
	}
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

func (s *Server) getUserByToken(token string) (string, error) {
	ctx := context.Background()
	uid, err := s.userClient.client.GetUserByToken(ctx,
		&user.GetUserByTokenRequest{ApiToken: s.userClient.token, UserToken: token},
	)
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.Unauthenticated {
			err := s.updateUserToken()
			if err != nil {
				return "", err
			}
			uid, err = s.userClient.client.GetUserByToken(ctx,
				&user.GetUserByTokenRequest{ApiToken: s.userClient.token, UserToken: token},
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

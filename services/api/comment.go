package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	comment "github.com/andreymgn/RSOI/services/comment/proto"
	post "github.com/andreymgn/RSOI/services/post/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/gorilla/mux"
)

func (s *Server) getPostComments() http.HandlerFunc {
	type c struct {
		UID        string
		PostUID    string
		Body       string
		ParentUID  string
		CreatedAt  time.Time
		ModifiedAt time.Time
	}

	type response struct {
		Comments   []c
		PageSize   int32
		PageNumber int32
	}

	return func(w http.ResponseWriter, r *http.Request) {
		page, size := r.URL.Query().Get("page"), r.URL.Query().Get("size")
		var pageNum, sizeNum int32 = 0, 10
		if page != "" {
			n, err := strconv.Atoi(page)
			if err != nil {
				http.Error(w, "can't parse query parameter `page`", http.StatusBadRequest)
				return
			}
			pageNum = int32(n)
		}

		if size != "" {
			n, err := strconv.Atoi(size)
			if err != nil {
				http.Error(w, "can't parse query parameter `size`", http.StatusBadRequest)
				return
			}
			sizeNum = int32(n)
		}

		vars := mux.Vars(r)
		uid := vars["uid"]
		postUID := vars["postuid"]

		ctx := r.Context()
		checkExistsResponse, err := s.postClient.client.CheckExists(ctx, &post.CheckExistsRequest{Uid: postUID})
		if err != nil {
			handleRPCError(w, err)
			return
		}

		if !checkExistsResponse.Exists {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		commentsResponse, err := s.commentClient.ListComments(ctx, &comment.ListCommentsRequest{PostUid: postUID, CommentUid: uid, PageSize: sizeNum, PageNumber: pageNum})
		if err != nil {
			handleRPCError(w, err)
			return
		}

		comments := make([]c, len(commentsResponse.Comments))
		for i, singleComment := range commentsResponse.Comments {
			comments[i].UID = singleComment.Uid
			comments[i].PostUID = singleComment.PostUid
			comments[i].Body = singleComment.Body
			comments[i].ParentUID = singleComment.ParentUid
			comments[i].CreatedAt, err = ptypes.Timestamp(singleComment.CreatedAt)
			if err != nil {
				handleRPCError(w, err)
				return
			}

			comments[i].ModifiedAt, err = ptypes.Timestamp(singleComment.ModifiedAt)
			if err != nil {
				handleRPCError(w, err)
				return
			}
		}

		resp := response{comments, sizeNum, pageNum}

		json, err := json.Marshal(resp)
		if err != nil {
			handleRPCError(w, err)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(json)
	}
}

func (s *Server) createComment() http.HandlerFunc {
	type request struct {
		Body      string `json:"body"`
		ParentUID string `json:"parent_uid"`
	}

	type response struct {
		UID        string
		PostUID    string
		Body       string
		ParentUID  string
		CreatedAt  time.Time
		ModifiedAt time.Time
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

		vars := mux.Vars(r)
		postUID := vars["postuid"]

		ctx := r.Context()
		checkExistsResponse, err := s.postClient.client.CheckExists(ctx, &post.CheckExistsRequest{Uid: postUID})
		if err != nil {
			handleRPCError(w, err)
			return
		}

		if !checkExistsResponse.Exists {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		c, err := s.commentClient.CreateComment(ctx, &comment.CreateCommentRequest{PostUid: postUID, Body: req.Body, ParentUid: req.ParentUID})
		if err != nil {
			handleRPCError(w, err)
			return
		}

		createdAt, err := ptypes.Timestamp(c.CreatedAt)
		if err != nil {
			handleRPCError(w, err)
			return
		}

		modifiedAt, err := ptypes.Timestamp(c.ModifiedAt)
		if err != nil {
			handleRPCError(w, err)
			return
		}

		response := response{c.Uid, c.PostUid, c.Body, c.ParentUid, createdAt, modifiedAt}
		json, err := json.Marshal(response)
		if err != nil {
			handleRPCError(w, err)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write(json)
	}
}

func (s *Server) updateComment() http.HandlerFunc {
	type request struct {
		Body string `json:"body"`
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

		vars := mux.Vars(r)
		uid := vars["uid"]
		postUID := vars["postuid"]

		ctx := r.Context()
		checkExistsResponse, err := s.postClient.client.CheckExists(ctx, &post.CheckExistsRequest{Uid: postUID})
		if err != nil {
			handleRPCError(w, err)
			return
		}

		if !checkExistsResponse.Exists {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		_, err = s.commentClient.UpdateComment(ctx, &comment.UpdateCommentRequest{Uid: uid, Body: req.Body})
		if err != nil {
			handleRPCError(w, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func (s *Server) deleteComment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		uid := vars["uid"]
		postUID := vars["postuid"]

		ctx := r.Context()
		checkExistsResponse, err := s.postClient.client.CheckExists(ctx, &post.CheckExistsRequest{Uid: postUID})
		if err != nil {
			handleRPCError(w, err)
			return
		}

		if !checkExistsResponse.Exists {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		_, err = s.commentClient.DeleteComment(ctx, &comment.DeleteCommentRequest{Uid: uid})
		if err != nil {
			handleRPCError(w, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

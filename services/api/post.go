package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	comment "github.com/andreymgn/RSOI/services/comment/proto"
	post "github.com/andreymgn/RSOI/services/post/proto"
	poststats "github.com/andreymgn/RSOI/services/poststats/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/gorilla/mux"
)

func (s *Server) getPosts() http.HandlerFunc {
	type p struct {
		UID         string
		Title       string
		URL         string
		CreatedAt   time.Time
		ModifiedAt  time.Time
		NumLikes    int32
		NumDislikes int32
		NumViews    int32
	}

	type response struct {
		Posts []p
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

		ctx := r.Context()
		postResponse, err := s.postClient.ListPosts(ctx, &post.ListPostsRequest{PageSize: sizeNum, PageNumber: pageNum})
		if err != nil {
			handleRPCError(w, err)
			return
		}

		posts := make([]p, len(postResponse.Posts))
		for i, singlePostResponse := range postResponse.Posts {
			posts[i].UID = singlePostResponse.Uid
			posts[i].Title = singlePostResponse.Title
			posts[i].URL = singlePostResponse.Url
			posts[i].CreatedAt, err = ptypes.Timestamp(singlePostResponse.CreatedAt)
			if err != nil {
				handleRPCError(w, err)
				return
			}

			posts[i].ModifiedAt, err = ptypes.Timestamp(singlePostResponse.ModifiedAt)
			if err != nil {
				handleRPCError(w, err)
				return
			}

			postStats, err := s.postStatsClient.GetPostStats(ctx, &poststats.GetPostStatsRequest{PostUid: posts[i].UID})
			if err != nil {
				handleRPCError(w, err)
				return
			}

			posts[i].NumLikes = postStats.NumLikes
			posts[i].NumDislikes = postStats.NumDislikes
			posts[i].NumViews = postStats.NumViews
		}

		json, err := json.Marshal(response{posts})
		if err != nil {
			handleRPCError(w, err)
			return
		}

		w.Write(json)
	}
}

func (s *Server) createPost() http.HandlerFunc {
	type request struct {
		Title string `json:"title"`
		URL   string `json:"url"`
	}

	type response struct {
		UID         string
		Title       string
		URL         string
		CreatedAt   time.Time
		ModifiedAt  time.Time
		NumLikes    int32
		NumDislikes int32
		NumViews    int32
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
		p, err := s.postClient.CreatePost(ctx, &post.CreatePostRequest{Title: req.Title, Url: req.URL})
		if err != nil {
			handleRPCError(w, err)
			return
		}

		_, err = s.postStatsClient.CreatePostStats(ctx, &poststats.CreatePostStatsRequest{PostUid: p.Uid})
		if err != nil {
			handleRPCError(w, err)
			return
		}

		createdAt, err := ptypes.Timestamp(p.CreatedAt)
		if err != nil {
			handleRPCError(w, err)
			return
		}

		modifiedAt, err := ptypes.Timestamp(p.ModifiedAt)
		if err != nil {
			handleRPCError(w, err)
			return
		}

		response := response{p.Uid, p.Title, p.Url, createdAt, modifiedAt, 0, 0, 0}
		json, err := json.Marshal(response)
		if err != nil {
			handleRPCError(w, err)
			return
		}

		w.Write(json)
	}
}

func (s *Server) getPost() http.HandlerFunc {
	type response struct {
		UID         string
		Title       string
		URL         string
		CreatedAt   time.Time
		ModifiedAt  time.Time
		NumLikes    int32
		NumDislikes int32
		NumViews    int32
	}

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		uid := vars["uid"]

		ctx := r.Context()
		postResponse, err := s.postClient.GetPost(ctx, &post.GetPostRequest{Uid: uid})
		if err != nil {
			handleRPCError(w, err)
			return
		}

		var res response
		res.UID = postResponse.Uid
		res.Title = postResponse.Title
		res.URL = postResponse.Url
		res.CreatedAt, err = ptypes.Timestamp(postResponse.CreatedAt)
		if err != nil {
			handleRPCError(w, err)
			return
		}

		res.ModifiedAt, err = ptypes.Timestamp(postResponse.ModifiedAt)
		if err != nil {
			handleRPCError(w, err)
			return
		}

		postStats, err := s.postStatsClient.GetPostStats(ctx, &poststats.GetPostStatsRequest{PostUid: res.UID})
		if err != nil {
			handleRPCError(w, err)
			return
		}

		res.NumLikes = postStats.NumLikes
		res.NumDislikes = postStats.NumDislikes
		res.NumViews = postStats.NumViews

		json, err := json.Marshal(res)
		if err != nil {
			handleRPCError(w, err)
			return
		}

		_, err = s.postStatsClient.IncreaseViews(ctx, &poststats.IncreaseViewsRequest{PostUid: uid})
		if err != nil {
			handleRPCError(w, err)
			return
		}

		w.Write(json)
	}
}

func (s *Server) updatePost() http.HandlerFunc {
	type request struct {
		Title string `json:"title"`
		URL   string `json:"url"`
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

		ctx := r.Context()
		_, err = s.postClient.UpdatePost(ctx, &post.UpdatePostRequest{Uid: uid, Title: req.Title, Url: req.URL})
		if err != nil {
			handleRPCError(w, err)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (s *Server) deletePost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		uid := vars["uid"]

		ctx := r.Context()
		_, err := s.postClient.DeletePost(ctx, &post.DeletePostRequest{Uid: uid})
		if err != nil {
			handleRPCError(w, err)
			return
		}

		_, err = s.postStatsClient.DeletePostStats(ctx, &poststats.DeletePostStatsRequest{PostUid: uid})
		if err != nil {
			handleRPCError(w, err)
			return
		}

		comments, err := s.commentClient.ListComments(ctx, &comment.ListCommentsRequest{PostUid: uid})
		if err != nil {
			handleRPCError(w, err)
			return
		}

		for _, c := range comments.Comments {
			s.commentClient.DeleteComment(ctx, &comment.DeleteCommentRequest{Uid: c.Uid})
			if err != nil {
				handleRPCError(w, err)
				return
			}
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (s *Server) likePost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		uid := vars["uid"]

		ctx := r.Context()
		_, err := s.postStatsClient.LikePost(ctx, &poststats.LikePostRequest{PostUid: uid})
		if err != nil {
			handleRPCError(w, err)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (s *Server) dislikePost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		uid := vars["uid"]

		ctx := r.Context()
		_, err := s.postStatsClient.DislikePost(ctx, &poststats.DislikePostRequest{PostUid: uid})
		if err != nil {
			handleRPCError(w, err)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

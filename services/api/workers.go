package api

import (
	"context"
	"log"
	"time"

	comment "github.com/andreymgn/RSOI/services/comment/proto"
	post "github.com/andreymgn/RSOI/services/post/proto"
	poststats "github.com/andreymgn/RSOI/services/poststats/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type workerRequest struct {
	uid      string
	doneTime time.Time
}

func (s *Server) deletePostWorker() {
	ctx := context.Background()
	for {
		req := <-s.deletePostChannel

		if time.Now().Before(req.doneTime) {
			s.deletePostChannel <- req
			continue
		}

		_, err := s.postClient.client.DeletePost(ctx,
			&post.DeletePostRequest{Token: s.postClient.token, Uid: req.uid},
		)
		if err != nil {
			if st, ok := status.FromError(err); ok {
				switch st.Code() {
				case codes.Unauthenticated:
					err := s.updatePostToken()
					if err != nil {
						panic(err)
					}
					_, err = s.postStatsClient.client.DeletePostStats(ctx,
						&poststats.DeletePostStatsRequest{Token: s.postStatsClient.token, PostUid: req.uid},
					)
					if err != nil {
						panic(err)
					}
				case codes.Unavailable:
					newReq := workerRequest{req.uid, time.Now().Add(time.Second * 5)}
					s.deletePostChannel <- newReq
					log.Printf("DeletePost rabotyaga: retrying %s", req.uid)
				}
			}
		}
	}
}

func (s *Server) deletePostStatsWorker() {
	ctx := context.Background()
	for {
		req := <-s.deletePostStatsChannel

		if time.Now().Before(req.doneTime) {
			s.deletePostStatsChannel <- req
			continue
		}

		_, err := s.postStatsClient.client.DeletePostStats(ctx,
			&poststats.DeletePostStatsRequest{Token: s.postStatsClient.token, PostUid: req.uid},
		)
		if err != nil {
			if st, ok := status.FromError(err); ok {
				switch st.Code() {
				case codes.Unauthenticated:
					err := s.updatePostStatsToken()
					if err != nil {
						panic(err)
					}
					_, err = s.postStatsClient.client.DeletePostStats(ctx,
						&poststats.DeletePostStatsRequest{Token: s.postStatsClient.token, PostUid: req.uid},
					)
					if err != nil {
						panic(err)
					}
				case codes.Unavailable:
					newReq := workerRequest{req.uid, time.Now().Add(time.Second * 5)}
					s.deletePostStatsChannel <- newReq
					log.Printf("DeletePostStats rabotyaga: retrying %s", req.uid)
				}
			}
		}
	}
}

func (s *Server) deleteCommentWorker() {
	ctx := context.Background()
	for {
		req := <-s.deleteCommentChannel

		if time.Now().Before(req.doneTime) {
			s.deleteCommentChannel <- req
			continue
		}

		_, err := s.commentClient.client.DeleteComment(ctx,
			&comment.DeleteCommentRequest{Uid: req.uid, Token: s.commentClient.token},
		)
		if err != nil {
			if st, ok := status.FromError(err); ok {
				switch st.Code() {
				case codes.Unauthenticated:
					err := s.updatePostToken()
					if err != nil {
						panic(err)
					}
					_, err = s.commentClient.client.DeleteComment(ctx,
						&comment.DeleteCommentRequest{Uid: req.uid, Token: s.commentClient.token},
					)
					if err != nil {
						panic(err)
					}
				case codes.Unavailable:
					newReq := workerRequest{req.uid, time.Now().Add(time.Second * 5)}
					s.deleteCommentChannel <- newReq
					log.Printf("DeleteComment rabotyaga: retrying %s", req.uid)
				}
			}
		}
	}
}
